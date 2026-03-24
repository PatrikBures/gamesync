package main

import (
	cmdAdmin "gamesync/cmd/gamesync-server/admin"
	cmdAuth "gamesync/cmd/gamesync-server/auth"
	cmdWrapper "gamesync/cmd/gamesync-server/wrapper"
	"log"
	"os"
	"path/filepath"
)

var cmds = map[string]func() error {
	"gamesync-admin":    cmdAdmin.Execute,
	"gamesync-auth":     cmdAuth.Execute,
	"gamesync-wrapper":  cmdWrapper.Execute,
}
func main() {
	base := filepath.Base(os.Args[0])
	if _, ok := cmds[base]; !ok {
		log.Fatalf("'%s' is not a valid command", base)
	}
	if err := cmds[base](); err != nil {
		log.Fatalln(err)
	}
}
