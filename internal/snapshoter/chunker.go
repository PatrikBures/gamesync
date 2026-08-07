package snapshoter

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	api "go.pabu.dev/gamesync/internal/ogen"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/klauspost/compress/zstd"
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
		return fastcdc.New(nil, fastcdc.Opts{})
	},
}

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

type ChunkStream struct {
	Ch  <-chan FileResults
	err error
}

func (s *ChunkStream) Err() error {
	return s.err
}

// used to generate chunks
func NewChunkGen(chunkDir string) *chunkGen {
	return &chunkGen{chunkDir: chunkDir}
}

// chunks all files in repoDir
// 
// loop through the stream with ChunkStream.Ch
//
// after the loop, check any errors with ChunkStream.Err()
func (cg *chunkGen) ChunkFilesInDir(ctx context.Context, repoDir string) (*ChunkStream, error) {
	if info, err := os.Stat(repoDir); err != nil {
		return nil, fmt.Errorf("cannot access repo dir: %w", err)
	} else if !info.IsDir() {
		return nil, fmt.Errorf("repository is not a dir: %s", repoDir)
	}

	limit := runtime.NumCPU()
	g := errgroup.Group{}
	g.SetLimit(limit)

	out := make(chan FileResults, limit<<1)
	stream := &ChunkStream{Ch: out}

	go func() {
		defer close(out)

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
				fileHash, hashes, cerr := cg.chunkFile(path)
				relPath, ferr := filepath.Rel(repoDir, path)
				err := errors.Join(cerr, ferr)

				select {
				case out <- FileResults{Path: relPath, Hash: fileHash, Hashes: hashes, Err: err}:
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

		_ = g.Wait()
	}()

	return stream, nil
}

// creates chunks from the file at path to chunkDir
//
// and returns file hash and slice of the chunk hashes, including already existing ones
func (cg *chunkGen) chunkFile(path string) (hashHex string, chunkHashes []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", nil, err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	chunker := chunkFilePool.Get().(*fastcdc.Chunker)
	chunker.Reset(f)
	defer chunkFilePool.Put(chunker)

	fileHash := blake3.New(32, nil)

	for chunkCount := 0; ; chunkCount++ {
		chunk, err := chunker.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", nil, err
		}

		hashBytes, hashHex, chunkFile, err := cg.createChunk(chunk)
		if err != nil && !errors.Is(err, os.ErrExist) {
			return "", nil, err
		}

		if _, err = fileHash.Write(hashBytes); err != nil {
			return "", nil, fmt.Errorf("hashing chunk to form file hash: %w", errors.Join(err, chunkFile.Close()))
		}
		chunkHashes = append(chunkHashes, hashHex)

		// this means the chunk already exists
		if chunkFile == nil {
			atomic.AddInt64(&cg.Info.ChunksSkipped, 1)
			continue
		}

		werr := cg.writeChunk(chunk, chunkFile)
		cerr := chunkFile.Close()
		err = errors.Join(werr, cerr)
		if err != nil {
			return "", nil, fmt.Errorf("creating chunk for file %s, chunk nr %d, hash %s: %w", path, chunkCount+1, hashHex, err)
		}
		atomic.AddInt64(&cg.Info.ChunksCreated, 1)
	}

	return hex.EncodeToString(fileHash.Sum(nil)), chunkHashes, nil
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

// creates relative dirs in chunkDir and opens the
// file which it returns including its path.
//
// if the file already exist, error is os.ErrExists
func CreateChunkFile(chunkDir string, chunkHexHash string) (*os.File, string, error) {
	return createChunkFile(chunkDir, chunkHexHash, os.O_WRONLY|os.O_EXCL)
}

func CreateReadChunkFile(chunkDir string, chunkHexHash string) (*os.File, string, error) {
	return createChunkFile(chunkDir, chunkHexHash, os.O_RDWR)
}

func ReadChunkFile(chunkDir string, chunkHexHash string) (*os.File, string, error) {
	chunkFilePath := filepath.Join(chunkDir, DirsForChunk(chunkHexHash), chunkHexHash)

	chunkFile, err := os.Open(chunkFilePath)
	if err != nil {
		return nil, chunkFilePath, err
	}
	return chunkFile, chunkFilePath, nil
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
	}
	if err := stream.Err(); err != nil {
		return nil, fmt.Errorf("stream system error: %w", err)
	}

	return s, nil
}
