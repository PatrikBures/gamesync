package initServer

import (
	"context"
	"go.pabu.dev/gamesync/internal/dbx"
	"go.pabu.dev/gamesync/internal/server/dbm"
	"go.pabu.dev/gamesync/internal/server/permissions"
	"log/slog"
	"slices"
)

func CreateDefaultRoles(db *dbx.DB, skipRoles []string) error {
	ctx := context.Background()
	roles := []dbm.InsertRoleWithIdParams{
		{RoleID: 1, RoleName: "admin"},
		{RoleID: 50, RoleName: "standard"},
		{RoleID: 60, RoleName: "maintainer"},
		{RoleID: 99, RoleName: "none"},
	}
	for _, role := range roles {
		existingRoleCount, err := db.ReadQuery().GetRoleWithIdCount(ctx, role.RoleID)
		if err != nil {
			return err
		}

		if slices.Contains(skipRoles, role.RoleName) {
			if existingRoleCount > 0 {
				err := db.WriteQuery().DeleteRole(ctx, role.RoleID)
				if err != nil {
					return err
				}
				slog.Info("deleted role", "role", role)
				continue
			}
			slog.Info("skipped creating role", "role", role)
			continue
		}

		if existingRoleCount > 0 {
			slog.Info("role already exists", "roleID", role.RoleID, "roleName", role.RoleName)
			continue
		}

		if err := db.WriteQuery().InsertRoleWithId(ctx, role); err != nil {
			return err
		}

		slog.Info("created role", "role", role)
	}
	return nil
}

func CreateDefaultRolePerms(db *dbx.DB) (err error) {
	ctx := context.Background()
	rolePerms := []dbm.InsertRolePermsCopyParams{}

	rolePerms = append(rolePerms, permsToRolePermModel(1, permissions.AllPerms)...)
	rolePerms = append(rolePerms, permsToRolePermModel(50, permissions.Perms{
		permissions.PermRolesGet,
		permissions.PermUserGetOwn,
		permissions.PermUserNameUpdateOwn,
		permissions.PermSync,
	})...)

	qtx, tx, err := db.BeginTX(ctx)
	if err != nil {
		return
	}

	defer func() {
		if recover() != nil || err != nil {
			if e := tx.Rollback(ctx); e != nil {
				slog.Error("failed rollback", "error", e)
			}
		}
	}()

	rolesDeletedCount, err := qtx.DeleteRolePermsLT(ctx, 100)
	if err != nil {
		return err
	}

	permsInsertedCount, err := qtx.InsertRolePermsCopy(ctx, rolePerms)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	slog.Info("removed role perms under 100", "count", rolesDeletedCount)
	slog.Info("created perms for roles under 100", "total_perms_inserted", permsInsertedCount)

	return nil
}

func permsToRolePermModel(roleID int32, perms permissions.Perms) []dbm.InsertRolePermsCopyParams {
	m := make([]dbm.InsertRolePermsCopyParams, 0, len(perms))
	for _, p := range perms {
		rp := dbm.InsertRolePermsCopyParams{
			RoleID: roleID,
			PermID: int32(p),
		}
		m = append(m, rp)
	}

	return m
}
