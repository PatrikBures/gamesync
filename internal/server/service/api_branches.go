package service

import (
	"context"
	api "gamesync/internal/ogen"
	"gamesync/internal/server"
	"gamesync/internal/server/dbm"
	"log/slog"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Service) GetUserRepoBranches(ctx context.Context, params api.GetUserRepoBranchesParams) (api.GetUserRepoBranchesRes, error) {
	branches, err := s.conn.Queries.ListBranches(ctx, params.RepoName)
	if err != nil {
		return &api.GetUserRepoBranchesInternalServerError{}, server.ErrDatabase
	}
	branchesReturn := make(api.Branches, 0, len(branches))
	for _, b := range branches {
		branchesReturn = append(branchesReturn, b.BranchName)
	}
	return &branchesReturn, nil
}

func (s *Service) PutUserRepoBranch(ctx context.Context, params api.PutUserRepoBranchParams) (api.PutUserRepoBranchRes, error) {
	if err := s.conn.Queries.CreateBranch(ctx, dbm.CreateBranchParams{
		RepoName: params.RepoName,
		BranchName: params.BranchName,
	}); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return &api.PutUserRepoBranchConflict{}, nil
			}
		}
		slog.Error("failed creating branch", "BranchName", params.BranchName, "RepoName", params.RepoName, "UserID", params.UserID)
		return &api.PutUserRepoBranchInternalServerError{}, server.ErrDatabase
	}

	return &api.PutUserRepoBranchCreated{}, nil
}
