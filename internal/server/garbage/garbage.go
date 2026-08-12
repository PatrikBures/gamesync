package garbage

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"go.pabu.dev/gamesync/internal/dbx"
	"go.pabu.dev/gamesync/internal/server"
	"go.pabu.dev/gamesync/internal/server/storage"
)


type garbageCollector struct {
	db *dbx.DB
	storage storage.Chunk
	opts Opts
}

type Opts struct {
	DeleteFiles        time.Duration // deletes files older than
	MarkChunks         time.Duration // marks chunks older than
	DeleteMarkedChunks time.Duration // deletes chunks marked for more than
}

func New(db *dbx.DB, storage storage.Chunk, opts Opts) *garbageCollector {
	return &garbageCollector{
		db: db,
		storage: storage,
		opts: opts,
	}
}

func (gc *garbageCollector) CleanOrphanedChunks(ctx context.Context) error {

	// delete files
	qtx, tx, err := gc.db.BeginTX(ctx)
	if err != nil {
		return server.NewInternalError(err, "failed starting tx")
	}
	defer func() {
		if recover() != nil || err != nil {
			if e := tx.Rollback(ctx); e != nil { slog.Error("failed rollback", "error", e) }
		}
	}()

	var filesDeletedCount int
	for {
		filesDelete, err := qtx.GarbageDeleteOrphanedFiles(ctx, gc.opts.DeleteFiles.Microseconds())
		if err != nil {
			return fmt.Errorf("deleting orphaned files: %w", err)
		}
		if len(filesDelete) == 0 {
			break
		}
		filesDeletedCount += len(filesDelete)
		time.Sleep(time.Millisecond*50)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commiting file deltions from db, %w", err)
	}
	slog.Info("deleted files from db", "count", filesDeletedCount, "olderThan", gc.opts.DeleteFiles)




	// mark chunks
	qtx, tx, err = gc.db.BeginTX(ctx)
	if err != nil {
		return server.NewInternalError(err, "failed starting tx")
	}
	defer func() {
		if recover() != nil || err != nil {
			if e := tx.Rollback(ctx); e != nil { slog.Error("failed rollback", "error", e) }
		}
	}()

	var markCount int64
	for {
		count, err := qtx.GarbageMarkOrphanedChunks(ctx, gc.opts.MarkChunks.Microseconds())
		if err != nil {
			return fmt.Errorf("marking chunks for deletion: %w", err)
		}
		if count == 0 {
			break
		}
		markCount += count
		time.Sleep(time.Millisecond*50)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commiting marking chunks for deletion, %w", err)
	}
	slog.Info("marked chunks for deletion", "count", markCount, "olderThan", gc.opts.MarkChunks)



	time.Sleep(gc.opts.DeleteMarkedChunks + time.Second*5)

	chunksDeletedCount, err := gc.deleteMarkedChunks(ctx)
	if err != nil {
		return fmt.Errorf("deleting marked chunks: %w", err)
	}

	slog.Info("deleted chunks", "count", chunksDeletedCount, "markedLongerThan", gc.opts.DeleteMarkedChunks)

	return nil
}


func (gc *garbageCollector) deleteMarkedChunks(ctx context.Context) (totalChunkDeletedCount int64, err error) {
	var deleteErr error

	qtx, tx, err := gc.db.BeginTX(ctx)
	if err != nil {
		return 0, server.NewInternalError(err, "failed starting tx")
	}
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic", "panic", r)
			if e := tx.Rollback(ctx); e != nil {
				slog.Error("failed rollback after panic", "error", e)
			}
		}
		if err != nil || deleteErr != nil {
			if e := tx.Rollback(ctx); e != nil {
				slog.Error("failed rollback on error", "error", e)
			}
		}
	}()

	for {
		chunksDelete, err := qtx.GarbageListMarkedChunks(ctx, gc.opts.DeleteMarkedChunks.Microseconds())
		if err != nil {
			return totalChunkDeletedCount, fmt.Errorf("orphaned chunk: %w", err)
		}
		if len(chunksDelete) == 0 {
			break
		}

		chunksDeletedCount := 0
		for _, chunkHash := range chunksDelete {
			if err = gc.storage.Delete(ctx, hex.EncodeToString(chunkHash)); err != nil {
				slog.Error("failed deleting chunk", "nr", chunksDeletedCount+1, "chunkHash", chunkHash)
				deleteErr = err
				break
			}
			chunksDeletedCount++
		}

		if chunksDeletedCount > 0 {
			chunksDeletedCountDB, err := qtx.GarbageDeleteChunks(ctx, chunksDelete[:chunksDeletedCount])
			if err != nil {
				return totalChunkDeletedCount, fmt.Errorf("deleting chunks from db: %w", err)
			}
			totalChunkDeletedCount += chunksDeletedCountDB
		}

		if deleteErr != nil {
			return totalChunkDeletedCount, fmt.Errorf("deleting chunk: %w", err)
		}

		time.Sleep(time.Millisecond*50)
	}

	if err = tx.Commit(ctx); err != nil {
		return totalChunkDeletedCount, fmt.Errorf("commiting marking chunks for deletion, %w", err)
	}

	return totalChunkDeletedCount, nil
}
