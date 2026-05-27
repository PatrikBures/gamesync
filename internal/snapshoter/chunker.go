package snapshoter

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
	"github.com/restic/chunker"
)

func ChunkFilesInDir(repoDir string, chunkDir string) error {
	if err := filepath.WalkDir(repoDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		fmt.Println(path)

		if err := chunkFile(path, chunkDir); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

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
		hash := sha256.Sum256(chunk.Data)
		hashS := hex.EncodeToString(hash[:])

		chunkFilePath := filepath.Join(chunkDir, hashS)
		chunkFile, err := os.Create(chunkFilePath)
		if err != nil {
			return err
		}
		if err := writeChunk(chunk.Data, chunkFile); err != nil {
			return err
		}
		chunkCount++
	}
	
	return nil
}

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

