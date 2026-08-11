package service

import (
	"context"
	api "go.pabu.dev/gamesync/internal/ogen"
	"go.pabu.dev/gamesync/internal/server"
	"go.pabu.dev/gamesync/internal/server/dbm"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)


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

func (s *Service) PutRepo(ctx context.Context, params api.PutRepoParams) error {
	_, err := s.db.WriteQuery().CreateRepo(ctx, dbm.CreateRepoParams{
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
