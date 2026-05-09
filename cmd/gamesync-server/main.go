package main

import (
	"flag"
	"fmt"
	"gamesync/internal/dbx"
	api "gamesync/internal/ogen"
	"gamesync/internal/query"
	"gamesync/internal/server"
	serverConfig "gamesync/internal/server/config"
	"gamesync/internal/server/roles"
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

type config struct {
	appDir   string
	dbType   string
	dbUrl    string
	disabledRoles string
}

func start() error {
	c := config{}
	serverConfig.AddStringVar(&c.appDir, "app-dir", server.AppDir, "Path where files will be kept like the SQLite database")
	serverConfig.AddStringVar(&c.dbType, "db-type", "sqlite", "Either 'postgres' or 'sqlite'")
	serverConfig.AddStringVar(&c.dbUrl, "db-url", server.DefaultSQLitePath, "Url to postgres db or path to SQLite db-file")
	serverConfig.AddStringVar(&c.disabledRoles, "disabled-roles", "", "Default roles to disable, seperated by '|'")

	flag.Parse()

	for _, dir := range []string{c.appDir} {
		if err := os.MkdirAll(dir, 0775); err != nil {
			return fmt.Errorf("could not initialize dirs: %v", err)
		}
	}

	db, err := dbx.ConnectDb(c.dbType, c.dbUrl)
	if err != nil {
		return err
	}
	if c.dbType == "sqlite" {
		db.Exec("PRAGMA foreign_keys = ON")
	}
	q := query.Use(db)

	if err := roles.CreateDefaultRoles(q, serverConfig.StringToSlice(c.disabledRoles)); err != nil {
		return err
	}

	s := service.NewService(q)

	srv, err := api.NewServer(s, s, api.WithPathPrefix("/api/v1"))
	if err != nil {
		return fmt.Errorf("creating server: %v", err)
	}

	return http.ListenAndServe(":8080", srv)
}
