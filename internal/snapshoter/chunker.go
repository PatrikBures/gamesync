package snapshoter

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"go.pabu.dev/gamesync/internal/client/stater"
	api "go.pabu.dev/gamesync/internal/ogen"

	"github.com/klauspost/compress/zstd"
	"github.com/schollz/progressbar/v3"
	"go.pabu.dev/fastcdc"
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

var chunkFilePool = sync.Pool{
	New: func() any {
		return new(poolData)
	},
}

type poolData struct {
	chunker *fastcdc.Chunker
	encoder *zstd.Encoder
}

// handles the generation of chunks
//
// any chunks created will be in ChunkDir
type chunkGen struct {
	chunkDir string
	bar *progressbar.ProgressBar
	totalFileCount int64
	totalFilesProcessed int64
	Info chunkGenInfo
}

type chunkGenInfo struct {
	chunkCreated atomic.Int64
	chunkSkipped atomic.Int64
}
func (cgi *chunkGenInfo) Created() int64 { return cgi.chunkCreated.Load() }
func (cgi *chunkGenInfo) Skipped() int64 { return cgi.chunkSkipped.Load() }
func (cgi *chunkGenInfo) Print() {
	fmt.Printf("New created chunks: %d, Already existing chunks: %d\n", cgi.Created(),  cgi.Skipped())
}



// contains the path, hashes and error obtained when chunking file
type FileResults struct {
	Hash   string
	Path   string
	Hashes []string
	State  stater.FileState
	Err    error
}

type ChunkStream struct {
	Ch  <-chan FileResults
	err error
}

func (s *ChunkStream) Err() error {
	return s.err
}

// used to generate chunks
//
// remember to exit chunkGen and to not reuse the struct
func NewChunkGen(chunkDir string) *chunkGen {
	return &chunkGen{chunkDir: chunkDir}
}

// exits bar
func (cg *chunkGen) exit() error {
	if cg.bar == nil {
		return nil
	}
	err := cg.bar.Exit()
	cg.bar = nil
	fmt.Println()
	return err
}

func (cg *chunkGen) initBar(repoDir string) error {
	if cg.bar != nil {
		return fmt.Errorf("bar is not nil, Exit or Finish chunkGen")
	}

	totalSize, totalFileCount, err := dirInfo(repoDir)
	if err != nil {
		return err
	}
	cg.totalFileCount = totalFileCount

	cg.bar = progressbar.NewOptions64(totalSize,
		progressbar.OptionSetDescription(fmt.Sprintf("[0/%d] Starting chunking...", cg.totalFileCount)),
		progressbar.OptionShowCount(),
		progressbar.OptionShowBytes(true),
		progressbar.OptionUseIECUnits(true),
	)

	return nil
}

// increases counter
func (cg *chunkGen) ProcessedFile() {
	if cg.bar == nil {
		return
	}
	cg.totalFilesProcessed++
	cg.bar.Describe(fmt.Sprintf("[%d/%d] Chunking files", cg.totalFilesProcessed, cg.totalFileCount))
}

// Chunks all files in repoDir. Loop through the stream with ChunkStream.Ch
//
// During the loop check each FileResult error
// and after the loop, check any errors with ChunkStream.Err()
func (cg *chunkGen) ChunkFilesInDir(ctx context.Context, repoDir string) (*ChunkStream, error) {

	if info, err := os.Stat(repoDir); err != nil {
		return nil, fmt.Errorf("cannot access repo dir: %w", err)
	} else if !info.IsDir() {
		return nil, fmt.Errorf("repository is not a dir: %s", repoDir)
	}

	if err := cg.initBar(repoDir); err != nil {
		return nil, fmt.Errorf("initializing bar: %w", err)
	}

	if cg.totalFileCount == 0 {
		emptyCh := make(chan FileResults)
		close(emptyCh)
		return &ChunkStream{Ch: emptyCh, err: nil}, nil
	}

	limit := 1
	g := errgroup.Group{}
	g.SetLimit(limit)

	out := make(chan FileResults, limit<<1)
	stream := &ChunkStream{Ch: out}

	go func() {
		defer close(out)
		defer func() { _ = cg.exit() }()

		walkErr := filepath.WalkDir(repoDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if d.IsDir() {
				return nil
			}
			g.Go(func() error {
				fileHash, hashes, fileState, chunkErr := cg.chunkFile(path)

				relPath, relErr := filepath.Rel(repoDir, path)

				err := errors.Join(chunkErr, relErr)

				select {
				case out <- FileResults{
					Path: relPath,
					Hash: fileHash,
					Hashes: hashes,
					State: fileState,
					Err: err,
				}:
				case <-ctx.Done():
				}

				return nil
			})
			return nil
		})

		if walkErr != nil {
			stream.err = fmt.Errorf("dir traversal failed: %w", walkErr)
			return
		}

		cg.ProcessedFile()

		_ = g.Wait()
	}()

	return stream, nil
}

