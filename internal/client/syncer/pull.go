package syncer

import (
	"context"
	"errors"
	"fmt"
	api "go.pabu.dev/gamesync/internal/ogen"
	"go.pabu.dev/gamesync/internal/snapshoter"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
)

func (s *syncer) Pull(snapshotID int64) (err error) {

	snapshot, err := s.client.GetSnapshot(context.Background(), api.GetSnapshotParams{
		UserID: s.conf.Server.UserID,
		RepoName: s.profile.RepoName,
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
			client: s.client,
			chunkDir: s.conf.Global.ChunkDir,
		}
		// if it errors after this point, it is quite likely that it leaves the file in a broken state
		puller.file, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0664)
		if err != nil {
			return fmt.Errorf("opening file: %w", err)
		}

		puller.decoder, err = zstd.NewReader(nil)
		if err != nil {
			return fmt.Errorf("initializing zstd decoder: %w", err)
		}
		for _, chunkHash := range f.ChunkHashes {
			if err := puller.PullChunk(chunkHash); err != nil {
				return fmt.Errorf("pulling chunk for file '%s', chunk '%s': %w", filePath, chunkHash, err)
			}
		}
	}

	return nil
}


func (p *puller) PullChunk(chunkHash string) (err error) {

	var chunkReader io.Reader
	chunkExists := false

	chunkFile, _, err := snapshoter.CreateReadChunkFile(p.chunkDir, chunkHash)
	if err != nil {
		return fmt.Errorf("opening chunkfile: %w", err)
	}
	if chunkFile == nil {
		return fmt.Errorf("chunkFile is nil!")
	}
	defer func() {
		err = errors.Join(err, chunkFile.Close())
	}()

	if chunkStat, err := chunkFile.Stat(); err != nil {
		return fmt.Errorf("stating chunk file: %w", err)
	} else if chunkStat.Size() != 0 {
		chunkExists = true
	}

	fileWriterOffset := io.NewOffsetWriter(p.file, p.fileBytesWritten)

	var chunkWriter io.Writer
	// if the chunk does not exist, it will write the chunk to the chunk cache and the snapshot file
	// otherwise it will read from the existing cache into the snapshot file
	if !chunkExists {
		chunk, err := p.client.GetChunk(context.Background(), api.GetChunkParams{ChunkHash: chunkHash})
		if err != nil {
			return fmt.Errorf("failed getting chunk: %w", err)
		}
		chunkReader = io.TeeReader(chunk.Data, chunkFile)
		chunkWriter = fileWriterOffset
		fmt.Print("!")
	} else {
		chunkReader = chunkFile
		chunkWriter = fileWriterOffset
		fmt.Print("#")
	}

	if chunkWriter == nil {
		return fmt.Errorf("chunkWriter is nil!")
	}
	if chunkReader == nil {
		return fmt.Errorf("chunkReader is nil!")
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
