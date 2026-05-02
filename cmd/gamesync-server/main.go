package main

import (
	"fmt"
	"gamesync/internal/api"
	"gamesync/internal/dbx"
	"log"
	"log/slog"
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
	
	return api.NewHandler().Serve()
}
