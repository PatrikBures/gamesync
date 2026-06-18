package service

import (
	"context"
	"encoding/hex"
	api "gamesync/internal/ogen"
	"gamesync/internal/server"
	"gamesync/internal/server/dbm"
	"log/slog"
)


func (s *Service) PostChunk(ctx context.Context, req api.PostChunkReq, params api.PostChunkParams) (api.PostChunkRes, error)  {
	if req.Data == nil {
		slog.Warn("got nil chunk", "chunkHash", params.ChunkHash)
		return &api.PostChunkNotAcceptable{}, nil
	}

	hash, err := hex.DecodeString(params.ChunkHash)
	if err != nil {
		slog.Warn("invalid chunk hash received", "error", err)
		return &api.PostChunkNotAcceptable{}, nil
	}

	chunkExists, err := s.storage.Exists(ctx, params.ChunkHash)
	if err != nil {
		return &api.PostChunkInternalServerError{}, server.ErrFile
	}
	if chunkExists {
		return &api.PostChunkConflict{}, nil
	}
	bytes, err := s.storage.Store(ctx, params.ChunkHash, req.Data)
	switch err {
	case nil:
		break
	case server.ErrHashMismatch: 
		return &api.PostChunkUnprocessableEntity{}, nil
	case server.ErrChunkTooBig:
		return &api.PostChunkRequestEntityTooLarge{}, nil
	default:
		return &api.PostChunkInternalServerError{}, err
	}


	s.db.WriteQuery().CreateChunk(ctx, dbm.CreateChunkParams{
		ChunkHash: hash,
		Bytes: bytes,
	})

	slog.Info("stored chunk", "chunkHash", params.ChunkHash[:8], "bytes", bytes)

	return &api.PostChunkCreated{}, nil
}

