package service

import (
	"context"
	"gamesync/internal/server"
	"gamesync/internal/server/permissions"
	"slices"
)


// server.ErrContext if there is a issue loading the context
//
// server.ErrNotAuthorized if user does not have perm 
//
// nil if the user has perm
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
