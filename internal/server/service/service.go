package service

import (
	"context"
	"errors"
	"fmt"
	"gamesync/internal/model"
	api "gamesync/internal/ogen"
	"gamesync/internal/query"

	"gorm.io/gorm"
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
	return &api.UserNewReturn{
		UserId: user.UserID,
		UserName: user.UserName,
		RoleId: user.RoleID,
		Token: "asdf",
	}, nil
}

