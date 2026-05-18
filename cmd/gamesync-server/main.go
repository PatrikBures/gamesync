package main

import (
	"flag"
	"fmt"
	"gamesync/internal/dbx"
	api "gamesync/internal/ogen"
	"gamesync/internal/query"
	"gamesync/internal/server"
	serverConfig "gamesync/internal/server/config"
	initServer "gamesync/internal/server/initialize"
	middlewares "gamesync/internal/server/middleware"
	"gamesync/internal/server/service"
	"log"
	"log/slog"
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
	defaultRoleID int
}

func start() error {
	c := config{}
	serverConfig.AddStringVar(&c.appDir, "app-dir", server.AppDir, "Path where files will be kept like the SQLite database")
	serverConfig.AddStringVar(&c.dbType, "db-type", "sqlite", "Either 'postgres' or 'sqlite'")
	serverConfig.AddStringVar(&c.dbUrl, "db-url", server.DefaultSQLitePath, "Url to postgres db or path to SQLite db-file")
	serverConfig.AddStringVar(&c.disabledRoles, "disabled-roles", "", "Default roles to disable, seperated by '|'")
	serverConfig.AddIntVar(&c.defaultRoleID, "default-role-id", 50, "Default role id for newly created users. Role needs to exist for creation of new users to work")

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

	if err := initServer.EnsurePermissions(q); err != nil {
		return err
	}
	if err := initServer.CreateDefaultRoles(q, serverConfig.StringToSlice(c.disabledRoles)); err != nil {
		return err
	}
	if err := initServer.CreateDefaultRolePerms(q); err != nil {
		return err
	}
	if token, err := initServer.CreateAdmin(q); err == nil {
		if token == "" {
			slog.Info("admin already exists")
		} else {
			slog.Info("Created admin, make sure to update the token", "token", token)
		}
	}

	s := service.NewService(q, service.ServiceOpts{
		DefaultRoleID: int32(c.defaultRoleID),
	})

	srv, err := api.NewServer(
		s,
		s,
		api.WithMiddleware(
			middlewares.AuthzMiddleware(),
		),
		api.WithPathPrefix("/api/v1"),
	)
	if err != nil {
		return fmt.Errorf("creating server: %v", err)
	}

	return http.ListenAndServe(":8080", srv)
}
