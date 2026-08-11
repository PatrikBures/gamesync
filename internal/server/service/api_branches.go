package service

import (
	"context"
	"log/slog"

	api "go.pabu.dev/gamesync/internal/ogen"
	"go.pabu.dev/gamesync/internal/server"
	"go.pabu.dev/gamesync/internal/server/dbm"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) GetBranches(ctx context.Context, params api.GetBranchesParams) ([]api.Branch, error) {
	repoID, err := repoFromCtx(ctx)
	if err != nil { return nil, err }

	branches, err := s.db.ReadQuery().ListBranches(ctx, dbm.ListBranchesParams{
		RepoID: repoID,
		BranchName: pgtype.Text{
			String: params.BranchName.Value,
			Valid: params.BranchName.Set,
		},
	})
	if err != nil {
		return nil, server.NewInternalError(err, "failed listing branches",
			"userID", params.UserID,
			"repoID", repoID,
		)
	}

	branchesReturn := make([]api.Branch, 0, len(branches))
	for _, b := range branches {
		branchesReturn = append(branchesReturn, api.Branch{
			RepoID: b.RepoID,
			BranchID: b.BranchID,
			HeadSnapshotID: b.HeadSnapshotID,
			ParentSnapshotID: api.OptNilInt64{
				Value: b.ParentSnapshotID.Int64,
				Null: b.ParentSnapshotID.Valid,
				Set: true,
			},
			BranchName: b.BranchName,
		})
	}
	return branchesReturn, nil
}

func (s *Service) PutBranch(ctx context.Context, req *api.SnapshotID, params api.PutBranchParams) error {
	repoID, err := repoFromCtx(ctx)
	if err != nil {
		return server.NewInternalError(err, "")
	}
	if _, err := s.db.WriteQuery().CreateBranch(ctx, dbm.CreateBranchParams{
		RepoID:     repoID,
		BranchName: params.BranchName,
		HeadSnapshotID: req.SnapshotID,
	}); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return server.ErrBranchNameConflict
			case pgerrcode.ForeignKeyViolation:
				return server.ErrRepoNotFound
			}
		}
		return server.NewInternalError(err, "failed creating branch",
			"userID", params.UserID,
			"repoName", params.RepoName,
			"branchName", params.BranchName,
		)
	}

	return nil
}

func (s *Service) GetBranchHead(ctx context.Context, params api.GetBranchHeadParams) (*api.Snapshot, error) {
	branch, err := branchFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	snapshot, err := s.db.ReadQuery().GetSnapshot(ctx, dbm.GetSnapshotParams{
		SnapshotID: branch.HeadSnapshotID,
		RepoID: branch.RepoID,
	})
	if err != nil {
		return nil, server.NewInternalError(err, "failed getting snapshot while getting branch head", "snapshotID", branch.HeadSnapshotID)
	}

	return &api.Snapshot{
		SnapshotID: snapshot.SnapshotID,
		ParentSnapshotID: api.NilInt64{
			Value: snapshot.ParentSnapshotID.Int64,
			Null:  !snapshot.ParentSnapshotID.Valid,
		},
	}, nil
}

func (s *Service) DeleteBranch(ctx context.Context, params api.DeleteBranchParams) (err error) {
	branch, err := branchFromCtx(ctx)
	if err != nil {
		return
	}

	slog.Info("deleting branch", "repoID", branch.RepoID, "branchID", branch.BranchID)

	if err := s.db.WriteQuery().DeleteBranch(ctx, dbm.DeleteBranchParams{
		RepoID:   branch.RepoID,
		BranchID: branch.BranchID,
	}); err != nil {
		return server.NewInternalError(err, "deleting branch", "repoID", branch.RepoID, "branchID", branch.BranchID)
	}

	return nil
}
