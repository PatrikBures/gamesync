package roles

import (
	"context"
	"fmt"
	"gamesync/internal/model"
	"gamesync/internal/query"
	"log/slog"
)


func CreateDefaultRoles(q *query.Query) error {
	roles := []*model.Role{
		{ RoleID:    1, RoleName: "admin" },
		{ RoleID:  50, RoleName: "standard" },
		{ RoleID:  99, RoleName: "none" },
	}
	if err := q.Role.WithContext(context.Background()).Create(roles...); err != nil {
		return fmt.Errorf("creating %d default roles: %w", len(roles), err)	
	}
	for _, role := range roles {
		slog.Info("created", "role", *role)
	}
	return nil
}
