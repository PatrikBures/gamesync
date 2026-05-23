package dbx

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/libtnb/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDb(dbType string, dsn string, logLevel logger.LogLevel) (*gorm.DB, error) {
	var err error
	var db *gorm.DB

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel: logLevel,
			IgnoreRecordNotFoundError: true,
		},
	)

	config := &gorm.Config{
		TranslateError: true, 
		Logger: newLogger,
	}

	switch dbType{
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), config)
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), config)
	default:
		return nil, fmt.Errorf("unsupported database type '%s'", dbType)
	}

	return db, err
}
