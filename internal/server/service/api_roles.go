package service

import (
	"context"
	api "gamesync/internal/ogen"
	"gamesync/internal/server"
	"gamesync/internal/server/dbm"
	"log/slog"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Service) GetRoles(ctx context.Context) (api.GetRolesRes, error) {
	roles, err := s.db.ReadQuery().ListRoles(ctx)
	if err != nil {
		slog.Error("listing roles", "error", err)
		return &api.GetRolesInternalServerError{}, server.ErrDatabase
	}
	rolesReturn := make(api.GetRolesOKApplicationJSON, 0, len(roles))
	for _, role := range roles {
		rolesReturn = append(rolesReturn, api.Role{
			RoleID:   role.RoleID,
			RoleName: role.RoleName,
		})
	}
	return &rolesReturn, nil
}

func (s *Service) PostRoles(ctx context.Context, req api.OptRoleName) (api.PostRolesRes, error) {
	if !req.Set {
		return &api.PostRolesNotAcceptable{}, nil
	}
	roleID, err := s.db.WriteQuery().InsertRole(ctx, req.Value.RoleName)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return &api.PostRolesConflict{}, nil
			}
		}
		return &api.PostRolesInternalServerError{}, server.ErrDatabase
	}
	return &api.Role{RoleID: roleID, RoleName: req.Value.RoleName}, nil
}

func (s *Service) PutRoleName(ctx context.Context, req api.OptRoleName, params api.PutRoleNameParams) (api.PutRoleNameRes, error) {
	if !req.Set {
		return &api.PutRoleNameNotAcceptable{}, nil
	}
	if err := s.db.WriteQuery().UpdateRoleName(ctx, dbm.UpdateRoleNameParams{
		RoleID:   params.RoleID,
		RoleName: req.Value.RoleName,
	}); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return &api.PutRoleNameConflict{}, nil
			case pgerrcode.ForeignKeyViolation:
				return &api.PutRoleNameNotFound{}, nil
			}
		}
		return &api.PutRoleNameInternalServerError{}, server.ErrDatabase
	}
	return &api.PutRoleNameOK{}, nil
}
