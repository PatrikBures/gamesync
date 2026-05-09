package service

import (
	"context"
	"database/sql/driver"
	"errors"
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
func (s *Service) GetHealth(ctx context.Context) error {
	return nil
}
func (s *Service) GetRoles(ctx context.Context) (api.GetRolesRes, error) {
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


