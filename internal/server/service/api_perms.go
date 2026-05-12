package service

import (
	"context"
	api "gamesync/internal/ogen"
)


func (s *Service) GetPerms(ctx context.Context) (api.GetPermsRes, error) {
	perms, err := s.q.Permission.WithContext(ctx).Find()
	if err != nil {
		return &api.GetPermsInternalServerError{}, ErrDatabase
	}
	permsReturn := make(api.GetPermsOKApplicationJSON, 0, len(perms))
	for _, perm := range perms {
		permsReturn = append(permsReturn, api.PermWithName{
			PermID: api.Perm(perm.PermID),
			PermName: perm.PermName,
		})
	}
	return &permsReturn, nil
}
