package service

import (
	"context"
	api "go.pabu.dev/gamesync/internal/ogen"
	"go.pabu.dev/gamesync/internal/server"
	"go.pabu.dev/gamesync/internal/server/dbm"
	"log/slog"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

const defaultBranchName = "main"

func (s *Service) GetRepos(ctx context.Context, params api.GetReposParams) (api.Repos, error) {
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

func (s *Service) PutRepo(ctx context.Context, params api.PutRepoParams) (err error) {
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
		RepoID:     repoID,
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

func (s *Service) DeleteRepo(ctx context.Context, params api.DeleteRepoParams) error {
	repoID, err := repoFromCtx(ctx)
	if err != nil {
		return err
	}

	if err := s.db.WriteQuery().DeleteRepo(ctx, repoID); err != nil {
		return server.NewInternalError(err, "deleting repo", "repoID", repoID)
	}

	return nil
}
