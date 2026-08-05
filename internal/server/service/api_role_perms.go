package service

import (
	"context"
	"fmt"
	api "go.pabu.dev/gamesync/internal/ogen"
	"go.pabu.dev/gamesync/internal/server"
	"go.pabu.dev/gamesync/internal/server/dbm"
	"log/slog"
	"net/http"
	"slices"
)

func (s *Service) GetRolePerms(ctx context.Context, params api.GetRolePermsParams) (result api.PermNameArray, err error) {

	qtx, tx, err := s.db.BeginReadTX(ctx)
	if err != nil {
		return nil, server.NewInternalError(err, "failed starting tx")
	}
	defer func() {
		if recover() != nil || err != nil {
			if e := tx.Rollback(ctx); e != nil {
				slog.Error("failed rollback", "error", e)
			}
		}
	}()

	// we check first that the role exists before getting the perms it has.
	// if we do not do that, the next query will  return an empty slice,
	// even though the role does not exist.
	roleCount, err := qtx.GetRoleWithIdCount(ctx, params.RoleID)
	if err != nil {
		return nil, server.NewInternalError(err, "failed checking if role exists", "roleID", params.RoleID)
	}
	if roleCount < 1 {
		return nil, server.ErrRolePermsNotFound
	}

	permNames, err := qtx.ListRolePermNamesWithName(ctx, params.RoleID)
	if err != nil {
		return nil, server.NewInternalError(err, "failed listing perms for role", "roleID", params.RoleID)
	}
	result = slices.Grow(result, len(permNames))
	for _, n := range permNames {
		result = append(result, n)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, server.NewInternalError(err, "failed committing tx")
	}
	return result, nil
}

func (s *Service) PatchRolePerms(ctx context.Context, req *api.PermDiff, params api.PatchRolePermsParams) (result api.PatchRolePermsRes, err error) {

	permsAdd, err := s.db.ReadQuery().ListPermIDsWithNames(ctx, req.Add)
	if err != nil {
		return nil, server.NewInternalError(err, "failed listing 'add' perm names", "roleID", params.RoleID)
	}
	if len(permsAdd) != len(req.Add) {
		return nil, &api.GlobalErrorStatusCode{
			StatusCode: http.StatusUnprocessableEntity,
			Response: api.Error{
				Code:    http.StatusUnprocessableEntity,
				Message: fmt.Sprintf("There are %d invalid permission names in the add array", len(req.Add)-len(permsAdd)),
			},
		}
	}

	permsRemove, err := s.db.ReadQuery().ListPermIDsWithNames(ctx, req.Remove)
	if err != nil {
		return nil, server.NewInternalError(err, "failed listing 'delete' perm names", "roleID", params.RoleID)
	}
	if len(permsRemove) != len(req.Remove) {
		return nil, &api.GlobalErrorStatusCode{
			StatusCode: http.StatusUnprocessableEntity,
			Response: api.Error{
				Code:    http.StatusUnprocessableEntity,
				Message: fmt.Sprintf("There are %d invalid permission names in the delete array", len(req.Add)-len(permsAdd)),
			},
		}
	}

	// i decided to not include the 2 above queries as the permissions table does not really
	// update during runtime, so it is probably fine. i want to avoid querying the primary
	// database if not needed.
	qtx, tx, err := s.db.BeginTX(ctx)
	if err != nil {
		return nil, server.NewInternalError(err, "failed starting tx while patching role perms")
	}
	defer func() {
		if recover() != nil || err != nil {
			if e := tx.Rollback(ctx); e != nil {
				slog.Error("failed rollback", "error", e)
			}
		}
	}()

	if len(req.Add) > 0 {
		if err = qtx.InsertRolePermsNoConflict(ctx, dbm.InsertRolePermsNoConflictParams{
			RoleID:  params.RoleID,
			PermIds: permsAdd,
		}); err != nil {
			return nil, server.NewInternalError(err, "failed adding perms to role", "roleID", params.RoleID)
		}
	}

	if len(req.Remove) > 0 {
		if err = qtx.DeleteRolePermsWithId(ctx, dbm.DeleteRolePermsWithIdParams{
			RoleID:  params.RoleID,
			PermIds: permsRemove,
		}); err != nil {
			return nil, server.NewInternalError(err, "failed removing perms from role", "roleID", params.RoleID)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, server.NewInternalError(err, "failed committing perms for role")
	}

	return &api.PatchRolePermsOK{}, nil
}

func (s *Service) PutRolePerms(ctx context.Context, req api.PermNameArray, params api.PutRolePermsParams) (result api.PutRolePermsRes, err error) {
	qtx, tx, err := s.db.BeginTX(ctx)
	if err != nil {
		return nil, server.NewInternalError(err, "failed starting tx while setting role perms")
	}
	defer func() {
		if recover() != nil || err != nil {
			if e := tx.Rollback(ctx); e != nil {
				slog.Error("failed rollback", "error", e)
			}
		}
	}()

	perms, err := qtx.ListPermIDsWithNames(ctx, req)
	if err != nil {
		return nil, server.NewInternalError(err, "failed listing perm names while setting role perms", "roleID", params.RoleID)
	}
	if len(perms) != len(req) {
		return nil, &api.GlobalErrorStatusCode{
			StatusCode: http.StatusUnprocessableEntity,
			Response: api.Error{
				Code:    http.StatusUnprocessableEntity,
				Message: fmt.Sprintf("There are %d invalid permission names in the array", len(perms)-len(req)),
			},
		}
	}

	err = qtx.DeleteRolePermsWithRoleId(ctx, params.RoleID)
	if err != nil {
		return nil, server.NewInternalError(err, "failed deleting all perms for role", "roleID", params.RoleID)
	}

	if err = qtx.InsertRolePerms(ctx, dbm.InsertRolePermsParams{
		RoleID:  params.RoleID,
		PermIds: perms,
	}); err != nil {
		slog.Error("adding perms to role", "roleID", params.RoleID, "error", err)
		return nil, server.NewInternalError(err, "failed adding perms to role", "roleID", params.RoleID)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, server.NewInternalError(err, "failed committing tx while setting perms for role", "roleID", params.RoleID)
	}

	return &api.PutRolePermsOK{}, nil
}
