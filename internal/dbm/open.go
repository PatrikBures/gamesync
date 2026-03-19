package dbm

import (
	"database/sql"
	"fmt"
	"gamesync/internal/vars"
	_ "modernc.org/sqlite"
)

func OpenSQLite() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "file:"+vars.RemoteSQLiteDb+"?_busy_timeout=1000")
	if err != nil {
		return nil, fmt.Errorf("opening SQLite db: %w", err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("enabling foreign keys: %w", err)
	}
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("enabling WAL: %w", err)
	}
	return db, nil
}

func CloseDB(db *sql.DB, err *error) {
	if cerr := db.Close(); err != nil {
		if *err != nil{
			*err = fmt.Errorf("%v; additionally, %w", *err, cerr)
		} else {
			*err = cerr
		}
	}
}
