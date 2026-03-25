package cmdWrapper

import (
	"flag"
	"fmt"
	"gamesync/internal/dbm"
	"gamesync/internal/rrsync"
	"gamesync/internal/util"
	"gamesync/internal/vars"
	"os"
	"path"
	"strconv"
	"syscall"
)

func Execute() error {
	return m()
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

	cmdParsed := util.ParseArgs(cmd)

	switch cmdParsed[0] {
	case "rsync":
		rDir := path.Join(vars.RemoteSaveDir, strconv.Itoa(user.ID))
		rrsync.Run(rDir, cmdParsed[1:])
		return nil
	case "admin":
		cmdCombinded := make([]string, 0, len(cmdParsed))
		cmdCombinded = append(cmdCombinded, "/usr/local/bin/gamesync-admin")
		cmdCombinded = append(cmdCombinded, cmdParsed[1:]...)

		if err := syscall.Exec("/usr/local/bin/gamesync-admin", cmdCombinded, os.Environ()); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid option cmd '%s'", cmdParsed[0])
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

