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
	roles, err := s.conn.Queries.ListRoles(ctx)
	if err != nil {
		slog.Error("listing roles", "error", err)
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
	roleID, err := s.conn.Queries.InsertRole(ctx, req.Value.RoleName)
	if err != nil {
		var pgErr *pgconn.PgError
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return &api.PostRolesConflict{}, server.ErrDuplicateKey
		}
		return &api.PostRolesInternalServerError{}, server.ErrDatabase
	}
	return &api.Role{RoleID: roleID, RoleName: req.Value.RoleName}, nil
}

func (s *Service) PutRoleName(ctx context.Context, req api.OptRoleName, params api.PutRoleNameParams) (api.PutRoleNameRes, error) {
	if !req.Set {
		return &api.PutRoleNameNotAcceptable{}, server.ErrMissingBody
	}
	if err := s.conn.Queries.UpdateRoleName(ctx, dbm.UpdateRoleNameParams{
		RoleID: params.RoleID,
		RoleName: req.Value.RoleName,
	}); err != nil {
		var pgErr *pgconn.PgError
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return &api.PutRoleNameConflict{}, server.ErrDuplicateKey
		}
	}
	return &api.PutRoleNameOK{}, nil
}
