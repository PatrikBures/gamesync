package service

import (
	"context"
	"gamesync/internal/model"
	api "gamesync/internal/ogen"

	"gorm.io/gorm/clause"
)


func (s *Service) GetRolePerms(ctx context.Context, params api.GetRolePermsParams) (api.GetRolePermsRes, error) {
	return nil, nil
}

func (s *Service) PatchRolePerms(ctx context.Context, req api.OptPermDiff, params api.PatchRolePermsParams) (result api.PatchRolePermsRes, err error) {
	if ! req.Set {
		return &api.PatchRolePermsNotAcceptable{}, ErrMissingBody
	}
	permsAdd, err := s.q.Permission.WithContext(ctx).Where(s.q.Permission.PermName.In(req.Value.Add...)).Find()
	if err != nil {
		return &api.PatchRolePermsInternalServerError{}, ErrDatabase
	}
	if len(permsAdd) != len(req.Value.Add) {
		return &api.PatchRolePermsUnprocessableEntity{}, ErrPermNotFound
	}

	permsRemove, err := s.q.Permission.WithContext(ctx).Where(s.q.Permission.PermName.In(req.Value.Remove...)).Find()
	if err != nil {
		return &api.PatchRolePermsInternalServerError{}, ErrDatabase
	}
	if len(permsRemove) != len(req.Value.Remove) {
		return &api.PatchRolePermsUnprocessableEntity{}, ErrPermNotFound
	}


	tx := s.q.Begin()
	   defer func() {
	       if recover() != nil || err != nil {
	           _ = tx.Rollback()
	       }
	   }()

	if len(req.Value.Add) > 0 {
		rolePerms := permToRolePerm(permsAdd, params.RoleID)
		err = tx.RolePermission.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(rolePerms...)
		if err != nil {
			return &api.PatchRolePermsInternalServerError{}, ErrDatabase
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
			return &api.PatchRolePermsInternalServerError{}, ErrDatabase
		}
	}

	if err := tx.Commit(); err != nil {
		return &api.PatchRolePermsInternalServerError{}, ErrDatabase
	}

	return &api.PatchRolePermsCreated{}, nil
}

func (s *Service) PutRolePerms(ctx context.Context, req api.PermNameArray, params api.PutRolePermsParams) (api.PutRolePermsRes, error) {
	return nil, nil
}

func permToRolePerm(perms []*model.Permission, roleID int32) []*model.RolePermission {
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
