package middlewares

import (
	"context"
	"errors"
	"gamesync/internal/dbx"
	"gamesync/internal/server"
	"gamesync/internal/server/dbm"
	"gamesync/internal/server/service"

	"github.com/jackc/pgx/v5"
	"github.com/ogen-go/ogen/middleware"
)

// loads all values in path into context
func LoadPathData(db *dbx.DB) middleware.Middleware {
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
			return middleware.Response{}, server.NewInternalError(nil, "repoName is not string", "path", req.Raw.URL.Path)
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
			return middleware.Response{}, server.NewInternalError(err, "failed getting repo",
				"userID", user.UserID,
				"repoName", "repoName",
			)
		}
		req.Context = context.WithValue(req.Context, service.CkRepoID, repo.RepoID)



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
			return middleware.Response{}, server.NewInternalError(nil, "branchName is not string", "path", req.Raw.URL.Path)
		}

		branch, err := db.ReadQuery().GetBranchWithName(req.Context, dbm.GetBranchWithNameParams{
			RepoID: repo.RepoID,
			BranchName: branchName,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return middleware.Response{}, server.ErrBranchNotFound
			}
			return middleware.Response{}, server.NewInternalError(err, "failed getting branch",
				"userID", user.UserID,
				"repoID", repo.RepoID,
				"branchName", branchName,
			)
		}
		req.Context = context.WithValue(req.Context, service.CkBranch, branch)




		targetSnapshotID, ok := req.Params.Path("snapshotID")
		if !ok {
			return next(req)
		}
		var snapshotID int64
		switch i := targetSnapshotID.(type) {
		case int64:
			snapshotID = i
		default:
			return middleware.Response{}, server.NewInternalError(nil, "snapshotID is not int64", "path", req.Raw.URL.Path)
		}

		snapshot, err := db.ReadQuery().GetSnapshot(req.Context, snapshotID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return middleware.Response{}, server.ErrSnapshotNotFound
			}
			return middleware.Response{}, server.NewInternalError(err, "failed getting snapshot",
				"userID", user.UserID,
				"repoID", repo.RepoID,
				"branchID", branch.BranchID,
				"snapshotID", snapshotID,
			)
		}
		req.Context = context.WithValue(req.Context, service.CkSnapshot, snapshot)
		
		return next(req)
	}
}
