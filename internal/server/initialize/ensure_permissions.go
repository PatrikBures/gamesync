package initServer

import (
	"context"
	"gamesync/internal/model"
	"gamesync/internal/query"
	"gamesync/internal/server/permissions"
	"log/slog"
	"slices"
)



func EnsurePermissions(q *query.Query) (err error) {
	ctx := context.Background()
	tx := q.Begin()
	defer func() {
		if recover() != nil || err != nil {
			_ = tx.Rollback()
		}
	}()

	currentPerms, err := q.Permission.WithContext(ctx).Find()
	if err != nil {
		return err
	}
	currentPermIds := make([]int32, len(currentPerms))
	for _, c := range currentPerms {
		currentPermIds = append(currentPermIds, c.PermID)
	}

	expectedPermsInt32 := make([]int32, len(permissions.AllPerms))

	for _, e := range permissions.AllPerms {
		i := int32(e)
		expectedPermsInt32 = append(expectedPermsInt32, i)
		perm := &model.Permission{PermID: i, PermName: e.String()}
		if slices.Contains(currentPermIds, i) {
			slog.Info("already exists", "permission", *perm)
			continue
		}
		tx.Permission.WithContext(ctx).Create(perm)
		slog.Info("created", "permission", *perm)
	}

	for _, c := range currentPermIds {
		if slices.Contains(expectedPermsInt32, c) {
			continue
		}
		perm := &model.Permission{PermID: c}
		tx.Permission.WithContext(ctx).Delete(perm)
		slog.Info("deleted", "permission", *perm)
	}
	return tx.Commit()
}
