package storage

import (
	"context"
	"io"
)

type Chunk interface {

	// Checks if chunk exists.
	//
	// if there was an error, then error is not nil and bool is false.
	Exists(ctx context.Context, hash string) (bool, error)

	// Writes compressed data from reader.
	//
	// Checks if the decompressed data hash matches provided hash.
	//
	// Returns number of uncompressed bytes the chunk has
	Store(ctx context.Context, hash string, data io.Reader) (int64, error)

	// Provides a reader for chunk
	Download(ctx context.Context, hash string) (io.Reader, error)
}
