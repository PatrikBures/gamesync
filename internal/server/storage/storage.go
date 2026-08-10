package storage

import (
	"context"
	"io"
)

type Chunk interface {

	// Checks if chunk exists.
	Exists(ctx context.Context, hash string) (bool, error)

	// Writes compressed data from reader.
	//
	// Checks if the decompressed data hash matches provided hash.
	//
	// Returns number of uncompressed and compressed bytes the chunk has
	Store(ctx context.Context, hash string, data io.Reader) (int64, int64, error)

	// Provides a reader for chunk
	Download(ctx context.Context, hash string) (io.Reader, error)

	// Deletes chunk
	Delete(ctx context.Context, hash string) error
}
