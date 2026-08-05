package service

import (
	"context"
	"go.pabu.dev/gamesync/internal/server"
	"go.pabu.dev/gamesync/internal/server/dbm"
)

// err is ServerError
func branchFromCtx(ctx context.Context) (branch dbm.Branch, err error) {
	branch, ok := ctx.Value(CkBranch).(dbm.Branch)
	if !ok {
		err = server.NewInternalError(nil, "branch is not mapped", "branchID", branch.BranchID)
	}
	return
}

// err is ServerError
func repoFromCtx(ctx context.Context) (repoID int64, err error) {
	repoID, ok := ctx.Value(CkRepoID).(int64)
	if !ok {
		err = server.NewInternalError(nil, "repoID is not mapped", "repoID", repoID)
	}
	return
}

// err is ServerError
func snapshotFromCtx(ctx context.Context) (snapshot dbm.Snapshot, err error) {
	snapshot, ok := ctx.Value(CkSnapshot).(dbm.Snapshot)
	if !ok {
		err = server.NewInternalError(nil, "snapshot is not mapped", "snapshotID", snapshot.SnapshotID)
	}
	return
}
