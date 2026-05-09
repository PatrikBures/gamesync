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
		{ RoleID: 60, RoleName: "maintainer" },
		{ RoleID: 99, RoleName: "none" },
	}
	for _, role := range roles {
		existingRoleCount, err := q.Role.WithContext(ctx).Where(q.Role.RoleID.Eq(role.RoleID)).Count()
		if err != nil {
			return err
		}

		if slices.Contains(skipRoles, role.RoleName) {
			if existingRoleCount > 0 {
				info, err := q.Role.WithContext(ctx).Delete(role)
				if err != nil {
					if errors.Is(err, gorm.ErrForeignKeyViolated) {
						slog.Error("could not delete role. there are probably users with that role preventing deletion.")
						return err
					}

					return err
				}
				slog.Info("deleted", "role", *role, "info", info)
				continue
			}
			slog.Info("skipped creating", "role", *role)
			continue
		}

		if existingRoleCount > 0 {
			slog.Info("already exists", "role", *role)
			continue
		}

		if err := q.Role.WithContext(ctx).Create(role); err != nil {
			return err
		}

		slog.Info("created", "role", *role)
	}
	return nil
}
