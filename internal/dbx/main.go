package dbx

import (
	"fmt"
	"os"
	"slices"

	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

type User struct {
	Id int64
	Name string `xorm:"varchar 25 not null unique 'usr_name"`
}

func ConnectDb() (*xorm.Engine, error) {
	dbType := os.Getenv("GAMESYNC_DB_TYPE")
	var dbPath string

	if dbType == "" || dbType == "sqlite" {
		dbType = "sqlite"
		dbPath = "./db.sqlite"
	} else {
		dbPath = os.Getenv("GAMESYNC_DB_HOST")
		if !slices.Contains([]string{"postgres"}, dbType) {
			return nil, fmt.Errorf("invalid database type '%s'", dbType)
		}
	}

	engine, err := xorm.NewEngine(dbType, dbPath)
	if err != nil {
		return nil, err
	}

	return engine, nil
}
