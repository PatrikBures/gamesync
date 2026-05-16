package service

import (
	"context"
	"errors"
	"gamesync/internal/model"
	api "gamesync/internal/ogen"
	"gamesync/internal/server/permissions"

	"gorm.io/gorm"
)

func (s *Service) GetRoles(ctx context.Context) (api.GetRolesRes, error) {
	if err := hasPerm(ctx, permissions.PermRolesGet); err != nil {
		return &api.GetRolesUnauthorized{}, err
	}

	roles, err := s.q.Role.WithContext(ctx).Find()
	if err != nil {
		return &api.GetRolesInternalServerError{}, ErrDatabase
	}
	rolesReturn := make(api.GetRolesOKApplicationJSON, 0, len(roles))
	for _, role := range roles {
		rolesReturn = append(rolesReturn, api.Role{
			RoleID: role.RoleID,
			RoleName: role.RoleName,
		})
	}
	return &rolesReturn, nil
}

func (s *Service) PostRoles(ctx context.Context, req api.OptRoleName) (api.PostRolesRes, error) {
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
	return &api.Role{RoleID: role.RoleID, RoleName: req.Value.RoleName}, nil
}
