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

func (s *Service) GetUserRepos(ctx context.Context, params api.GetUserReposParams) (api.Repos, error) {
	repos, err := s.db.ReadQuery().ListRepos(ctx, params.UserID)
	if err != nil {
		return nil, server.NewInternalError(err, "failed listing repos", "userID", params.UserID)
	}
	reposReturn := make(api.Repos, 0, len(repos))
	for _, r := range repos {
		reposReturn = append(reposReturn, r.RepoName)
	}
	return reposReturn, nil
}

func (s *Service) PutUserRepo(ctx context.Context, params api.PutUserRepoParams) (err error) {
	qtx, tx, err := s.db.BeginTX(ctx)
	if err != nil {
		return server.NewInternalError(err, "failed starting tx")
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
				return server.ErrRepoNameConflict
			}
		}
		return server.NewInternalError(err, "failed creating repo", "userID", params.UserID, "repoName", params.RepoName)
	}

	if err = qtx.CreateBranch(ctx, dbm.CreateBranchParams{
		RepoID: repoID,
		BranchName: defaultBranchName,
	}); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return server.NewInternalError(err, "default branch already exists on new branch", "userID", params.UserID, "repoName", params.RepoName)
			}
		}
		return server.NewInternalError(err, "failed creating branch", "repoID", repoID, "branchName", defaultBranchName)
	}

	if err = tx.Commit(ctx); err != nil {
		return server.NewInternalError(err, "failed committing tx")
	}

	return nil
}
