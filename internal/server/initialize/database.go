package initServer

import (
	"fmt"
	"gamesync/internal/dbx"
	"log/slog"
)

func InitDatabase(dbUrl string, disabledRoles []string) (conn dbx.DBconn, err error) {
	conn, err = dbx.ConnectDb(dbUrl)
	if err != nil {
		err = fmt.Errorf("connecting db: %w", err)
		return
	}

	if err = EnsurePermissions(conn); err != nil {
		return
	}
	if err = CreateDefaultRoles(conn, disabledRoles); err != nil {
		return
	}
	if err = CreateDefaultRolePerms(conn); err != nil {
		return
	}

	token, err := CreateAdmin(conn)
	if err != nil {
		err = fmt.Errorf("creating admin: %w", err)
		return
	}

	if token == "" {
		slog.Info("admin already exists")
	} else {
		slog.Info("Created admin, make sure to update the token", "token", token)
	}

	return conn, nil
}