// creates chunks from the file at path to chunkDir
//
// returns file hash and slice of the chunk hashes (including already existing) and the fileState.
func (cg *chunkGen) chunkFile(path string) (string, []string, stater.FileState, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", nil, stater.FileState{}, err
	}
	r := progressbar.NewReader(f, cg.bar)

	pd := chunkFilePool.Get().(*poolData)
	defer chunkFilePool.Put(pd)
	if pd.encoder == nil {
		enc, err := zstd.NewWriter(nil)
		if err != nil {
			return "", nil, stater.FileState{}, fmt.Errorf("initializing zstd encoder: %w", err)
		}
		pd.encoder = enc
	}
	if pd.chunker == nil {
		pd.chunker = fastcdc.New(nil, fastcdc.Opts{})
	}
	pd.chunker.Reset(&r)

	fileHash := blake3.New(32, nil)

	var chunkHashes []string

	for chunkCount := 0; ; chunkCount++ {
		chunk, err := pd.chunker.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", nil,stater.FileState{}, err
		}

		hashBytes, hashHex, chunkFile, err := cg.createChunk(chunk)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return "", nil, stater.FileState{}, err
		}

		if _, err = fileHash.Write(hashBytes); err != nil {
			return "", nil, stater.FileState{}, fmt.Errorf("hashing chunk to form file hash: %w", errors.Join(err, chunkFile.Close()))
		}
		chunkHashes = append(chunkHashes, hashHex)

		// this means the chunk already exists
		if chunkFile == nil {
			cg.Info.chunkSkipped.Add(1)
			continue
		}

		pd.encoder.Reset(chunkFile)
		werr := cg.writeChunk(chunk, pd.encoder)
		cerr := chunkFile.Close()
		err = errors.Join(werr, cerr)
		if err != nil {
			return "", nil, stater.FileState{}, fmt.Errorf("creating chunk for file %s, chunk nr %d, hash %s: %w", path, chunkCount+1, hashHex, err)
		}
		cg.Info.chunkCreated.Add(1)
	}

	stat, err := f.Stat()
	if err != nil {
		return "", nil, stater.FileState{}, err
	}
	fileState := stater.FileState{
		ModTime: stat.ModTime().Unix(),
		Size: stat.Size(),
	}

	return hex.EncodeToString(fileHash.Sum(nil)), chunkHashes, fileState, nil
}

// hashes chunk, creates the relative dirs, and opens the file (does not put any data in)
func (cg *chunkGen) createChunk(chunk []byte) ([]byte, string, *os.File, error) {
	hashBytes := blake3.Sum256(chunk)
	hashSlice := hashBytes[:]
	hashHex := hex.EncodeToString(hashSlice)

	chunkFile, _, err := CreateChunkFile(cg.chunkDir, hashHex)

	return hashSlice, hashHex, chunkFile, err
}

