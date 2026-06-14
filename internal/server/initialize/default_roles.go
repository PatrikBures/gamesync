package initServer

import (
	"context"
	"gamesync/internal/dbx"
	"gamesync/internal/server/dbm"
	"gamesync/internal/server/permissions"
	"log/slog"
	"slices"
)


func CreateDefaultRoles(conn dbx.DBconn, skipRoles []string) error {
	ctx := context.Background()
	roles := []dbm.InsertRoleWithIdParams{
		{ RoleID:  1, RoleName: "admin" },
		{ RoleID: 50, RoleName: "standard" },
		{ RoleID: 60, RoleName: "maintainer" },
		{ RoleID: 99, RoleName: "none" },
	}
	for _, role := range roles {
		existingRoleCount, err := conn.Queries.GetRoleWithIdCount(ctx, role.RoleID)
		if err != nil {
			return err
		}

		if slices.Contains(skipRoles, role.RoleName) {
			if existingRoleCount > 0 {
				err := conn.Queries.DeleteRole(ctx, role.RoleID)
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
			slog.Info("role already exists", "role", role)
			continue
		}

		if err := conn.Queries.InsertRoleWithId(ctx, role); err != nil {
			return err
		}

		slog.Info("created role", "role", role)
	}
	return nil
}

func CreateDefaultRolePerms(conn dbx.DBconn) (err error) {
	ctx := context.Background()
	rolePerms := []dbm.InsertRolePermsParams{}

	
	rolePerms = append(rolePerms, permsToRolePermModel(1, permissions.AllPerms)...)
	rolePerms = append(rolePerms, permsToRolePermModel(50, permissions.Perms{
		permissions.PermRolesGet,
		permissions.PermUserGetOwn,
		permissions.PermUserNameUpdateOwn,
	})...)

	tx, err := conn.Pool.Begin(ctx)
	if err != nil {
		return
	}
	qtx := conn.Queries.WithTx(tx)

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

	permsInsertedCount, err := qtx.InsertRolePerms(ctx, rolePerms)
	if err != nil {
		return err
	}
	
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	slog.Info("removed role perms under 100", "count", rolesDeletedCount)
	slog.Info("created perms for roles under 100", "roles", len(rolePerms), "total_perms_inserted", permsInsertedCount)

	return nil
}

func permsToRolePermModel(roleID int32, perms permissions.Perms) []dbm.InsertRolePermsParams {
	m := make([]dbm.InsertRolePermsParams, 0, len(perms))
	for _, p := range perms {
		rp := dbm.InsertRolePermsParams{
			RoleID: roleID,
			PermID: int32(p),
		}
		m = append(m, rp)
	}

	return m
}
