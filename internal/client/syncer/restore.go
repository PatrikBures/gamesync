package syncer

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	api "go.pabu.dev/gamesync/internal/ogen"
	"go.pabu.dev/gamesync/internal/snapshoter"
	"go.pabu.dev/packer"

	"github.com/klauspost/compress/zstd"
)

type puller struct {
	client           *api.Client
	decoder          *zstd.Decoder
	p                *packer.Packager
	buf              []byte
	fileBytesWritten int64
}

func (s *syncer) Restore(snapshotID int64, currentFilePaths map[string]struct{}) (err error) {

	snapshot, err := s.client.GetSnapshot(context.Background(), api.GetSnapshotParams{
		UserID:     s.conf.Server.UserID,
		RepoName:   s.profile.RepoName,
		BranchName: s.profile.BranchName,
		SnapshotID: snapshotID,
	})
	if err != nil {
		return fmt.Errorf("getting snapshot: %w", err)
	}

	// we can add one puller which is only used if there is one chunk in the file
	// this is to reuse the decoder and reduce allocations. (when goroutines are added)

	// if it errors anywhere in the loop. then the sync will not be up to date. if it errors there
	// should be a defer which restores its state the the previous snapshot.
	for _, f := range snapshot.Files {
		// currently only downloads the chunks. needs to also create the files from the chunks
		// when checking if the files exists, it will just delete it and write a new one for now.
		// until there is a state which stores the current chunk hash, size and modtime, which
		// can be used to quickly check if it has changed.
		//
		// when it writes the chunk, it needs to write the chunk and also append to the file.
		// if the chunk already exists, it just needs to append to the file.

		filePath := filepath.Join(s.profile.Dir, f.Path)

		fmt.Printf("\nwriting to: %s\n%s\n", filePath, f.Hash)

		dir := filepath.Dir(filePath)
		if err = os.MkdirAll(dir, 0775); err != nil {
			return fmt.Errorf("creating dir: %w", err)
		}

		if len(f.ChunkHashes) == 0 {
			if _, err := os.OpenFile(filePath, os.O_CREATE, 0664); err != nil {
				return fmt.Errorf("creating empty file: %w", err)
			}
			fmt.Println("created emtpy file:", filePath)
			continue
		}

		puller := puller{
			client:   s.client,
			p:        s.p,
			buf:      make([]byte, 128<<10),
		}

		puller.decoder, err = zstd.NewReader(nil)
		if err != nil {
			return fmt.Errorf("initializing zstd decoder: %w", err)
		}
		for _, chunkHash := range f.ChunkHashes {
			chunkHashHex, err := hex.DecodeString(chunkHash)
			if err != nil {
				return fmt.Errorf("decoding chunk hash: %w", err)
			}

			if err := puller.PullChunk(chunkHashHex); err != nil {
				return fmt.Errorf("pulling chunk for file '%s', chunk '%s': %w", filePath, chunkHash, err)
			}
		}

		delete(currentFilePaths, filePath)
	}

	return cleanUp(currentFilePaths)
}

func (p *puller) PullChunk(chunkHash []byte) (error) {

	var chunk []byte
	chunkExists := false

	if localChunk, err := p.p.Get(chunkHash, p.buf); !errors.Is(err, packer.ErrChunkNotFound) && err != nil {
		return fmt.Errorf("getting chunk localy: %w", err)
	} else if err == nil {
		chunk = localChunk
		chunkExists = true
	}

	if !chunkExists {
		chunkReader, err := p.client.GetChunk(context.Background(), api.GetChunkParams{ChunkHash: hex.EncodeToString(chunkHash)})
		if err != nil {
			return fmt.Errorf("failed getting chunk: %w", err)
		}
		fmt.Print("!")
		readBytes, err := io.ReadFull(chunkReader, p.buf)
		if err != nil {
			return fmt.Errorf("reading chunk: %w", err)
		}
		chunk = p.buf[:readBytes]
	} else {
		fmt.Print("#")
	}


	// need to verify that the hash matches the content
	if err := p.decoder.Reset(chunkReader); err != nil {
		return fmt.Errorf("reseting decoder: %w", err)
	}

	bytesWritten, err := p.decoder.WriteTo(chunkWriter)
	if err != nil {
		return fmt.Errorf("decoding chunk, wrote %d bytes: %w", bytesWritten, err)
	}

	p.fileBytesWritten += bytesWritten

	fmt.Println("Downloaded:", chunkHash)

	return nil
}

func cleanUp(deleteFiles map[string]struct{}) error {
	for path := range deleteFiles {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("removing file: %w", err)
		}
		dir := filepath.Dir(path)
		if err := os.Remove(dir); err != nil && !errors.Is(err, os.ErrExist) {
			return fmt.Errorf("removing empty dir: %w", err)
		}
	}
	return nil
}
