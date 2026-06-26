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
