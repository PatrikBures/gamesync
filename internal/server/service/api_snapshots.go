package service

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"go.pabu.dev/gamesync/internal/dbx"
	api "go.pabu.dev/gamesync/internal/ogen"
	"go.pabu.dev/gamesync/internal/server"
	"go.pabu.dev/gamesync/internal/server/dbm"
	"log/slog"
	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) PostSnapshot(ctx context.Context, req *api.Files, params api.PostSnapshotParams) (result api.PostSnapshotRes, err error) {

	files, chunkCount, err := reqFilesStruct(req.Files)
	if err != nil {
		return nil, err
	}

	allChunkHashes := make([][]byte, 0, chunkCount)
	for _, f := range files {
		allChunkHashes = append(allChunkHashes, f.chunkHashes...)
	}

	// sorts and removes duplicate chunks
	slices.SortFunc(allChunkHashes, bytes.Compare)
	allChunkHashes = slices.CompactFunc(allChunkHashes, func(a []byte, b []byte) bool {
		return bytes.Equal(a, b)
	})

	existingChunks, err := s.db.ReadQuery().GetChunkHashes(ctx, allChunkHashes)
	if err != nil {
		return nil, server.NewInternalError(err, "getting existing chunks")
	}

	allChunkHashes = deleteExistingChunks(allChunkHashes, existingChunks)
	if len(allChunkHashes) > 0 {
		return &api.PostSnapshotFailedDependency{ChunkHashes: bytesToHex(allChunkHashes)}, nil
	}

	newSnapshot, err := createSnapshot(s.db, files, ctx)
	if err != nil {
		return nil, err
	}

	return newSnapshot, nil
}

type file struct {
	path        string
	hash        []byte
	chunkHashes [][]byte
}

func createSnapshot(db *dbx.DB, files []file, ctx context.Context) (*api.Snapshot, error) {
	branch, err := branchFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	qtx, tx, err := db.BeginTX(ctx)
	if err != nil {
		return nil, server.NewInternalError(err, "failed starting tx")
	}
	defer func() {
		if recover() != nil || err != nil {
			if e := tx.Rollback(ctx); e != nil {
				slog.Error("failed rollback", "error", e)
			}
		}
	}()

	newSnapshot, err := qtx.CreateSnapshot(ctx, dbm.CreateSnapshotParams{
		RepoID:   branch.RepoID,
		BranchID: branch.BranchID,
	})
	if err != nil {
		return nil, server.NewInternalError(err, "failed creating new snapshot row", "repoID", branch.RepoID, "branchID", branch.BranchID)
	}

	connectSnapshotParams := make([]dbm.ConnectSnapshotWithFilesParams, 0, len(files))
	fileHashes := make([][]byte, 0, len(files))
	for _, f := range files {
		fileHashes = append(fileHashes, f.hash)
	}
	fileHashes, err = qtx.ListFileHashes(ctx, fileHashes)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			fileHashes = fileHashes[:0]
		} else {
			return nil, server.NewInternalError(err, "failed listing existing file hashes")
		}
	}

	for _, f := range files {
		connectSnapshotParams = append(connectSnapshotParams, dbm.ConnectSnapshotWithFilesParams{
			SnapshotID: newSnapshot.SnapshotID,
			FileHash:   f.hash,
			FilePath:   f.path,
		})

		if slices.ContainsFunc(fileHashes, func(item []byte) bool {
			return bytes.Equal(item, f.hash)
		}) {
			continue
		}

		if err := qtx.CreateFile(ctx, dbm.CreateFileParams{
			FileHash:  f.hash,
			ChunkHash: f.chunkHashes,
		}); err != nil {
			return nil, server.NewInternalError(err, "failed inserting file", "fileHash", hex.EncodeToString(f.hash))
		}

		if _, err := qtx.ConnectFileWithChunks(ctx, connectFileChunksParams(f.hash, f.chunkHashes)); err != nil {
			return nil, server.NewInternalError(err, "failed connecting file with chunks", "fileHash", hex.EncodeToString(f.hash))
		}
		fileHashes = append(fileHashes, f.hash)
	}

	if _, err := qtx.ConnectSnapshotWithFiles(ctx, connectSnapshotParams); err != nil {
		return nil, server.NewInternalError(err, "failed connecting snapshot with files", "snapshotID", newSnapshot.SnapshotID)
	}

	if err := qtx.UpdateBranchHead(ctx, dbm.UpdateBranchHeadParams{
		BranchID:       branch.BranchID,
		HeadSnapshotID: pgtype.Int8{Int64: newSnapshot.SnapshotID, Valid: true},
	}); err != nil {
		return nil, server.NewInternalError(err, "failed updating branch head", "branchID", branch.BranchID, "newSnapshotID", newSnapshot.SnapshotID)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, server.NewInternalError(err, "failed committing creation of snapshot", "repoID", branch.RepoID, "branchID", branch.BranchID)
	}

	return &api.Snapshot{
		ParentSnapshotID: api.NilInt64{Value: newSnapshot.ParentSnapshotID.Int64, Null: !newSnapshot.ParentSnapshotID.Valid},
		SnapshotID: newSnapshot.SnapshotID,
	}, nil
}

