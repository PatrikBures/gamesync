package service

import (
	"context"
	"gamesync/internal/model"
	api "gamesync/internal/ogen"
	"gamesync/internal/query"
	"gamesync/internal/server"
	"log/slog"

	"gorm.io/gorm/clause"
)

type GetRolePermsResult struct {
	PermName string
}

func (s *Service) GetRolePerms(ctx context.Context, params api.GetRolePermsParams) (api.GetRolePermsRes, error) {
	roleCount, err := s.q.Role.WithContext(ctx).
		Where(s.q.Role.RoleID.Eq(params.RoleID)).Count()
	if err != nil {
		return &api.GetRolePermsInternalServerError{}, server.ErrDatabase
	}
	if roleCount < 1 {
		return &api.GetRolePermsNotFound{}, server.ErrNotFound
	}

	var result []GetRolePermsResult
	err = s.q.RolePermission.WithContext(ctx).
		Select(s.q.Permission.PermName).
		Join(s.q.Permission, s.q.RolePermission.PermID.EqCol(s.q.Permission.PermID)).
		Where(s.q.RolePermission.RoleID.Eq(params.RoleID)).
		Scan(&result)
	if err != nil {
		return &api.GetRolePermsInternalServerError{}, server.ErrDatabase
	}
	permNames := make(api.PermNameArray, 0, len(result))
	for _, r := range result {
		permNames = append(permNames, r.PermName)
	}
	return &permNames, nil
}

func (s *Service) PatchRolePerms(ctx context.Context, req api.OptPermDiff, params api.PatchRolePermsParams) (result api.PatchRolePermsRes, err error) {
	if ! req.Set {
		return &api.PatchRolePermsNotAcceptable{}, server.ErrMissingBody
	}
	permsAdd, err := s.q.Permission.WithContext(ctx).Where(s.q.Permission.PermName.In(req.Value.Add...)).Find()
	if err != nil {
		return &api.PatchRolePermsInternalServerError{}, server.ErrDatabase
	}
	if len(permsAdd) != len(req.Value.Add) {
		return &api.PatchRolePermsUnprocessableEntity{}, server.ErrPermNotFound
	}

	permsRemove, err := s.q.Permission.WithContext(ctx).Where(s.q.Permission.PermName.In(req.Value.Remove...)).Find()
	if err != nil {
		return &api.PatchRolePermsInternalServerError{}, server.ErrDatabase
	}
	if len(permsRemove) != len(req.Value.Remove) {
		return &api.PatchRolePermsUnprocessableEntity{}, server.ErrPermNotFound
	}


	tx := s.q.Begin()
	defer func() {
		if recover() != nil || err != nil {
			if e := tx.Rollback(); e != nil {
				slog.Error("failed rollback", "error", e)
			}
		}
	}()

	if len(req.Value.Add) > 0 {
		rolePerms := permToRolePerm(params.RoleID, permsAdd)
		err = tx.RolePermission.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(rolePerms...)
		if err != nil {
			return &api.PatchRolePermsInternalServerError{}, server.ErrDatabase
		}
	}

	if len(req.Value.Remove) > 0 {
		remove := make([]int32, 0, len(req.Value.Remove))
		for _, pr := range permsRemove {
			remove = append(remove, pr.PermID)
		}
		_, err = tx.RolePermission.WithContext(ctx).
			Where(
				s.q.RolePermission.RoleID.Eq(params.RoleID),
				s.q.RolePermission.PermID.In(remove...),
			).Delete()
		if err != nil {
			return &api.PatchRolePermsInternalServerError{}, server.ErrDatabase
		}
	}

	if err := tx.Commit(); err != nil {
		return &api.PatchRolePermsInternalServerError{}, server.ErrDatabase
	}

	return &api.PatchRolePermsOK{}, nil
}

func (s *Service) PutRolePerms(ctx context.Context, req api.PermNameArray, params api.PutRolePermsParams) (api.PutRolePermsRes, error) {
	perms, err := s.q.Permission.WithContext(ctx).Where(s.q.Permission.PermName.In(req...)).Find()
	if err != nil {
		return &api.PutRolePermsInternalServerError{}, server.ErrDatabase
	}
	if len(perms) != len(req) {
		return &api.PutRolePermsUnprocessableEntity{}, server.ErrPermNotFound
	}

	err = s.q.Transaction(func(tx *query.Query) error {
		if _, err := tx.RolePermission.WithContext(ctx).Where(tx.RolePermission.RoleID.Eq(params.RoleID)).Delete(); err != nil {
			return err
		}
		permIDs := permToRolePerm(params.RoleID, perms)
		if err := tx.RolePermission.WithContext(ctx).Create(permIDs...); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return &api.PutRolePermsInternalServerError{}, server.ErrDatabase
	}

	return nil, nil
}

func permToRolePerm(roleID int32, perms []*model.Permission) []*model.RolePermission {
	rolePerms := make([]*model.RolePermission, 0, len(perms))
	for _, p := range perms {
		rolePerm := model.RolePermission{
			PermID: p.PermID, 
			RoleID: roleID,
		}
		rolePerms = append(rolePerms, &rolePerm)
	}
	return rolePerms
}
