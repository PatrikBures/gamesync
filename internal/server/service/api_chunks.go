package service

import (
	"context"
	api "gamesync/internal/ogen"
	"gamesync/internal/server"
	"log/slog"
)


func (s *Service) PostChunk(ctx context.Context, req api.PostChunkReq, params api.PostChunkParams) (api.PostChunkRes, error)  {
	if req.Data == nil {
		slog.Warn("got nil chunk", "chunkHash", params.ChunkHash)
		return &api.PostChunkNotAcceptable{}, nil
	}

	chunkExists, err := s.storage.Exists(ctx, params.ChunkHash)
	if err != nil {
		return &api.PostChunkInternalServerError{}, server.ErrFile
	}
	if chunkExists {
		return &api.PostChunkConflict{}, nil
	}

	switch err := s.storage.Store(ctx, params.ChunkHash, req.Data); err {
	case nil:
		break
	case server.ErrHashMismatch: 
		return &api.PostChunkUnprocessableEntity{}, nil
	case server.ErrChunkTooBig:
		return &api.PostChunkRequestEntityTooLarge{}, nil
	default:
		return &api.PostChunkInternalServerError{}, err
	}

	return &api.PostChunkCreated{}, nil
}

