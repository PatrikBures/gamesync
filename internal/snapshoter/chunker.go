package snapshoter

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
	"github.com/restic/chunker"
	"lukechampine.com/blake3"
)

type chunkHash struct {
	bytes [32]byte
	hex string
}

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

// chunks all files in repoDir
//
// all chunks will be in chunkDir
func ChunkFilesInDir(repoDir string, chunkDir string) error {
	if err := filepath.WalkDir(repoDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if err := chunkFile(path, chunkDir); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

// creates chunk from the file at path to chunkDir
func chunkFile(path string, chunkDir string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	c := chunker.New(f, chunker.Pol(0x3DA3358B4DC173))

	buf := make([]byte, 8*1024*1024)

	chunkCount := 0
	for {
		chunk, err := c.Next(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		_, err = createChunk(chunk, chunkDir)
		if err != nil {
			return fmt.Errorf("creating chunk: %w", err)
		}
		chunkCount++
	}
	
	return nil
}

func createChunk(chunk chunker.Chunk, chunkDir string) (chunkHash, error) {
	ch := chunkHash{}
	ch.bytes = blake3.Sum256(chunk.Data)
	ch.hex = hex.EncodeToString(ch.bytes[:])

	hashDir := filepath.Join(chunkDir, dirsForChunk(ch.hex))
	chunkFilePath := filepath.Join(hashDir, ch.hex)

	if cf, _ := os.Stat(chunkFilePath); cf != nil {
		return ch, nil
	}

	if err := os.MkdirAll(hashDir, 0775); err != nil {
		return ch, fmt.Errorf("creating dir '%s': %w", hashDir, err)
	}

	chunkFile, err := os.Create(chunkFilePath)
	if err != nil {
		return ch, err
	}
	if err := writeChunk(chunk.Data, chunkFile); err != nil {
		return ch, err
	}
	return ch, nil
}

// writes bytes to out compressed using zstd
func writeChunk(bytes []byte, out io.Writer) (err error) {
	w, err := zstd.NewWriter(out)
	if err != nil {
		return err
	}
	defer func() {
		if errc := w.Close(); errc != nil {
			err = errors.Join(errc, err)
		}
	}()

	if _, err = w.Write(bytes); err != nil {
		return err
	}
	return nil
}

// returns the dirs that the chunk should be in based on dirQty and dirLen
func dirsForChunk(hash string) (dirs string) {
	a := dirQty * dirLen
	for i := 0; i < a; i += dirLen {
		dirs = filepath.Join(dirs, hash[i:i+dirLen])
	}
	return dirs
}
