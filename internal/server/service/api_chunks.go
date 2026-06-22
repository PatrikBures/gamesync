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

	hash, err := hex.DecodeString(params.ChunkHash)
	if err != nil {
		slog.Warn("invalid chunk hash received", "error", err)
		return &api.PutChunkNotAcceptable{}, nil
	}

	chunkExists, err := s.db.ReadQuery().CheckChunk(ctx, hash)
	if err != nil {
		slog.Error("failed checking if chunk exists", "chunkHash", params.ChunkHash, "error", err)
		return &api.PutChunkInternalServerError{}, server.ErrFile
	}
	if chunkExists {
		return &api.PutChunkOK{}, nil
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


	if err := s.db.WriteQuery().CreateChunk(ctx, dbm.CreateChunkParams{
		ChunkHash: hash,
		Bytes: bytes,
	}); err != nil {
		slog.Error("inserting chunk to database", "chunkHash", params.ChunkHash, "error", err)
		return &api.PutChunkInternalServerError{}, server.ErrDatabase
		// TODO: should also delete the chunk file
	}

	return &api.PutChunkCreated{}, nil
}

