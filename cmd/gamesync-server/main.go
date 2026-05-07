package main

import (
	"fmt"
	"gamesync/internal/dbx"
	api "gamesync/internal/ogen"
	"gamesync/internal/query"
	"gamesync/internal/server"
	"gamesync/internal/server/service"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := start(); err != nil {
		log.Fatal(err)
	}
}

func start() error {
	appDir := os.Getenv("GAMESYNC_APP_DIR")
	if appDir == "" {
		appDir = server.AppDir
	}

	for _, dir := range []string{appDir} {
		if err := os.MkdirAll(dir, 0775); err != nil {
			return fmt.Errorf("could not initialize dirs: %v", err)
		}
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

	s := service.NewService(q)

	srv, err := api.NewServer(s, api.WithPathPrefix("/api/v1"))
	if err != nil {
		return fmt.Errorf("creating server: %v", err)
	}

	return http.ListenAndServe(":8080", srv)
}
