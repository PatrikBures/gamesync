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
	dbPrimaryUrl    string
	dbReplicaUrls   string
	disabledRoles   string
	defaultRoleID   int
	requestLogs     bool
}

func start() error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	c := config{}
	serverConfig.AddStringVar(&c.dbPrimaryUrl, "db-primary-url", "", "Url to primary postgres db, if you only have no replicas use this")
	serverConfig.AddStringVar(&c.dbReplicaUrls, "db-replica-urls", "", "Urls to replica postgres dbs, seperated by '|")
	serverConfig.AddStringVar(&c.disabledRoles, "disabled-roles", "", "Default roles to disable, seperated by '|'")
	serverConfig.AddIntVar(&c.defaultRoleID, "default-role-id", 50, "Default role id for newly created users. Role needs to exist for creation of new users to work")
	serverConfig.AddBoolVar(&c.requestLogs, "request-logs", false, "Enable logs for requests")

	flag.Parse()

	if err := validateConfig(&c); err != nil {
		return fmt.Errorf("validating config: %w", err)
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


	s := service.NewService(db, service.ServiceOpts{
		DefaultRoleID: int32(c.defaultRoleID),
	})


	mw := []api.Middleware{
		middlewares.AuthzMiddleware(),
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
	return nil
}
