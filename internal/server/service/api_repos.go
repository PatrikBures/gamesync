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

const defaultBranchName = "main"

func (s *Service) GetUserRepos(ctx context.Context, params api.GetUserReposParams) (api.GetUserReposRes, error) {
	repos, err := s.db.ReadQuery().ListRepos(ctx, params.UserID)
	if err != nil {
		return &api.GetUserReposInternalServerError{}, server.ErrDatabase
	}
	reposReturn := make(api.Repos, 0, len(repos))
	for _, r := range repos {
		reposReturn = append(reposReturn, r.RepoName)
	}
	return &reposReturn, nil
}

func (s *Service) PutUserRepo(ctx context.Context, params api.PutUserRepoParams) (result api.PutUserRepoRes, err error) {
	qtx, tx, err := s.db.BeginTX(ctx)
	if err != nil {
		return &api.PutUserRepoInternalServerError{}, server.ErrDatabase
	}
	defer func() {
		if recover() != nil || err != nil {
			if e := tx.Rollback(ctx); e != nil {
				slog.Error("failed rollback", "error", e)
			}
		}
	}()
	repoID, err := qtx.CreateRepo(ctx, dbm.CreateRepoParams{
		RepoName: params.RepoName,
		UserID:   params.UserID,
	})
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return &api.PutUserRepoConflict{}, nil
			}
		}
		slog.Error("failed creating repo", "RepoName", params.RepoName, "UserID", params.UserID)
		return &api.PutUserRepoInternalServerError{}, server.ErrDatabase
	}

	if err = qtx.CreateBranch(ctx, dbm.CreateBranchParams{
		RepoID: repoID,
		BranchName: defaultBranchName,
	}); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				slog.Error("default branch somehow exists on new repo", "UserID", params.UserID, "RepoName", params.RepoName)
				return &api.PutUserRepoConflict{}, nil
			}
		}
		slog.Error("failed creating branch", "RepoID", repoID, "BranchName", defaultBranchName)
		return &api.PutUserRepoInternalServerError{}, nil
	}

	if err = tx.Commit(ctx); err != nil {
		slog.Error("committing tx", "error", err)
		return &api.PutUserRepoInternalServerError{}, server.ErrDatabase
	}

	return &api.PutUserRepoCreated{}, nil
}
