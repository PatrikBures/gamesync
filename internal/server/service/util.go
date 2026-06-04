package service

import (
	"context"
	"errors"
	"gamesync/internal/server"
	"gamesync/internal/server/dbm"
	"gamesync/internal/server/permissions"
	"log/slog"
	"slices"
)

// returns either ErrContext, ErrNotAuthorized or nil
func CheckPerm(ctx context.Context, perm permissions.Perm) error {
	perms, ok := ctx.Value(ckRolePerms).(permissions.Perms)
	if !ok {
		return server.ErrContext
	}
	if !slices.Contains(perms, perm) {
		return server.ErrNotAuthorized
	}
	return nil
}



func isUserSelf(ctx context.Context, userID int64) (dbm.User, error) {
	currentUser, ok := ctx.Value(ckUser).(dbm.User)
	if !ok {
		slog.Error("mapping user context", "UserID", userID)
		return dbm.User{}, server.ErrContext
	}
	if userID != currentUser.UserID {
		return currentUser, server.ErrNotAuthorized
	}
	return currentUser, nil
}


// checks if the user is trying to access itself, 
// if it is it will return nil
//
// if it is trying to access a user id which is not itself, 
// it will check if it has the perm for that
//
// returns either ErrContext, ErrNotAuthorized or nil
func isUserSelfAndPerm(ctx context.Context, userID int64, perm permissions.Perm) error {
	_, errUser := isUserSelf(ctx, userID)
	if errors.Is(errUser, server.ErrContext) {
		return errUser
	} else if errUser == nil {
		return nil
	}

	errPerm := CheckPerm(ctx, perm)
	if errors.Is(errPerm, server.ErrContext) {
		return errPerm
	} else if errPerm != nil {
		return server.ErrNotAuthorized
	}

	return nil
}
