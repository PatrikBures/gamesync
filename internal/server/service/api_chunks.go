package service

import (
	"context"
	"encoding/hex"
	api "go.pabu.dev/gamesync/internal/ogen"
	"go.pabu.dev/gamesync/internal/server"
	"go.pabu.dev/gamesync/internal/server/dbm"
)

func (s *Service) PutChunk(ctx context.Context, req api.PutChunkReq, params api.PutChunkParams) (api.PutChunkRes, error) {

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

	bytes, bytesCompressed, err := s.storage.Store(ctx, params.ChunkHash, req.Data)
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
		Bytes:     bytes,
		BytesCompressed: bytesCompressed,
	}); err != nil {
		derr := s.storage.Delete(ctx, params.ChunkHash)
		if derr != nil {
			return nil, server.NewInternalError(err, "failed inserting chunk to database and deleting chunk", "chunkHash", params.ChunkHash, "delete err", derr)
		}
		return nil, server.NewInternalError(err, "failed inserting chunk to database", "chunkHash", params.ChunkHash)
	}

	return &api.PutChunkCreated{}, nil
}

func (s *Service) GetChunk(ctx context.Context, params api.GetChunkParams) (api.GetChunkOK, error) {
	hash, err := hex.DecodeString(params.ChunkHash)
	if err != nil {
		return api.GetChunkOK{}, server.ErrInvalidHash
	}

	chunkExists, err := s.db.ReadQuery().CheckChunk(ctx, hash)
	if err != nil {
		return api.GetChunkOK{}, server.NewInternalError(err, "failed checking if chunk exists", "chunkHash", params.ChunkHash, "error", err)
	}
	if !chunkExists {
		return api.GetChunkOK{}, server.ErrChunkNotFound
	}

	reader, err := s.storage.Download(ctx, params.ChunkHash)
	if err != nil {
		return api.GetChunkOK{}, err
	}

	return api.GetChunkOK{Data: reader}, nil
}
