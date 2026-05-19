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

	tokenRecord, err := s.q.Token.WithContext(ctx).Where(s.q.Token.TokenHash.Eq(HashBytes(tokenHashSlice))).First()
	if err != nil {
		return ctx, server.ErrNotAuthorized
	}
	user, err := s.q.User.WithContext(ctx).Where(s.q.User.UserID.Eq(tokenRecord.UserID)).First()
	if err != nil {
		return ctx, server.ErrNotAuthorized
	}
	userRolePerms, err := s.q.RolePermission.WithContext(ctx).Where(s.q.RolePermission.RoleID.Eq(user.RoleID)).Find()
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

	// TODO: centralized authorization using middleware

	ctx = context.WithValue(ctx, ckUser, user)
	ctx = context.WithValue(ctx, ckRolePerms, perms)
	return ctx, nil
}
