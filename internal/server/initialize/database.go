package initServer

import (
	"gamesync/internal/dbx"
	"gamesync/internal/query"
	"log/slog"
)

func InitDatabase(dbType string, dbUrl string, disabledRoles []string) (*query.Query, error) {
	db, err := dbx.ConnectDb(dbType, dbUrl)
	if err != nil {
		return nil, err
	}
	if dbType == "sqlite" {
		db.Exec("PRAGMA foreign_keys = ON")
	}
	q := query.Use(db)

	if err := EnsurePermissions(q); err != nil {
		return nil, err
	}
	if err := CreateDefaultRoles(q, disabledRoles); err != nil {
		return nil, err
	}
	if err := CreateDefaultRolePerms(q); err != nil {
		return nil, err
	}
	if token, err := CreateAdmin(q); err == nil {
		if token == "" {
			slog.Info("admin already exists")
		} else {
			slog.Info("Created admin, make sure to update the token", "token", token)
		}
	}
	return q, nil
}
