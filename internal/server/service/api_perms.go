package service

import (
	"context"
	api "go.pabu.dev/gamesync/internal/ogen"
	"go.pabu.dev/gamesync/internal/server"
)

func (s *Service) GetPerms(ctx context.Context) ([]api.PermWithName, error) {
	perms, err := s.db.ReadQuery().ListPermissions(ctx)
	if err != nil {
		return nil, server.NewInternalError(err, "failed listing permissions", "error", err)
	}
	permsReturn := make([]api.PermWithName, 0, len(perms))
	for _, perm := range perms {
		permsReturn = append(permsReturn, api.PermWithName{
			PermID:   api.Perm(perm.PermID),
			PermName: perm.PermName,
		})
	}
	return permsReturn, nil
}
