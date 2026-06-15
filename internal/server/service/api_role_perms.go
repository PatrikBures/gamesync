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
	// permsAdd, err := s.q.Permission.WithContext(ctx).Where(s.q.Permission.PermName.In(req.Value.Add...)).Find()
	// if err != nil {
	// 	return &api.PatchRolePermsInternalServerError{}, server.ErrDatabase
	// }
	// if len(permsAdd) != len(req.Value.Add) {
	// 	return &api.PatchRolePermsUnprocessableEntity{}, nil
	// }
	//
	// permsRemove, err := s.q.Permission.WithContext(ctx).Where(s.q.Permission.PermName.In(req.Value.Remove...)).Find()
	// if err != nil {
	// 	return &api.PatchRolePermsInternalServerError{}, server.ErrDatabase
	// }
	// if len(permsRemove) != len(req.Value.Remove) {
	// 	return &api.PatchRolePermsUnprocessableEntity{}, nil
	// }
	//
	//
	// tx := s.q.Begin()
	// defer func() {
	// 	if recover() != nil || err != nil {
	// 		if e := tx.Rollback(); e != nil {
	// 			slog.Error("failed rollback", "error", e)
	// 		}
	// 	}
	// }()
	//
	// if len(req.Value.Add) > 0 {
	// 	rolePerms := permToRolePermInsert(params.RoleID, permsAdd)
	// 	err = tx.RolePermission.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(rolePerms...)
	// 	if err != nil {
	// 		return &api.PatchRolePermsInternalServerError{}, server.ErrDatabase
	// 	}
	// }
	//
	// if len(req.Value.Remove) > 0 {
	// 	remove := make([]int32, 0, len(req.Value.Remove))
	// 	for _, pr := range permsRemove {
	// 		remove = append(remove, pr.PermID)
	// 	}
	// 	_, err = tx.RolePermission.WithContext(ctx).
	// 		Where(
	// 			s.q.RolePermission.RoleID.Eq(params.RoleID),
	// 			s.q.RolePermission.PermID.In(remove...),
	// 		).Delete()
	// 	if err != nil {
	// 		return &api.PatchRolePermsInternalServerError{}, server.ErrDatabase
	// 	}
	// }
	//
	// if err := tx.Commit(); err != nil {
	// 	return &api.PatchRolePermsInternalServerError{}, server.ErrDatabase
	// }

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
