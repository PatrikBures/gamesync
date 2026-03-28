package dbm

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/sqlite/*sql
var embedSQLiteMigrations embed.FS

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(embedSQLiteMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		return fmt.Errorf("setting goose dialect: %v", err)
	}
	if err := goose.Up(db, "migrations/sqlite"); err != nil {
		return err
	}
	return nil
}
