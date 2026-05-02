package main

import (
	"fmt"
	"gamesync/internal/dbx"
	"gamesync/internal/server/api"
	"log"
	"log/slog"
	"os"
)

func main() {
	if err := start(); err != nil {
		log.Fatal(err)
	}
}

func start() error {
	engine, err := dbx.ConnectDb()
	if err != nil {
		return err
	}
	defer engine.Close()

	slog.Info("created database", "driverName", engine.DriverName())
	if err := engine.Ping(); err != nil {
		return fmt.Errorf("database not reachable: %v", err)
	}
	
	handlerOpts := api.HandlerOpts{
		Logging: false,
	}
	if os.Getenv("GAMESYNC_LOGGING") == "true" {
		handlerOpts.Logging = true
	}
	return api.NewHandler(handlerOpts, engine).Serve()
}
