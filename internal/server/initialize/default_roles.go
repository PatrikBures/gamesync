package initServer

import (
	"context"
	"errors"
	"gamesync/internal/model"
	"gamesync/internal/query"
	"gamesync/internal/server/permissions"
	"log/slog"
	"slices"

	"gorm.io/gen"
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

func CreateDefaultRolePerms(q *query.Query) error {
	ctx := context.Background()
	rolePerms := []*model.RolePermission{}
	
	rolePerms = append(rolePerms, permsToRolePermModel(1, permissions.AllPerms)...)
	rolePerms = append(rolePerms, permsToRolePermModel(50, permissions.Perms{
		permissions.PermRolesGet,
		permissions.PermUserGetOwn,
		permissions.PermUserNameUpdateOwn,
	})...)


	var rolesRemovedCount gen.ResultInfo
	err := q.Transaction(func(tx *query.Query) error {
		var err error
		rolesRemovedCount, err = tx.RolePermission.WithContext(ctx).Where(q.RolePermission.RoleID.Lt(100)).Delete()
		if err != nil {
			return err
		}

		if err := tx.RolePermission.WithContext(ctx).Create(rolePerms...); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil
	}
	slog.Info("removed role perms under 100", "count", rolesRemovedCount.RowsAffected)
	slog.Info("created perms for roles under 100", "count", len(rolePerms))

	return nil
}

func permsToRolePermModel(roleID int32, perms permissions.Perms) []*model.RolePermission {
	m := make([]*model.RolePermission, 0, len(perms))
	for _, p := range perms {
		rp := &model.RolePermission{
			RoleID: roleID,
			PermID: int32(p),
		}
		m = append(m, rp)
	}

	return m
}
