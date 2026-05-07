package service

import (
	"context"
	api "gamesync/internal/ogen"
	"gamesync/internal/query"
)

type Service struct {
	q *query.Query
}
func NewService(query *query.Query) *Service {
	return &Service{
		q: query,
	}
}

func (s *Service) GetHealth(ctx context.Context) error {
	return nil
}
func (s *Service) GetRoles(ctx context.Context) (api.GetRolesRes, error) {
	return nil, nil
}
func (s *Service) GetUserID(ctx context.Context, params api.GetUserIDParams) error {
	return nil
}
func (s *Service) GetUsers(ctx context.Context) (api.GetUsersRes, error) {
	return nil, nil
}
func (s *Service) PostRoles(ctx context.Context, req api.OptRole) (api.PostRolesRes, error) {
	return nil, nil
}
func (s *Service) PostUsers(ctx context.Context, req api.OptUserNew) (api.PostUsersRes, error) {
	return nil, nil
}

