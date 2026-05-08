package roles

import (
	"context"
	"errors"
	"gamesync/internal/model"
	"gamesync/internal/query"
	"log/slog"
	"slices"

	"gorm.io/gorm"
)


func CreateDefaultRoles(q *query.Query, skipRoles []string) error {
	ctx := context.Background()
	roles := []*model.Role{
		{ RoleID:  1, RoleName: "admin" },
		{ RoleID: 50, RoleName: "standard" },
		{ RoleID: 99, RoleName: "none" },
	}
	for _, role := range roles {
		if slices.Contains(skipRoles, role.RoleName) {
			slog.Info("skipped creating", "role", *role)
			continue
		}

		if err := q.Role.WithContext(ctx).Create(role); err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				slog.Info("already exists", "role", *role)
				continue
			}
			return err
		}

		slog.Info("created", "role", *role)
	}
	return nil
}
