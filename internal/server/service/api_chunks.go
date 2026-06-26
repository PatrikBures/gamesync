package service

import (
	"context"
	"encoding/hex"
	api "gamesync/internal/ogen"
	"gamesync/internal/server"
	"gamesync/internal/server/dbm"
)


func (s *Service) PutChunk(ctx context.Context, req api.PutChunkReq, params api.PutChunkParams) (api.PutChunkRes, error)  {

	hash, err := hex.DecodeString(params.ChunkHash)
	if err != nil {
		return nil, server.ErrInvalidHash
	}

	chunkExists, err := s.db.ReadQuery().CheckChunk(ctx, hash)
	if err != nil {
		return nil, server.NewInternalError(err, "failed checking if chunk exists", "chunkHash", params.ChunkHash, "error", err)
	}
	if chunkExists {
		return &api.PutChunkOK{}, nil
	}

	bytes, err := s.storage.Store(ctx, params.ChunkHash, req.Data)
	switch err {
	case nil:
		break
	case server.ErrHashMismatch, server.ErrChunkTooBig: 
		return nil, err
	default:
		return nil, server.NewInternalError(err, "failed storing chunk", "chunkHash", params.ChunkHash, "error", err)
	}


	if err := s.db.WriteQuery().CreateChunk(ctx, dbm.CreateChunkParams{
		ChunkHash: hash,
		Bytes: bytes,
	}); err != nil {
		return nil, server.NewInternalError(err, "failed inserting chunk to database", "chunkHash", params.ChunkHash, "error", err)
		// TODO: should also delete the chunk file
	}

	return &api.PutChunkCreated{}, nil
}

