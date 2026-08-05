package service

import (
	"context"
	api "go.pabu.dev/gamesync/internal/ogen"
	"go.pabu.dev/gamesync/internal/server"
	"go.pabu.dev/gamesync/internal/server/dbm"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Service) GetRoles(ctx context.Context) ([]api.Role, error) {
	roles, err := s.db.ReadQuery().ListRoles(ctx)
	if err != nil {
		return nil, server.NewInternalError(err, "failed listing roles")
	}
	rolesReturn := make([]api.Role, 0, len(roles))
	for _, role := range roles {
		rolesReturn = append(rolesReturn, api.Role{
			RoleID:   role.RoleID,
			RoleName: role.RoleName,
		})
	}
	return rolesReturn, nil
}

func (s *Service) PostRoles(ctx context.Context, req *api.RoleName) (*api.Role, error) {
	roleID, err := s.db.WriteQuery().InsertRole(ctx, req.RoleName)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return nil, server.ErrRoleNameConflict
			}
		}
		return nil, server.NewInternalError(err, "failed inserting role", "roleName", req.RoleName)
	}
	return &api.Role{RoleID: roleID, RoleName: req.RoleName}, nil
}

func (s *Service) PutRoleName(ctx context.Context, req *api.RoleName, params api.PutRoleNameParams) (error) {
	if err := s.db.WriteQuery().UpdateRoleName(ctx, dbm.UpdateRoleNameParams{
		RoleID:   params.RoleID,
		RoleName: req.RoleName,
	}); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return server.ErrRoleNameConflict
			case pgerrcode.ForeignKeyViolation:
				return server.ErrRoleNotFound
			}
		}
		return server.NewInternalError(err, "failed updating role name", "roleID", params.RoleID, "newRoleName", req.RoleName)
	}
	return nil
}
