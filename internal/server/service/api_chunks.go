package service

import (
	"context"
	"encoding/hex"
	api "gamesync/internal/ogen"
	"gamesync/internal/server"
	"gamesync/internal/server/dbm"
	"log/slog"
)


func (s *Service) PutChunk(ctx context.Context, req api.PutChunkReq, params api.PutChunkParams) (api.PutChunkRes, error)  {
	if req.Data == nil {
		slog.Warn("got nil chunk", "chunkHash", params.ChunkHash)
		return &api.PutChunkNotAcceptable{}, nil
	}

	hash, err := hex.DecodeString(params.ChunkHash)
	if err != nil {
		slog.Warn("invalid chunk hash received", "error", err)
		return &api.PutChunkNotAcceptable{}, nil
	}

	chunkExists, err := s.storage.Exists(ctx, params.ChunkHash)
	if err != nil {
		return &api.PutChunkInternalServerError{}, server.ErrFile
	}
	if chunkExists {
		return &api.PutChunkConflict{}, nil
	}
	bytes, err := s.storage.Store(ctx, params.ChunkHash, req.Data)
	switch err {
	case nil:
		break
	case server.ErrHashMismatch: 
		return &api.PutChunkUnprocessableEntity{}, nil
	case server.ErrChunkTooBig:
		return &api.PutChunkRequestEntityTooLarge{}, nil
	default:
		return &api.PutChunkInternalServerError{}, err
	}


	s.db.WriteQuery().CreateChunk(ctx, dbm.CreateChunkParams{
		ChunkHash: hash,
		Bytes: bytes,
	})

	slog.Info("stored chunk", "chunkHash", params.ChunkHash[:8], "bytes", bytes)

	return &api.PutChunkCreated{}, nil
}

