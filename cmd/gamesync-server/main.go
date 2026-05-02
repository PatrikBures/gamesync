package main

import (
	"fmt"
	"gamesync/internal/dbx"
	"gamesync/internal/query"
	"gamesync/internal/server"
	"gamesync/internal/server/api"
	"log"
	"os"
)

func main() {
	if err := start(); err != nil {
		log.Fatal(err)
	}
}

func start() error {
	if err := server.CreateDirs(); err != nil {
		return fmt.Errorf("could not initialize dirs: %v", err)
	}

	dbType := os.Getenv("GAMESYNC_DB_TYPE")
	if dbType == "" {
		dbType = "sqlite"
	}
	dbUrl := os.Getenv("GAMESYNC_DB_URL")
	if dbUrl == "" {
		dbUrl = server.DefaultSQLitePath
	}

	db, err := dbx.ConnectDb(dbType, dbUrl)
	if err != nil {
		return err
	}
	if dbType == "sqlite" {
		db.Exec("PRAGMA foreign_keys = ON")
	}
	q := query.Use(db)

	handlerOpts := api.HandlerOpts{
		Logging: false,
	}
	if os.Getenv("GAMESYNC_LOGGING") == "true" {
		handlerOpts.Logging = true
	}
	return api.NewHandler(handlerOpts, q).Serve()
}
