package service

import (
	"context"
	"crypto/rand"
	"gamesync/internal/server"
	"gamesync/internal/server/permissions"
	"slices"
)

func generateToken() ([]byte, error) {
	b := make([]byte, 33)
	_, err := rand.Read(b)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

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
