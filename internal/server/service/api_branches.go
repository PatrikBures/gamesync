package service

import (
	"context"
	"log/slog"

	api "go.pabu.dev/gamesync/internal/ogen"
	"go.pabu.dev/gamesync/internal/server"
	"go.pabu.dev/gamesync/internal/server/dbm"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Service) GetBranches(ctx context.Context, params api.GetBranchesParams) (api.Branches, error) {
	branches, err := s.db.ReadQuery().ListBranches(ctx, params.RepoName)
	if err != nil {
		return nil, server.NewInternalError(err, "failed listing branches",
			"userID", params.UserID,
			"repoName", params.RepoName,
		)
	}
	branchesReturn := make(api.Branches, 0, len(branches))
	for _, b := range branches {
		branchesReturn = append(branchesReturn, b.BranchName)
	}
	return branchesReturn, nil
}

func (s *Service) PutBranch(ctx context.Context, params api.PutBranchParams) error {
	repoID, err := repoFromCtx(ctx)
	if err != nil {
		return server.NewInternalError(err, "")
	}
	if err := s.db.WriteQuery().CreateBranch(ctx, dbm.CreateBranchParams{
		RepoID:     repoID,
		BranchName: params.BranchName,
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

func (s *Service) GetBranchHead(ctx context.Context, params api.GetBranchHeadParams) (result api.GetBranchHeadRes, err error) {
	branch, err := branchFromCtx(ctx)
	if err != nil {
		return
	}

	if !branch.HeadSnapshotID.Valid {
		return &api.GetBranchHeadNotFound{}, nil
	}
	snapshot, err := s.db.ReadQuery().GetSnapshot(ctx, dbm.GetSnapshotParams{
		SnapshotID: branch.HeadSnapshotID.Int64,
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
