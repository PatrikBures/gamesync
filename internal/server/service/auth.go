package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	api "gamesync/internal/ogen"
)

func (s *Service) HandleBearerAuth(ctx context.Context, operationName api.OperationName, t api.BearerAuth) (context.Context, error) {
	tokenRaw := make([]byte, tokenLen)
	if _, err := base64.URLEncoding.Decode(tokenRaw, []byte(t.Token)); err != nil {
		return ctx, ErrToken
	}
	if len(tokenRaw) != tokenLen {
		return ctx, ErrToken
	}

	tokenHash := sha256.Sum256(tokenRaw)
	tokenHashSlice := tokenHash[:]

	tokenRecord, err := s.q.Token.WithContext(ctx).Where(s.q.Token.TokenHash.Eq(HashBytes(tokenHashSlice))).First()
	if err != nil {
		return ctx, ErrNotAuthorized
	}
	user, err := s.q.User.WithContext(ctx).Where(s.q.User.UserID.Eq(tokenRecord.UserID)).First()
	if err != nil {
		return ctx, ErrNotAuthorized
	}
	ctx = context.WithValue(ctx, userContextKey, user)
	return ctx, nil
}
