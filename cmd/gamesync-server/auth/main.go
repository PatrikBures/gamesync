package cmdAuth

import (
	"fmt"
	"gamesync/internal/dbm"
	"os"
)

func Execute() error {
	return printKey()
}

func printKey() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("wrong number of args, wants only 1")
	}
	fingerprint := os.Args[1]

	db, err := dbm.OpenSQLite()
	if err != nil { return err }
	defer dbm.CloseDB(db, &err)

	key, userID, err := dbm.KeyGetByFingerprint(db, fingerprint)
	if err != nil {
		return fmt.Errorf("getting keys: %w", err)
	}
	user, err := dbm.UserGetFromID(db, userID)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}
	fmt.Printf(`command="/usr/local/bin/gamesync-wrapper -username %s -userid %d" %s`, user.Name, user.ID, key)
	return nil
}
