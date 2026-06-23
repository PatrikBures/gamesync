package service

import (
	"bytes"
	"context"
	"encoding/hex"
	api "gamesync/internal/ogen"
	"gamesync/internal/server"
	"log/slog"
	"slices"
)

func (s *Service) PostUserRepoBranchSnapshot(ctx context.Context, req *api.SnapshotNew, params api.PostUserRepoBranchSnapshotParams) (api.PostUserRepoBranchSnapshotRes, error) {

	totalChunkAmount := 0
	for _, f := range req.Files {
		totalChunkAmount += len(f.ChunkHashes)
	}

	allChunkHashes := make([][]byte, 0, totalChunkAmount)

	for _, f := range req.Files {
		for _, c := range f.ChunkHashes {
			s, err := hex.DecodeString(c)
			if err != nil {
				return &api.PostUserRepoBranchSnapshotUnprocessableEntity{}, nil
			}
			allChunkHashes = append(allChunkHashes, s)
		}
	}

	// sorts and removes duplicate chunks
	slices.SortFunc(allChunkHashes, bytes.Compare)
	allChunkHashes = slices.CompactFunc(allChunkHashes, func(a []byte, b []byte) bool {
		return bytes.Compare(a, b) == 0
	})

	existingChunks, err := s.db.ReadQuery().GetChunkHashes(ctx, allChunkHashes)
	if err != nil {
		slog.Error("getting existing chunks", "error", err)
		return &api.PostUserRepoBranchSnapshotInternalServerError{}, server.ErrDatabase
	}

	allChunkHashes = deleteExistingChunks(allChunkHashes, existingChunks)
	if len(allChunkHashes) > 0 {
		return &api.PostUserRepoBranchSnapshotFailedDependency{ChunkHashes: bytesToHex(allChunkHashes)}, nil
	}

	return &api.PostUserRepoBranchSnapshotCreated{}, nil
}



// removes every []byte from a if it exists in b
//
// a and b need to be sorted and have no duplicates
func deleteExistingChunks(a [][]byte, b [][]byte) [][]byte {
	if len(b) == 0 {
		return a
	}
	eci := 0
	a = slices.DeleteFunc(a, func(a []byte) bool {
		if eci == len(b) {
			return false
		}
		if !bytes.Equal(a, b[eci]) {
			return false
		}
		eci++
		return true
	})
	return a
}

func bytesToHex(bytes [][]byte) []string {
	s := make([]string, 0, len(bytes))
	for _, b := range bytes {
		s = append(s, hex.EncodeToString(b))
	}
	return s
}
