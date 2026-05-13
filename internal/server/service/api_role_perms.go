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

func (s *Service) PatchRolePerms(ctx context.Context, req api.OptPermDiff, params api.PatchRolePermsParams) (api.PatchRolePermsRes, error) {
	add := make([]*model.RolePermission, 0, len(req.Value.Add))
	for _, a := range req.Value.Add {
		r := &model.RolePermission{
			RoleID: params.RoleID,
			PermID: int32(a),
		}
		add = append(add, r)
	}
	if err := s.q.RolePermission.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(add...); err != nil {
		return &api.PatchRolePermsInternalServerError{}, ErrDatabase
	}
	return &api.PatchRolePermsCreated{}, nil
}

func (s *Service) PutRolePerms(ctx context.Context, req api.PermArray, params api.PutRolePermsParams) (api.PutRolePermsRes, error) {
	return nil, nil
}
