package middlewares

import (
	"context"
	"errors"
	"gamesync/internal/dbx"
	"gamesync/internal/server"
	"gamesync/internal/server/dbm"
	"gamesync/internal/server/service"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/ogen-go/ogen/middleware"
)

func RepoBranch(db *dbx.DB) middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {

		// skip checking if repo exists when creating one 
		if req.OperationName == "PutUserRepo" {
			return next(req)
		}

		targetRepoName, ok := req.Params.Path("repoName")
		if !ok {
			return next(req)
		}


		var repoName string
		switch i := targetRepoName.(type) {
		case string:
			repoName = i
		default:
			slog.Error("repoName is not string", "path", req.Raw.URL.Path)
			return middleware.Response{}, server.ErrAuth
		}

		user := req.Context.Value(service.CkUser).(dbm.User)
		repo, err := db.ReadQuery().GetRepoWithName(req.Context, dbm.GetRepoWithNameParams{
			UserID: user.UserID,
			RepoName: repoName,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return middleware.Response{}, server.ErrRepoNotFound
			}
			slog.Error("getting repo", "userID", user.UserID, "repoName", repoName, "error", err)
			return middleware.Response{}, server.ErrDatabase
		}
		context.WithValue(req.Context, service.CkRepoID, repo.RepoID)



		// skip checking if branch exists when creating one 
		if req.OperationName == "PutUserRepoBranch" {
			return next(req)
		}
		targetBranchName, ok := req.Params.Path("branchName")
		if !ok {
			return next(req)
		}

		var branchName string
		switch i := targetBranchName.(type) {
		case string:
			branchName = i
		default:
			slog.Error("branchName is not string", "path", req.Raw.URL.Path)
			return middleware.Response{}, server.ErrAuth
		}

		branch, err := db.ReadQuery().GetBranchWithName(req.Context, dbm.GetBranchWithNameParams{
			RepoID: repo.RepoID,
			BranchName: branchName,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return middleware.Response{}, server.ErrBranchNotFound
			}
			slog.Error("getting branch", "repoID", repo.RepoID, "branchName", branchName, "error", err)
			return middleware.Response{}, server.ErrDatabase
		}
		context.WithValue(req.Context, service.CkBranchID, branch.BranchID)
		
		return next(req)
	}
}
