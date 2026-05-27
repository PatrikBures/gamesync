package main

import (
	"flag"
	"fmt"
	api "gamesync/internal/ogen"
	"gamesync/internal/server"
	serverConfig "gamesync/internal/server/config"
	initServer "gamesync/internal/server/initialize"
	middlewares "gamesync/internal/server/middleware"
	"gamesync/internal/server/service"
	"log"
	"log/slog"
	"net/http"
	"os"

	"gorm.io/gorm/logger"
)

func main() {
	if err := start(); err != nil {
		log.Fatal(err)
	}
}

type config struct {
	appDir          string
	dbType          string
	dbUrl           string
	dbLogLevel      string
	disabledRoles   string
	defaultRoleID   int
	requestLogs     bool
}

func start() error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	c := config{}
	serverConfig.AddStringVar(&c.appDir, "app-dir", server.AppDir, "Path where files will be kept like the SQLite database")
	serverConfig.AddStringVar(&c.dbType, "db-type", "sqlite", "Either 'postgres' or 'sqlite'")
	serverConfig.AddStringVar(&c.dbUrl, "db-url", server.DefaultSQLitePath, "Url to postgres db or path to SQLite db-file")
	serverConfig.AddStringVar(&c.disabledRoles, "disabled-roles", "", "Default roles to disable, seperated by '|'")
	serverConfig.AddIntVar(&c.defaultRoleID, "default-role-id", 50, "Default role id for newly created users. Role needs to exist for creation of new users to work")
	serverConfig.AddBoolVar(&c.requestLogs, "request-logs", false, "Enable logs for requests")
	serverConfig.AddStringVar(&c.dbLogLevel, "db-log-level", "error", "Change log level for db (error, info, warn, silent) default=error")

	flag.Parse()

	for _, dir := range []string{c.appDir} {
		if err := os.MkdirAll(dir, 0775); err != nil {
			return fmt.Errorf("could not initialize dirs: %v", err)
		}
	}

	logLevel, err := parseDbLogLevel(c.dbLogLevel)
	if err != nil {
		return err
	}

	q, err := initServer.InitDatabase(
		c.dbType, c.dbUrl, logLevel,
		serverConfig.StringToSlice(c.disabledRoles),
	)
	if err != nil {
		return fmt.Errorf("initializing database: %w", err)
	}


	s := service.NewService(q, service.ServiceOpts{
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


func parseDbLogLevel(level string) (logger.LogLevel, error) {
	levelMap := map[string]logger.LogLevel {
		"error":   logger.Error,
		"warn":    logger.Warn,
		"info":    logger.Info,
		"silent":  logger.Silent,
	}
	logLevel, ok := levelMap[level]
	if !ok {
		return logger.Error, fmt.Errorf("invalid db log level '%s'. valid levels: error, warn, info, silent", level)
	}
	return logLevel, nil
}
