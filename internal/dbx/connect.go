package dbx

import (
	"fmt"

	"github.com/libtnb/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDb(dbType string, dsn string) (*gorm.DB, error) {
	var err error
	var db *gorm.DB

	switch dbType{
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database type '%s'", dbType)
	}

	return db, err
}
