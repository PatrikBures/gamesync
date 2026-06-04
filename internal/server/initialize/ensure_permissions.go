package initServer

import (
	"context"
	"gamesync/internal/dbx"
	"gamesync/internal/server/dbm"
	"gamesync/internal/server/permissions"
	"log/slog"
	"slices"
)



func EnsurePermissions(conn dbx.DBconn) (err error) {
	ctx := context.Background()

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

	currentPerms, err := qtx.ListPermissions(ctx)
	if err != nil {
		return err
	}
	currentPermIds := make([]int32, 0, len(currentPerms))
	for _, c := range currentPerms {
		currentPermIds = append(currentPermIds, c.PermID)
	}

	expectedPermsInt32 := make([]int32, 0, len(permissions.AllPerms))

	loop: for _, e := range permissions.AllPerms {
		i := int32(e)
		expectedPermsInt32 = append(expectedPermsInt32, i)
		perm := dbm.Permission{PermID: i, PermName: e.String()}

		for _, cp := range currentPerms {
			if cp.PermID == perm.PermID {
				if cp.PermName != perm.PermName {
					err = qtx.UpdatePermissionName(ctx, dbm.UpdatePermissionNameParams(perm))
					if err != nil {
						return
					}
					slog.Info("updated perm name", "permission", perm, "old_permission", cp)
				}
				slog.Info("already exists", "permission", perm)
				continue loop
			}
		}

		err = qtx.InsertPermission(ctx, dbm.InsertPermissionParams(perm))
		if err != nil {
			return
		}
		slog.Info("created permission", "permission", perm)
	}

	for _, c := range currentPermIds {
		if slices.Contains(expectedPermsInt32, c) {
			continue
		}
		var perm dbm.Permission
		perm, err = qtx.DeletePermission(ctx, c)
		if err != nil {
			return
		}
		slog.Info("deleted permission", "permission", perm)
	}
	return tx.Commit(ctx)
}
