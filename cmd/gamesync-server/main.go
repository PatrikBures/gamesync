package main

import (
	"context"
	"flag"
	"fmt"
	api "gamesync/internal/ogen"
	serverConfig "gamesync/internal/server/config"
	initServer "gamesync/internal/server/initialize"
	middlewares "gamesync/internal/server/middleware"
	"gamesync/internal/server/service"
	"gamesync/internal/server/storage"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	if err := start(); err != nil {
		log.Fatal(err)
	}
}

type config struct {
	dbPrimaryUrl  string
	dbReplicaUrls string
	disabledRoles string
	defaultRoleID int
	requestLogs   bool
	chunkBaseDir  string
	chunkTmpDir   string

}

func start() error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	c := config{}
	serverConfig.AddStringVar(&c.dbPrimaryUrl, "db-primary-url", "", "Url to primary postgres db, if you only have no replicas use this")
	serverConfig.AddStringVar(&c.dbReplicaUrls, "db-replica-urls", "", "Urls to replica postgres dbs, seperated by '|")
	serverConfig.AddStringVar(&c.disabledRoles, "disabled-roles", "", "Default roles to disable, seperated by '|'")
	serverConfig.AddStringVar(&c.chunkBaseDir, "chunk-dir", "/var/lib/gamesync/chunks", "Location where chunks are stored")
	serverConfig.AddIntVar(&c.defaultRoleID, "default-role-id", 50, "Default role id for newly created users. Role needs to exist for creation of new users to work")
	serverConfig.AddBoolVar(&c.requestLogs, "request-logs", false, "Enable logs for requests")

	flag.Parse()

	if err := validateConfig(&c); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}

	for _, d := range []string{
		c.chunkBaseDir,
	} {
		if err :=os.MkdirAll(d, 0775); err != nil {
			return err
		}
		slog.Info("ensured dir exists", "path", d)
	}

	db, err := initServer.InitDatabase(
		context.Background(),
		c.dbPrimaryUrl,
		serverConfig.StringToSlice(c.dbReplicaUrls),
		serverConfig.StringToSlice(c.disabledRoles),
	)
	if err != nil {
		return fmt.Errorf("initializing database: %w", err)
	}

	localStorage, err := storage.NewLocal(c.chunkBaseDir)
	if err != nil {
		return fmt.Errorf("starting new local storage: %w", err)
	}

	s := service.NewService(
		db, localStorage,
		service.ServiceOpts{
			DefaultRoleID: int32(c.defaultRoleID),
		},
	)

	mw := []api.Middleware{
		middlewares.AuthzMiddleware(),
		middlewares.UserAuthz(),
	}

	if c.requestLogs {
		requestLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		mw = append(mw, middlewares.Logging(requestLogger))
	}

	srv, err := api.NewServer(
		s,
		s,
		api.WithMiddleware(mw...),
		api.WithPathPrefix("/api/v1"),
	)
	if err != nil {
		return fmt.Errorf("creating server: %v", err)
	}

	return http.ListenAndServe(":8080", srv)
}

func validateConfig(c *config) error {
	if c.dbPrimaryUrl == "" {
		return fmt.Errorf("primary url needs to be set")
	}
	if c.chunkBaseDir == "" {
		return fmt.Errorf("chunk dir needs to be set")
	}

	c.chunkTmpDir = filepath.Join(c.chunkBaseDir, ".tmp")

	return nil
}
