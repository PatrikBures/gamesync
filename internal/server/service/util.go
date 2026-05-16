package service

import (
	"context"
	"crypto/rand"
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

func hasPerm(ctx context.Context, perm permissions.Perm) error {
	perms, ok := ctx.Value(ckRolePerms).(permissions.Perms)
	if !ok {
		return ErrContext
	}
	if !slices.Contains(perms, perm) {
		return ErrNotAuthorized
	}
	return nil
}
