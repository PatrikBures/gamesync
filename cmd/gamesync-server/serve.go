package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	api "go.pabu.dev/gamesync/internal/ogen"
	config "go.pabu.dev/gamesync/internal/server/config"
	"go.pabu.dev/gamesync/internal/server/garbage"
	initServer "go.pabu.dev/gamesync/internal/server/initialize"
	middlewares "go.pabu.dev/gamesync/internal/server/middleware"
	"go.pabu.dev/gamesync/internal/server/service"
	"go.pabu.dev/gamesync/internal/server/storage"

	"go.pabu.dev/bytesize"
)

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
	gcEnabled        bool
	gcCron           string
}

func serve() error {
	appCtx, stopApp := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopApp()


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
	config.AddBoolVar(&c.gcEnabled, "gc-enabled", true, "Enables background goroutine for garbage collector")
	config.AddStringVar(&c.gcCron, "gc-cron", "0 2 * * *", "Cron schedule for garbage collector")

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
		appCtx,
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

	var cr *cron.Cron
	if c.gcEnabled {
		garbageCollector := garbage.New(db)
		cr = cron.New()

		if _, err := cr.AddFunc(c.gcCron, func() {
			jobCtx, cancelJob := context.WithTimeout(appCtx, 15*time.Minute)
			defer cancelJob()

			slog.Info("Started GC job")
			if err := garbageCollector.CleanOrphanedChunks(jobCtx); err != nil {
				if jobCtx.Err() == context.Canceled {
					slog.Error("GC job canceled due to application shutdown")
				}
				slog.Error("GC job error", "error", err)
			}
		}); err != nil {
			return fmt.Errorf("adding cron job: %w", err)
		}
		cr.Start()
		slog.Info("GC cron scheduler started", "schedule", c.gcCron)
	}

	httpSrv := &http.Server{
		Addr:         ":8080",
		Handler:      srv,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	serverErrChan := make(chan error, 1)
	go func() {
		slog.Info("starting HTTP server", "addr", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrChan <- err
		}
	}()

	select {
	case err := <- serverErrChan:
		return fmt.Errorf("HTTP server error: %w", err)
	case <- appCtx.Done():
		slog.Info("shutdown signal received, starting graceful teardown...")
	}

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server forced shutdown", "error", err)
	}

	if cr != nil {
		slog.Info("stopping cron scheduler...")
		cronCtx := cr.Stop()
		<-cronCtx.Done()
	}

	slog.Info("server stopped cleanly")
	return nil
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
