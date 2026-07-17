package service

import (
	"context"
	api "gamesync/internal/ogen"
	"gamesync/internal/server"
	"gamesync/internal/server/dbm"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Service) GetUserRepoBranches(ctx context.Context, params api.GetUserRepoBranchesParams) (api.Branches, error) {
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

func (s *Service) PutUserRepoBranch(ctx context.Context, params api.PutUserRepoBranchParams) error {
	if err := s.db.WriteQuery().CreateBranchWithRepoName(ctx, dbm.CreateBranchWithRepoNameParams{
		UserID:     params.UserID,
		RepoName:   params.RepoName,
		BranchName: params.BranchName,
	}); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return server.ErrBranchNameConflict
			case pgerrcode.ForeignKeyViolation:
				return server.ErrBranchNotFound
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
	if err != nil { return }

	if !branch.HeadSnapshotID.Valid {
		return &api.GetBranchHeadNotFound{}, nil
	}
	snapshot, err := s.db.ReadQuery().GetSnapshot(ctx, branch.HeadSnapshotID.Int64)
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
