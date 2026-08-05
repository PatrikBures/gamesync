package initServer

import (
	"context"
	"fmt"
	"go.pabu.dev/gamesync/internal/dbx"
	"log/slog"
)

func InitDatabase(ctx context.Context, dbWriteUrl string, dbReadUrls []string, disabledRoles []string) (db *dbx.DB, err error) {
	db, err = dbx.NewDB(ctx, dbWriteUrl, dbReadUrls)
	if err != nil {
		err = fmt.Errorf("connecting db: %w", err)
		return
	}

	if err = EnsurePermissions(db); err != nil {
		return
	}
	if err = CreateDefaultRoles(db, disabledRoles); err != nil {
		return
	}
	if err = CreateDefaultRolePerms(db); err != nil {
		return
	}

	token, err := CreateAdmin(db)
	if err != nil {
		err = fmt.Errorf("creating admin: %w", err)
		return
	}

	if token == "" {
		slog.Info("admin already exists")
	} else {
		slog.Info("Created admin, make sure to update the token", "token", token)
	}

	return db, nil
}
