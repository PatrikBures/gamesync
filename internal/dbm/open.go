package dbm

import (
	"database/sql"
	"fmt"
	"gamesync/internal/vars"
	_ "modernc.org/sqlite"
)

func OpenSQLite() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "file:"+vars.RemoteSQLiteDb)
	if err != nil {
		return nil, fmt.Errorf("opening SQLite db: %v", err)
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
