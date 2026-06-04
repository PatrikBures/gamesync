package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	api "gamesync/internal/ogen"
	"gamesync/internal/server"
	"gamesync/internal/server/permissions"
	"log/slog"
)

func (s *Service) HandleBearerAuth(ctx context.Context, operationName api.OperationName, t api.BearerAuth) (context.Context, error) {
	tokenRaw := make([]byte, tokenLen)
	if _, err := base64.URLEncoding.Decode(tokenRaw, []byte(t.Token)); err != nil {
		return ctx, server.ErrToken
	}
	if len(tokenRaw) != tokenLen {
		return ctx, server.ErrToken
	}

	tokenHash := sha256.Sum256(tokenRaw)
	tokenHashSlice := tokenHash[:]

	userID, err := s.conn.Queries.GetUserIdFromToken(ctx, tokenHashSlice)
	if err != nil {
		slog.Warn("authorizing user, getting token", "userID", userID, "error", err)
		return ctx, server.ErrNotAuthorized
	}
	user, err := s.conn.Queries.GetUser(ctx, userID)
	if err != nil {
		slog.Warn("authorizing user, getting user", "userID", userID, "error", err)
		return ctx, server.ErrNotAuthorized
	}
	userRolePerms, err := s.conn.Queries.ListRolePerms(ctx, user.RoleID)
	if err != nil {
		return ctx, server.ErrDatabase
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
		slog.Error("Role has invalid RoleID", "RoleID", user.RoleID)
	}

	ctx = context.WithValue(ctx, ckUser, user)
	ctx = context.WithValue(ctx, ckRolePerms, perms)
	return ctx, nil
}