func createChunkFile(chunkDir string, chunkHexHash string, flags int) (*os.File, string, error) {
	hashDir := filepath.Join(chunkDir, DirsForChunk(chunkHexHash))
	chunkFilePath := filepath.Join(hashDir, chunkHexHash)

	if err := os.MkdirAll(hashDir, 0775); err != nil {
		return nil, chunkFilePath, err
	}

	chunkFile, err := os.OpenFile(chunkFilePath, os.O_CREATE|flags, 0664)
	if err != nil {
		return nil, chunkFilePath, err
	}
	return chunkFile, chunkFilePath, nil
}

// Creates relative dirs in chunkDir and opens the
// file which it returns including its path.
//
// if the file already exist, error is os.ErrExists
//
// File is opened in readonly mode
func CreateChunkFile(chunkDir string, chunkHexHash string) (*os.File, string, error) {
	return createChunkFile(chunkDir, chunkHexHash, os.O_WRONLY|os.O_EXCL)
}

// Creates relative dirs in chunkDir and opens the
// file which it returns including its path.
// The file is created if it already does not exist. 
// 
// File is opened int read and write mode.
func CreateReadChunkFile(chunkDir string, chunkHexHash string) (*os.File, string, error) {
	return createChunkFile(chunkDir, chunkHexHash, os.O_RDWR)
}

// Opens chunkfile in read only mode also returning its path.
func ReadChunkFile(chunkDir string, chunkHexHash string) (*os.File, string, error) {
	chunkFilePath := filepath.Join(chunkDir, DirsForChunk(chunkHexHash), chunkHexHash)

	chunkFile, err := os.Open(chunkFilePath)
	if err != nil {
		return nil, chunkFilePath, err
	}
	return chunkFile, chunkFilePath, nil
}

// writes bytes to encoder and closes it
func (cg *chunkGen) writeChunk(bytes []byte, encoder *zstd.Encoder) error {
	_, writeErr := encoder.Write(bytes)
	closeErr := encoder.Close()
	if writeErr != nil {
		writeErr = fmt.Errorf("writing chunk: %w", writeErr)
	}
	if closeErr != nil {
		closeErr = fmt.Errorf("closing chunk: %w", closeErr)
	}

	return errors.Join(writeErr, closeErr)
}

// returns the dirs that the chunk should be in
func DirsForChunk(hash string) (dirs string) {
	dq := dirQty
	dl := dirLen
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

func (cg *chunkGen) ChunkFilesInDirSlice(ctx context.Context, repoDir string) ([]FileResults, error) {
	stream, err := cg.ChunkFilesInDir(ctx, repoDir)
	if err != nil {
		return nil, err
	}

	s := []FileResults{}

	for fr := range stream.Ch {
		if fr.Err != nil {
			return nil, fmt.Errorf("chunk '%s': %w", fr.Hash, fr.Err)
		}
		if fr.Hash == "" {
			return nil, fmt.Errorf("file hash is empty: %s", fr.Path)
		}
		s = append(s, fr)
		cg.ProcessedFile()
	}
	if err := stream.Err(); err != nil {
		return nil, fmt.Errorf("stream system error: %w", err)
	}

	return s, nil
}

func (cg *chunkGen) ChunkFilesInDirApiFile(ctx context.Context, repoDir string) ([]api.File, error) {
	stream, err := cg.ChunkFilesInDir(ctx, repoDir)
	if err != nil {
		return nil, err
	}

	s := []api.File{}

	for fr := range stream.Ch {
		if fr.Err != nil {
			return nil, fmt.Errorf("chunk '%s': %w", fr.Hash, fr.Err)
		}
		if fr.Hash == "" {
			return nil, fmt.Errorf("file hash is empty: %s", fr.Path)
		}
		s = append(s, api.File{
			ChunkHashes: fr.Hashes,
			Hash:        fr.Hash,
			Path:        fr.Path,
		})
		cg.ProcessedFile()
	}
	if err := stream.Err(); err != nil {
		return nil, fmt.Errorf("stream system error: %w", err)
	}

	return s, nil
}



func dirInfo(dir string) (totalSize int64, fileCount int64, err error) {
	if err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		totalSize += info.Size()
		fileCount++
		
		return nil
	}); err != nil {
		return 0, 0, err
	}

	return totalSize, fileCount, nil
}
