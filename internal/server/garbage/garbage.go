package garbage

import (
	"context"

	"go.pabu.dev/gamesync/internal/dbx"
)


type garbageCollector struct {
	db *dbx.DB
}

func New(db *dbx.DB) *garbageCollector {
	return &garbageCollector{
		db: db,
	}
}

func (gc *garbageCollector) CleanOrphanedChunks(ctx context.Context) error  {
	return nil
}

