package service

import (
	"context"
	"errors"
	"gamesync/internal/model"
	api "gamesync/internal/ogen"
	"gamesync/internal/server"

	"gorm.io/gorm"
)

func (s *Service) GetRoles(ctx context.Context) (api.GetRolesRes, error) {
	roles, err := s.q.Role.WithContext(ctx).Find()
	if err != nil {
		return &api.GetRolesInternalServerError{}, server.ErrDatabase
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
		return &api.PostRolesNotAcceptable{}, server.ErrMissingBody
	}
	role := model.Role{RoleName: req.Value.RoleName}
	if err := s.q.Role.WithContext(ctx).Create(&role); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &api.PostRolesConflict{}, server.ErrDuplicateKey
		}
		return &api.PostRolesInternalServerError{}, server.ErrDatabase
	}
	return &api.Role{RoleID: role.RoleID, RoleName: req.Value.RoleName}, nil
}

func (s *Service) PutRoleName(ctx context.Context, req api.OptRoleName, params api.PutRoleNameParams) (api.PutRoleNameRes, error) {
	if !req.Set {
		return &api.PutRoleNameNotAcceptable{}, server.ErrMissingBody
	}
	if _, err := s.q.Role.WithContext(ctx).
	Where(s.q.Role.RoleID.Eq(params.RoleID)).
	Update(s.q.Role.RoleName, req.Value.RoleName);
	err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &api.PutRoleNameConflict{}, server.ErrDuplicateKey
		}
	}
	return &api.PutRoleNameOK{}, nil
}
