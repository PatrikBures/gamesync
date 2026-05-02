package dbx

import (
	"fmt"
	"os"

	_ "modernc.org/sqlite"
	_ "github.com/jackc/pgx/v5/stdlib"
	"xorm.io/xorm"
)

const defaultSQLitePath = "/db.sqlite"

type User struct {
	Id int64
	Name string `xorm:"varchar 25 not null unique 'usr_name"`
}

func ConnectDb() (*xorm.Engine, error) {
	dbType := os.Getenv("GAMESYNC_DB_TYPE")
	dbUrl := os.Getenv("GAMESYNC_DB_URL")

	switch dbType{
	case "", "sqlite":
		dbType = "sqlite"
		if dbUrl == "" {
			dbUrl = defaultSQLitePath
		}
	case "postgres":
		dbType = "pgx"
	default:
		return nil, fmt.Errorf("invalid database type '%s'", dbType)
	}

	engine, err := xorm.NewEngine(dbType, dbUrl)
	if err != nil {
		return nil, err
	}

	return engine, nil
}
