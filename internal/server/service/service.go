package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"
	"gamesync/internal/model"
	api "gamesync/internal/ogen"
	"gamesync/internal/query"

	"gorm.io/gorm"
)

const tokenLen = 33

// This type is requred as the gorm driver does not support inserting bytes
type HashBytes []byte
func (h HashBytes) Value() (driver.Value, error) {
	return []byte(h), nil
}

type contextKey int
const (
	userContextKey contextKey = iota
)

type Service struct {
	q *query.Query
}
func NewService(query *query.Query) *Service {
	return &Service{
		q: query,
	}
}
func (s *Service) HandleBearerAuth(ctx context.Context, operationName api.OperationName, t api.BearerAuth) (context.Context, error) {
	tokenRaw := make([]byte, tokenLen, tokenLen)
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
	return nil, nil
}
func (s *Service) GetHealth(ctx context.Context) error {
	return nil
}
func (s *Service) GetRoles(ctx context.Context) (api.GetRolesRes, error) {
	return nil, nil
}
func (s *Service) GetUserID(ctx context.Context, params api.GetUserIDParams) error {
	return nil
}
func (s *Service) GetUsers(ctx context.Context) (api.GetUsersRes, error) {
	return nil, nil
}
func (s *Service) PostRoles(ctx context.Context, req api.OptRoleNew) (api.PostRolesRes, error) {
	if !req.Set {
		return &api.PostRolesNotAcceptable{}, ErrMissingBody
	}
	role := model.Role{RoleName: req.Value.RoleName}
	if err := s.q.Role.WithContext(ctx).Create(&role); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &api.PostRolesConflict{}, ErrDuplicateKey
		}
		return &api.PostRolesInternalServerError{}, ErrDatabase
	}
	return &api.Role{RoleId: role.RoleID, RoleName: req.Value.RoleName}, nil
}
func (s *Service) PostUsers(ctx context.Context, req api.OptUserNew) (api.PostUsersRes, error) {
	if !req.Set {
		return &api.PostUsersNotAcceptable{}, ErrMissingBody
	}
	user := model.User{
		UserName: req.Value.UserName,
		RoleID: req.Value.RoleId,
	}
	if err := s.q.User.WithContext(ctx).Create(&user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &api.PostUsersConflict{}, ErrDuplicateKey
		} else if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return &api.PostUsersNotAcceptable{}, fmt.Errorf("invalid role")
		}
		return &api.PostUsersInternalServerError{}, ErrDatabase
	}


	token, err := generateToken()
	if err != nil {
		return &api.PostUsersInternalServerError{}, ErrToken
	}
	token64 := base64.URLEncoding.EncodeToString(token)
	tokenHash := sha256.Sum256(token)

	t := model.Token{
		UserID: user.UserID,
		TokenHash: tokenHash[:],
	}
	if err := s.q.Token.WithContext(ctx).Create(&t); err != nil {
		return &api.PostUsersInternalServerError{}, ErrDatabase
	}

	return &api.UserNewReturn{
		UserId: user.UserID,
		UserName: user.UserName,
		RoleId: user.RoleID,
		Token: token64,
	}, nil
}


func generateToken() ([]byte, error) {
	b := make([]byte, 33)
	_, err := rand.Read(b)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}
