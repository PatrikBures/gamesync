package main

import (
	"flag"
	"fmt"
	"gamesync/internal/dbm"
	"log"
	"os"
	"strconv"
	"syscall"
)

func main() {
	if err := m(); err != nil {
		log.Fatalln(err)
	}
}

func m() error {
	user := parseUser()

	if err := os.Setenv("GAMESYNC_USER", user.Name); err != nil {
		return fmt.Errorf("setting GAMESYNC_USER: %w", err)
	}
	if err := os.Setenv("GAMESYNC_USER_ID", strconv.Itoa(user.ID)); err != nil {
		return fmt.Errorf("setting GAMESYNC_USER_ID: %w", err)
	}

	cmd := os.Getenv("SSH_ORIGINAL_COMMAND")

	if err := syscall.Exec("/usr/local/bin/gamesync-admin", []string{cmd}, os.Environ()); err != nil {
		return nil
	}

	return nil
}

func parseUser() *dbm.User {
	user := dbm.User{}

	flag.IntVar(&user.ID, "userid", -1, "Sets the user id for user to run, should match username")
	flag.StringVar(&user.Name, "username", "", "Sets the username for user to run, should match userid")
	flag.Parse()

	return &user
}

