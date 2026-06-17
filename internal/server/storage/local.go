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
	baseDir string
	tmpDir string
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

func NewLocal(baseDir string) (*local, error) {
	tmpDir := filepath.Join(baseDir, ".tmp")
	if err := os.MkdirAll(tmpDir, 0775); err != nil {
		return nil, err
	}
	return &local{
		baseDir: baseDir,
		tmpDir: tmpDir,
	}, nil
}

func (l *local) Exists(ctx context.Context, hash string) (bool, error) {
	chunkFilePath := l.chunkFilePath(hash)
	d, err := os.Stat(chunkFilePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		slog.Error("stating chunk file", "error", err)
		return false, server.ErrFile
	}
	if d != nil {
		return true, nil
	}
	return false, nil

}

func (l *local) Store(ctx context.Context, hash string, data io.Reader) (err error) {

	chunkFilePathTmp := filepath.Join(l.tmpDir, hash)

	chunkFileTmp, err := os.OpenFile(chunkFilePathTmp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0664)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			slog.Warn("multiple uploads of same chunk occured at the same time", "chunkHash", hash)
			return server.ErrDuplicateKey
		}
		slog.Error("opening chunk", "chunkHash", hash, "error", err)
		return server.ErrFile
	}
	defer func() {
		if chunkFileTmp != nil {
			cerr := chunkFileTmp.Close()
			if cerr != nil {
				slog.Error("closing chunk file", "error", cerr)
				err = errors.Join(err, server.ErrFileClose)
			}
		}
	}()

	dataCopy := io.TeeReader(data, chunkFileTmp)

	decoder, err := zstd.NewReader(dataCopy)
	if err != nil {
		slog.Error("starting ztsd reader", "chunkHash", hash, "error", err)
		return server.ErrDecompress
	}
	defer decoder.Close()

	hasher := blake3.New(32, nil)

	if _, err := io.Copy(hasher, decoder); err != nil {
		slog.Error("reading/hashing chunk", "chunkHash", hash, "error", err)
		return server.ErrHashing
	}

	defer func() {
		if err == nil {
			return
		}
		rerr := os.Remove(chunkFilePathTmp) 
		if rerr != nil {
			slog.Error("failed removing file", "error", rerr)
		}
	}()

	if err := chunkFileTmp.Sync(); err != nil {
		slog.Error("syncing tmp file", "error", err)
		return server.ErrFile
	}

	if err = chunkFileTmp.Close(); err != nil {
		slog.Error("closing chunk file", "error", err)
		return server.ErrFileClose
	}
	chunkFileTmp = nil

	hashExpected := hasher.Sum(nil)
	hashHex := hex.EncodeToString(hashExpected)
	if hashHex != hash {
		slog.Warn("uploaded chunk hash does not match", "chunkHash", hash, "expectedChunkHash", hashHex)
		return server.ErrHashMismatch
	}


	chunkFileDir := l.chunkFileDir(hash)
	chunkFilePath := l.chunkFilePathDir(chunkFileDir, hash)

	// creates dir and rename file to the correct spot
	if err := os.MkdirAll(chunkFileDir, 0775); err != nil {
		slog.Error("creating chunk file dir", "error", err)
		return server.ErrFile
	}
	if err := os.Rename(chunkFilePathTmp, chunkFilePath); err != nil {
		slog.Error("renaming chunk", "error", err)
		return server.ErrFile
	}

	return nil
}

