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


type chunkHash struct {
	bytes [32]byte
	hex string
}

// handles the generation of chunks
//
// any chunks created will be in ChunkDir
type chunkGen struct {
	chunkDir string
	Info chunkGenInfo
}

type chunkGenInfo struct {
	FilesChunked  int64
	FilesErr      int64
	ChunksCreated int64
	ChunksSkipped int64
}

func NewChunkGen(chunkDir string) *chunkGen {
	return &chunkGen{chunkDir: chunkDir}
}


// chunks all files in repoDir
func (cg *chunkGen) ChunkFilesInDir(repoDir string) error {
	g := errgroup.Group{}
	g.SetLimit(runtime.NumCPU())

	if err := filepath.WalkDir(repoDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		g.Go(func() error {
			err := cg.chunkFile(path)
			if err != nil {
				atomic.AddInt64(&cg.Info.FilesErr, 1)
			} else {
				atomic.AddInt64(&cg.Info.FilesChunked, 1)
			}
			return err
		})

		return nil
	}); err != nil {
		return err
	}
	if err := g.Wait(); err != nil {
		return fmt.Errorf("chunking file: %w", err)
	}
	return nil
}

// creates chunk from the file at path to chunkDir
func (cg *chunkGen) chunkFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	c := chunker.New(f, chunker.Pol(0x3DA3358B4DC173))

	buf := make([]byte, 8*1024*1024)

	for chunkCount := 0; ; chunkCount++ {
		chunk, err := c.Next(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		chunkHash, chunkFile, err := cg.createChunk(chunk)
		if err != nil {
			return err
		}
		// this means the chunk already exists
		if chunkFile == nil {
			atomic.AddInt64(&cg.Info.ChunksSkipped, 1)
			continue
		}

		werr := cg.writeChunk(chunk.Data, chunkFile)
		cerr := chunkFile.Close()
		err = errors.Join(werr, cerr)
		if err != nil {
			return fmt.Errorf("creating chunk for file %s, chunk nr %d, hash %s: %w", path, chunkCount+1, chunkHash.hex, err)
		}
		atomic.AddInt64(&cg.Info.ChunksCreated, 1)
	}
	
	return nil
}

func (cg *chunkGen) createChunk(chunk chunker.Chunk) (chunkHash, *os.File, error) {
	ch := chunkHash{}
	ch.bytes = blake3.Sum256(chunk.Data)
	ch.hex = hex.EncodeToString(ch.bytes[:])

	hashDir := filepath.Join(cg.chunkDir, dirsForChunk(ch.hex))
	chunkFilePath := filepath.Join(hashDir, ch.hex)

	if err := os.MkdirAll(hashDir, 0775); err != nil {
		return ch, nil, fmt.Errorf("creating dir '%s': %w", hashDir, err)
	}

	chunkFile, err := os.OpenFile(chunkFilePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return ch, nil, nil
		}
		return ch, nil, err
	}
	return ch, chunkFile, nil
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
func dirsForChunk(hash string) (dirs string) {
	a := dirQty * dirLen
	parts := make([]string, dirQty)
	for i := 0; i < a; i += dirLen {
		parts[i/dirLen] = hash[i:i+dirLen]
	}
	return filepath.Join(parts...)
}
