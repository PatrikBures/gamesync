package service

import (
	"context"
	"errors"
	"gamesync/internal/model"
	api "gamesync/internal/ogen"

	"gorm.io/gorm"
)

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
