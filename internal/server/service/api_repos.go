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

func (s *Service) PutUserRepo(ctx context.Context, params api.PutUserRepoParams) (api.PutUserRepoRes, error) {
	if err := s.db.WriteQuery().CreateRepo(ctx, dbm.CreateRepoParams{
		RepoName: params.RepoName,
		UserID:   params.UserID,
	}); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return &api.PutUserRepoConflict{}, nil
			}
		}
		slog.Error("failed creating repo", "RepoName", params.RepoName, "UserID", params.UserID)
		return &api.PutUserRepoInternalServerError{}, server.ErrDatabase
	}
	return &api.PutUserRepoCreated{}, nil
}
