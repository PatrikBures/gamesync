package dbx

import (
	"fmt"
	"os"
	"slices"

	_ "modernc.org/sqlite"
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

	if dbType == "" || dbType == "sqlite" {
		dbType = "sqlite"
		if dbUrl == "" {
			dbUrl = defaultSQLitePath
		}
	} else if !slices.Contains([]string{"postgres"}, dbType) {
		return nil, fmt.Errorf("invalid database type '%s'", dbType)
	}

	engine, err := xorm.NewEngine(dbType, dbUrl)
	if err != nil {
		return nil, err
	}

	return engine, nil
}
