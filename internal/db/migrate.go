package db

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed migrations/sqlite/*sql
var embedSQLiteMigrations embed.FS

func Migrate() error {
	db, err := sql.Open("sqlite", "gamesync.db")
	if err != nil {
		return fmt.Errorf("opening db: %v", err)
	}

	goose.SetBaseFS(embedSQLiteMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		return fmt.Errorf("setting goose dialect: %v", err)
	}
	if err := goose.Up(db, "migrations/sqlite"); err != nil {
		return fmt.Errorf("migrating: %v", err)
	}
	return nil
}
