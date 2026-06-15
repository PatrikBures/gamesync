package service

import (
	"context"
	api "gamesync/internal/ogen"
	"gamesync/internal/server"
	"gamesync/internal/server/dbm"
	"log/slog"
)

type GetRolePermsResult struct {
	PermName string
}

func (s *Service) GetRolePerms(ctx context.Context, params api.GetRolePermsParams) (result api.GetRolePermsRes, err error) {

	qtx, tx, err := s.db.BeginReadTX(ctx)
	if err != nil {
		return
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
		return &api.GetRolePermsInternalServerError{}, server.ErrDatabase
	}
	if roleCount < 1 {
		return &api.GetRolePermsNotFound{}, nil
	}

	permNames, err := qtx.ListRolePermNamesWithName(ctx, params.RoleID)
	if err != nil {
		return &api.GetRolePermsInternalServerError{}, server.ErrDatabase
	}
	permNamesReturn := make(api.PermNameArray, 0, len(permNames))
	for _, n := range permNames {
		permNamesReturn = append(permNamesReturn, n)
	}
	if err = tx.Commit(ctx); err != nil {
		return &api.GetRolePermsInternalServerError{}, server.ErrDatabase
	}
	return &permNamesReturn, nil
}

func (s *Service) PatchRolePerms(ctx context.Context, req api.OptPermDiff, params api.PatchRolePermsParams) (result api.PatchRolePermsRes, err error) {
	if ! req.Set {
		return &api.PatchRolePermsNotAcceptable{}, nil
	}

	permsAdd, err := s.db.ReadQuery().ListPermIDsWithNames(ctx, req.Value.Add)
	if err != nil {
		slog.Error("listing add values", "error", err)
		return &api.PatchRolePermsInternalServerError{}, server.ErrDatabase
	}
	if len(permsAdd) != len(req.Value.Add) {
		slog.Warn("add perm(s) not found", "expected", len(req.Value.Add), "received", len(permsAdd), "received", permsAdd)
		return &api.PatchRolePermsUnprocessableEntity{}, nil
	}

	permsRemove, err := s.db.ReadQuery().ListPermIDsWithNames(ctx, req.Value.Remove)
	if err != nil {
		slog.Error("listing remove values", "error", err)
		return &api.PatchRolePermsInternalServerError{}, server.ErrDatabase
	}
	if len(permsRemove) != len(req.Value.Remove) {
		slog.Warn("remove perm(s) not found", "expectedQty", len(req.Value.Remove), "receivedQty", len(permsRemove), "received", permsRemove)
		return &api.PatchRolePermsUnprocessableEntity{}, nil
	}


	// i decided to not include the 2 above queries as the permissions table does not really
	// update during runtime, so it is probably fine. i want to avoid querying the primary
	// database if not needed.
	qtx, tx, err := s.db.BeginTX(ctx)
	if err != nil {
		slog.Error("starting tx", "error", err)
		return &api.PatchRolePermsInternalServerError{}, server.ErrDatabase
	}
	defer func() {
		if recover() != nil || err != nil {
			if e := tx.Rollback(ctx); e != nil {
				slog.Error("failed rollback", "error", e)
			}
		}
	}()

	if len(req.Value.Add) > 0 {
		if err = qtx.InsertRolePermsNoConflict(ctx, dbm.InsertRolePermsNoConflictParams{
			RoleID: params.RoleID,
			PermIds: permsAdd,
		}); err != nil {
			slog.Error("adding role perms", "error", err)
			return &api.PatchRolePermsInternalServerError{}, server.ErrDatabase
		}
	}

	if len(req.Value.Remove) > 0 {
		if err = qtx.DeleteRolePermsWithId(ctx, dbm.DeleteRolePermsWithIdParams{
			RoleID: params.RoleID,
			PermIds: permsRemove,
		}); err != nil {
			slog.Error("removing role perms", "error", err)
			return &api.PatchRolePermsInternalServerError{}, server.ErrDatabase
		}
	}

	if err := tx.Commit(ctx); err != nil {
			slog.Error("committing role perms", "error", err)
		return &api.PatchRolePermsInternalServerError{}, server.ErrDatabase
	}

	return &api.PatchRolePermsOK{}, nil
}

func (s *Service) PutRolePerms(ctx context.Context, req api.PermNameArray, params api.PutRolePermsParams) (result api.PutRolePermsRes, err error) {
	qtx, tx, err := s.db.BeginTX(ctx)
	if err != nil {
		slog.Error("starting tx", "error", err)
		return &api.PutRolePermsInternalServerError{}, server.ErrDatabase
	}
	defer func() {
		if recover() != nil || err != nil {
			if e := tx.Rollback(ctx); e != nil {
				slog.Error("failed rollback", "error", e)
			}
		}
	}()

	perms, err := qtx.ListPermsWithNames(ctx, req)
	if err != nil {
		return &api.PutRolePermsInternalServerError{}, server.ErrDatabase
	}
	if len(perms) != len(req) {
		return &api.PutRolePermsUnprocessableEntity{}, nil
	}

	err = qtx.DeleteRolePermsWithRoleId(ctx, params.RoleID)
	if err != nil {
		slog.Error("deleting all perms from role", "roleID", params.RoleID, "error", err)
		return &api.PutRolePermsInternalServerError{}, server.ErrDatabase
	}

	permIDs := permToRolePermInsert(params.RoleID, perms)
	_, err = qtx.InsertRolePerms(ctx, permIDs)
	if err != nil {
		slog.Error("adding perms to role", "roleID", params.RoleID, "error", err)
		return &api.PutRolePermsInternalServerError{}, server.ErrDatabase
	}

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("commiting transaction", "error", err)
		return &api.PutRolePermsInternalServerError{}, server.ErrDatabase
	}

	return &api.PutRolePermsOK{}, nil
}

// converts permissions slice to a slice which can be inserted to database
//
// every object in returned slice will have roleID
func permToRolePermInsert(roleID int32, perms []dbm.Permission) []dbm.InsertRolePermsParams {
	rolePerms := make([]dbm.InsertRolePermsParams, 0, len(perms))
	for _, p := range perms {
		rolePerm := dbm.InsertRolePermsParams{
			PermID: p.PermID, 
			RoleID: roleID,
		}
		rolePerms = append(rolePerms, rolePerm)
	}
	return rolePerms
}

func permToSlice(perms []dbm.Permission) []int32 {
	permsInt32 := make([]int32, 0, len(perms))
	for _, p := range perms {
		permsInt32 = append(permsInt32, p.PermID)
	}
	return permsInt32
}
