package storage

import (
	"context"
	"encoding/hex"
	"errors"
	"gamesync/internal/server"
	"gamesync/internal/snapshoter"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
	"lukechampine.com/blake3"
)

type local struct {
	baseDir       string
	tmpDir        string
	maxChunkBytes int64
}

func (l *local) chunkFilePathDir(chunkFileDir string, hash string) string {
	return filepath.Join(chunkFileDir, hash)
}
func (l *local) chunkFilePath(hash string) string {
	return filepath.Join(l.chunkFileDir(hash), hash)
}
func (l *local) chunkFileDir(hash string) string {
	return filepath.Join(l.baseDir, snapshoter.DirsForChunk(2, 2, hash))
}

func NewLocal(baseDir string, maxChunkBytes int64) (*local, error) {
	tmpDir := filepath.Join(baseDir, ".tmp")
	if err := os.MkdirAll(tmpDir, 0775); err != nil {
		return nil, err
	}
	return &local{
		baseDir: baseDir,
		tmpDir: tmpDir,
		maxChunkBytes: maxChunkBytes,
	}, nil
}

func (l *local) Exists(ctx context.Context, hash string) (bool, error) {
	chunkFilePath := l.chunkFilePath(hash)
	d, err := os.Stat(chunkFilePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return false, server.NewInternalError(err, "failed stating chunk file")
	}
	if d != nil {
		return true, nil
	}
	return false, nil

}

func (l *local) Store(ctx context.Context, hash string, data io.Reader) (uncompressedBytes int64, err error) {
	tmpPath := filepath.Join(l.tmpDir, hash)

	tmpFile, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0664)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			// the file is already being uploaded so return no error
			return 0, nil
		}
		slog.Error("creating tmp file", "chunkHash", hash, "error", err)
		return 0, server.NewInternalError(err, "failed creating tmp file", "chunkHash", hash)
	}
	defer func() {
		if err != nil {
			if rerr := os.Remove(tmpPath); rerr != nil && !os.IsNotExist(rerr) {
				slog.Error("removing tmp file", "path", tmpPath, "error", rerr)
			}
		}
	}()


	uncompressedBytes, actualHashBytes, err := l.writeAndHash(data, tmpFile)
	if err != nil {
		return 0, err // it is fine to return the error here directly
	}

	// make sure the entire file is written to disk
	if err = tmpFile.Sync(); err != nil {
		return 0, server.NewInternalError(err, "failed syncing tmp chunk file")
	}

	if err = tmpFile.Close(); err != nil {
		slog.Error("closing chunk file", "error", err)
		return 0, server.NewInternalError(err, "failed closing chunk file")
	}

	actualHash := hex.EncodeToString(actualHashBytes)
	if actualHash != hash {
		slog.Warn("uploaded chunk hash does not match", "chunkHash", hash, "expectedChunkHash", actualHash)
		return 0, server.ErrHashMismatch
	}


	// move file to chunkdir
	chunkFileDir := l.chunkFileDir(hash)
	chunkFilePath := l.chunkFilePathDir(chunkFileDir, hash)
	if err = os.MkdirAll(chunkFileDir, 0775); err != nil {
		return 0, server.NewInternalError(err, "failed creating chunk file dir")
	}
	if err = os.Rename(tmpPath, chunkFilePath); err != nil {
		return 0, server.NewInternalError(err, "failed renaming chunk file to correct path")
	}

	// clear error just in case, so defer does not delete file
	err = nil
	return
}

func (l *local) Download(ctx context.Context, hash string) (io.Reader, error) {
	path := l.chunkFilePath(hash)
	f, err := os.Open(path)
	if err != nil {
		// do not add a specific check for if file exists, because the file should exist if it exists in db
		return nil, server.NewInternalError(err, "failed opening chunk file", "path", path)
	}
	return f, nil
}


// writes src to dst and returns Blake3 hash of zstd decompressed data, and the uncompressed amount of bytes
func (l *local) writeAndHash(src io.Reader, dst io.Writer) (int64, []byte, error) {
	dataCopy := io.TeeReader(src, dst)

	decoder, err := zstd.NewReader(dataCopy)
	if err != nil {
		return 0, nil, server.NewInternalError(err, "initializing zstd reader")
	}
	defer decoder.Close()

	hasher := blake3.New(32, nil)

	limitedData := &maxBytesReader{
		r: decoder,
		max: l.maxChunkBytes,
	}

	if _, err := io.Copy(hasher, limitedData); err != nil {
		if errors.Is(err, server.ErrChunkTooBig) {
			return 0, nil, err
		}
		return 0, nil, server.NewInternalError(err, "reading/hashing chunk")
	}

	return limitedData.n, hasher.Sum(nil), nil
}




type maxBytesReader struct {
	r   io.Reader
	max int64
	n   int64
}

func (m *maxBytesReader) Read(p []byte) (int, error) {
	n, err := m.r.Read(p)
	if n > 0 {
		m.n += int64(n)
		if m.n > m.max {
			slog.Warn("chunk too big", "reachedBytes", m.n, "limit", m.max)
			return n, server.ErrChunkTooBig
		}
	}
	return n, err
}
