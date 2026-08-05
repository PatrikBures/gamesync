package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	api "go.pabu.dev/gamesync/internal/ogen"
	"go.pabu.dev/gamesync/internal/server"
	"go.pabu.dev/gamesync/internal/server/permissions"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func (s *Service) HandleBearerAuth(ctx context.Context, operationName api.OperationName, t api.BearerAuth) (context.Context, error) {
	tokenRaw := make([]byte, tokenLen)
	if _, err := base64.URLEncoding.Decode(tokenRaw, []byte(t.Token)); err != nil {
		return ctx, server.ErrBase64
	}
	if len(tokenRaw) != tokenLen {
		return ctx, server.ErrTokenLength
	}

	tokenHash := sha256.Sum256(tokenRaw)
	tokenHashSlice := tokenHash[:]

	user, err := s.db.ReadQuery().GetUserFromToken(ctx, tokenHashSlice)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ctx, server.ErrNotAuthorized
		}
		return ctx, server.NewInternalError(err, "failed authorizing user", "userID", user.UserID)
	}
	userRolePerms, err := s.db.ReadQuery().ListRolePerms(ctx, user.RoleID)
	if err != nil {
		return ctx, server.NewInternalError(err, "failed listing role perms", "roleID", user.RoleID)
	}
	perms := make(permissions.Perms, 0, len(userRolePerms))
	for _, perm := range permissions.AllPerms {
		perm32 := int32(perm)
		for _, rolePerm := range userRolePerms {
			if rolePerm.PermID == perm32 {
				perms = append(perms, perm)
				break
			}
		}
	}
	if len(perms) != len(userRolePerms) {
		slog.Error("Role has unmapped permID(s)", "RoleID", user.RoleID)
	}

	ctx = context.WithValue(ctx, CkUser, user)
	ctx = context.WithValue(ctx, CkRolePerms, perms)
	return ctx, nil
}