func connectFileChunksParams(fileHash []byte, chunkHashes [][]byte) []dbm.ConnectFileWithChunksParams {
	s := make([]dbm.ConnectFileWithChunksParams, 0, len(chunkHashes))
	for i := int32(0); i < int32(len(chunkHashes)); i++ {
		s = append(s, dbm.ConnectFileWithChunksParams{
			FileHash:   fileHash,
			ChunkHash:  chunkHashes[i],
			ChunkOrder: i + 1,
		})
	}
	return s
}

func reqFilesStruct(reqFiles []api.File) ([]file, int, error) {
	files := make([]file, 0, len(reqFiles))
	chunkCount := 0

	for _, f := range reqFiles {
		nf := file{
			path:        f.Path,
			chunkHashes: make([][]byte, 0, len(f.ChunkHashes)),
		}
		h, err := hex.DecodeString(f.Hash)
		if err != nil {
			return nil, 0, server.ErrInvalidHash
		}
		nf.hash = h
		for _, c := range f.ChunkHashes {
			h, err := hex.DecodeString(c)
			if err != nil {
				return nil, 0, server.ErrInvalidHash
			}
			nf.chunkHashes = append(nf.chunkHashes, h)
			chunkCount++
		}
		files = append(files, nf)
	}
	return files, chunkCount, nil
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

func (s *Service) GetSnapshot(ctx context.Context, params api.GetSnapshotParams) (result *api.SnapshotFiles, err error) {
	snapshot, err := snapshotFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	qtx, tx, err := s.db.BeginReadTX(ctx)
	if err != nil {
		return nil, server.NewInternalError(err, "failed starting tx")
	}
	defer func() {
		if recover() != nil || err != nil {
			if e := tx.Rollback(ctx); e != nil {
				slog.Error("failed rollback", "error", e)
			}
		}
	}()

	snapshotFiles, err := qtx.GetSnapshotFiles(ctx, snapshot.SnapshotID)
	if err != nil {
		return nil, server.NewInternalError(err, "failed listing snapshot files")
	}

	result = &api.SnapshotFiles{
		Files: make([]api.File, 0, len(snapshotFiles)),
	}

	result.ParentSnapshotID.Value = snapshot.ParentSnapshotID.Int64
	result.ParentSnapshotID.Null = !snapshot.ParentSnapshotID.Valid

	for _, sf := range snapshotFiles {
		chunkHashes, err := qtx.GetFileChunkHashes(ctx, sf.FileHash)
		if err != nil {
			return nil, server.NewInternalError(err, "failed listing file chunk hashes", "snapshotID", snapshot.SnapshotID, "fileHash", sf.FileHash)
		}

		ChunkHashesHex := make([]string, 0, len(chunkHashes))
		for _, ch := range chunkHashes {
			ChunkHashesHex = append(ChunkHashesHex, hex.EncodeToString(ch))
		}

		result.Files = append(result.Files, api.File{
			ChunkHashes: ChunkHashesHex,
			Hash:        hex.EncodeToString(sf.FileHash),
			Path:        sf.FilePath,
		})
	}

	return result, nil
}
