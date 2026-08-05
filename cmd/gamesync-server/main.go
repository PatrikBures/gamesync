package main

import (
	"context"
	"flag"
	"fmt"
	api "go.pabu.dev/gamesync/internal/ogen"
	config "go.pabu.dev/gamesync/internal/server/config"
	initServer "go.pabu.dev/gamesync/internal/server/initialize"
	middlewares "go.pabu.dev/gamesync/internal/server/middleware"
	"go.pabu.dev/gamesync/internal/server/service"
	"go.pabu.dev/gamesync/internal/server/storage"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"go.pabu.dev/bytesize"
)

func main() {
	if err := start(); err != nil {
		log.Fatal(err)
	}
}

type opts struct {
	dbPrimaryUrl     string
	dbReplicaUrls    string
	disabledRoles    string
	defaultRoleID    int
	requestLogs      bool
	chunkBaseDir     string
	chunkTmpDir      string
	maxChunkSize     string
	quietLogRequests string
}

func start() error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	c := opts{}
	config.AddStringVar(&c.dbPrimaryUrl, "db-primary-url", "", "Url to primary postgres db, if you only have no replicas use this")
	config.AddStringVar(&c.dbReplicaUrls, "db-replica-urls", "", "Urls to replica postgres dbs, seperated by '|")
	config.AddStringVar(&c.disabledRoles, "disabled-roles", "", "Default roles to disable, seperated by '|'")
	config.AddStringVar(&c.chunkBaseDir, "chunk-dir", "/var/lib/gamesync/chunks", "Location where chunks are stored")
	config.AddIntVar(&c.defaultRoleID, "default-role-id", 50, "Default role id for newly created users. Role needs to exist for creation of new users to work")
	config.AddStringVar(&c.maxChunkSize, "max-chunk-size", "256Ki", "Max chunk size")
	config.AddBoolVar(&c.requestLogs, "log-requests", false, "Enable logs for requests")
	config.AddStringVar(&c.quietLogRequests, "log-requests-quiet", "GetHealth", "Hide requests from logger, seperated by '|'")

	flag.Parse()

	if err := validateOpts(&c); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}

	for _, d := range []string{
		c.chunkBaseDir,
	} {
		if err := os.MkdirAll(d, 0775); err != nil {
			return err
		}
		slog.Info("ensured dir exists", "path", d)
	}

	db, err := initServer.InitDatabase(
		context.Background(),
		c.dbPrimaryUrl,
		config.StringToSlice(c.dbReplicaUrls),
		config.StringToSlice(c.disabledRoles),
	)
	if err != nil {
		return fmt.Errorf("initializing database: %w", err)
	}

	maxChunkBytes, err := bytesize.Parse(c.maxChunkSize)
	if err != nil {
		return fmt.Errorf("parsing max chunks size: %w", err)
	}
	slog.Info("max chunk size", "size", c.maxChunkSize, "bytes", maxChunkBytes)
	localStorage, err := storage.NewLocal(c.chunkBaseDir, maxChunkBytes)
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
		mw = append(mw, middlewares.Logging(requestLogger, config.StringToMap(c.quietLogRequests)))
	}

	mw = append(mw,
		middlewares.LoadPathData(db),
	)

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

func validateOpts(c *opts) error {
	if c.dbPrimaryUrl == "" {
		return fmt.Errorf("primary url needs to be set")
	}
	if c.chunkBaseDir == "" {
		return fmt.Errorf("chunk dir needs to be set")
	}

	c.chunkTmpDir = filepath.Join(c.chunkBaseDir, ".tmp")

	return nil
}
