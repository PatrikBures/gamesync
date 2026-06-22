package snapshoter

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"

	"github.com/klauspost/compress/zstd"
	"github.com/restic/chunker"
	"golang.org/x/sync/errgroup"
	"lukechampine.com/blake3"
)

// how many dirs each hash will be in
//
// example:
//
// 2: as/df/asdf1234
//
// 3: as/df/12/asdf1234
const dirQty = 2

// how many chars each dir will have, is probably better that is to the power of 2
//
// example:
//
// 2: as/asdf1234
//
// 4: asdf/asdf1234
const dirLen = 2

// handles the generation of chunks
//
// any chunks created will be in ChunkDir
type chunkGen struct {
	chunkDir string
	Info     chunkGenInfo
}

type chunkGenInfo struct {
	FilesChunked  int64
	FilesErr      int64
	ChunksCreated int64
	ChunksSkipped int64
}

// contains the path, hashes and error obtained when chunking file
type FileResults struct {
	Hash   string
	Path   string
	Hashes []string
	Err    error
}

// used to generate chunks
func NewChunkGen(chunkDir string) *chunkGen {
	return &chunkGen{chunkDir: chunkDir}
}

// chunks all files in repoDir
func (cg *chunkGen) ChunkFilesInDir(repoDir string) ([]FileResults, error) {
	g := errgroup.Group{}
	g.SetLimit(runtime.NumCPU())
	results := make(chan FileResults, 10)
	waitErr := make(chan error, 1)

	if err := filepath.WalkDir(repoDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		g.Go(func() error {
			fileHash, hashes, err := cg.chunkFile(path)
			results <- FileResults{Path: path, Hash: fileHash, Hashes: hashes, Err: err}
			return err
		})

		return nil
	}); err != nil {
		return nil, err
	}
	go func() {
		waitErr <- g.Wait()
		close(results)
	}()

	var all []FileResults
	for r := range results {
		if r.Err == nil {
			cg.Info.FilesChunked++
		} else {
			cg.Info.FilesErr++
		}
		all = append(all, r)
	}
	if gErr := <-waitErr; gErr != nil {
		return all, fmt.Errorf("chunking files: %w", gErr)
	}
	return all, nil
}

// creates chunks from the file at path to chunkDir
//
// and returns file hash and slice of the chunk hashes, including already existing ones
func (cg *chunkGen) chunkFile(path string) (string, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", nil, err
	}
	c := chunker.New(f, chunker.Pol(0x3DA3358B4DC173))

	buf := make([]byte, 8*1024*1024)

	chunkHashes := []string{}
	fileHash := blake3.New(32, nil)

	for chunkCount := 0; ; chunkCount++ {
		chunk, err := c.Next(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return "", nil, err
		}

		hashBytes, hashHex, chunkFile, err := cg.createChunk(chunk)
		if err != nil {
			return "", nil, err
		}
		// this means the chunk already exists
		if chunkFile == nil {
			atomic.AddInt64(&cg.Info.ChunksSkipped, 1)
			chunkHashes = append(chunkHashes, hashHex)
			continue
		}

		werr := cg.writeChunk(chunk.Data, chunkFile)
		cerr := chunkFile.Close()
		err = errors.Join(werr, cerr)
		if err != nil {
			return "", nil, fmt.Errorf("creating chunk for file %s, chunk nr %d, hash %s: %w", path, chunkCount+1, hashHex, err)
		}
		if _, err := fileHash.Write(hashBytes); err != nil {
			return "", nil, fmt.Errorf("hashing chunk to form file hash: %w", err)
		}
		chunkHashes = append(chunkHashes, hashHex)
		atomic.AddInt64(&cg.Info.ChunksCreated, 1)
	}
	
	return hex.EncodeToString(fileHash.Sum(nil)), chunkHashes, nil
}

// hashes chunk, creates the relative dirs, and opens the file (does not put any data in)
func (cg *chunkGen) createChunk(chunk chunker.Chunk) ([]byte, string, *os.File, error) {
	hashBytes := blake3.Sum256(chunk.Data)
	hashSlice := hashBytes[:]
	hashHex := hex.EncodeToString(hashSlice)

	chunkFile, err := CreateChunkFile(cg.chunkDir, hashHex); 

	return hashSlice, hashHex, chunkFile, err
}

// creates relative dirs in chunkDir and opens the file.
//
// if file already exists, *os.File and error are nil.
func CreateChunkFile(chunkDir string, chunkHexHash string) (*os.File, error) {
	hashDir := filepath.Join(chunkDir, DirsForChunk(dirLen, dirQty, chunkHexHash))
	chunkFilePath := filepath.Join(hashDir, chunkHexHash)

	if err := os.MkdirAll(hashDir, 0775); err != nil {
		return nil, err
	}

	chunkFile, err := os.OpenFile(chunkFilePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0664)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, nil
		}
		return nil, err
	}
	return chunkFile, nil
}

// writes bytes to out compressed using zstd
func (cg *chunkGen) writeChunk(bytes []byte, out io.Writer) error {
	w, err := zstd.NewWriter(out)
	if err != nil {
		return fmt.Errorf("creating writer: %w", err)
	}

	_, writeErr := w.Write(bytes)
	closeErr := w.Close()
	if writeErr != nil {
		writeErr = fmt.Errorf("writing chunk: %w", writeErr)
	}
	if closeErr != nil {
		closeErr = fmt.Errorf("closing chunk: %w", closeErr)
	}

	return errors.Join(writeErr, closeErr)
}

// returns the dirs that the chunk should be in based on dirQty and dirLen
func DirsForChunk(dl int, dq int, hash string) (dirs string) {
	a := dq * dl
	parts := make([]string, dq)
	for i := 0; i < a; i += dl {
		parts[i/dl] = hash[i : i+dl]
	}
	return filepath.Join(parts...)
}

func PrintFileResultsErrors(result []FileResults) {
	for _, r := range result {
		if r.Err == nil {
			continue
		}
		fmt.Fprintf(os.Stderr, "Error chunking file path: %s, err: %v\n", r.Path, r.Err)
	}
}
