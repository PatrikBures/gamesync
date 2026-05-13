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
	tx := s.q.Begin()
    defer func() {
        if recover() != nil || err != nil {
            _ = tx.Rollback()
        }
    }()

	if len(req.Value.Add) > 0 {
		add := apiRoleToDbRole(req.Value.Add, params.RoleID)
		if err = tx.RolePermission.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(add...); err != nil {
			return &api.PatchRolePermsInternalServerError{}, ErrDatabase
		}
	}

	if len(req.Value.Remove) > 0 {
		remove := make([]int32, 0, len(req.Value.Remove))
		for _, r := range req.Value.Remove {
			remove = append(remove, int32(r))
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

func (s *Service) PutRolePerms(ctx context.Context, req api.PermArray, params api.PutRolePermsParams) (api.PutRolePermsRes, error) {
	return nil, nil
}

func apiRoleToDbRole(perms api.PermArray, roleID int32) []*model.RolePermission {
	newPerms := make([]*model.RolePermission, 0, len(perms))
	for _, a := range perms {
		r := &model.RolePermission{
			RoleID: roleID,
			PermID: int32(a),
		}
		newPerms = append(newPerms, r)
	}
	return newPerms
}
