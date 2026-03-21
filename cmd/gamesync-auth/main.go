package main

import (
	"fmt"
	"gamesync/internal/dbm"
	"log"
	"os"
)

func main() {
	if err := printKey(); err != nil {
		log.Fatalln(err)
	}
}

func printKey() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("wrong number of args, wants only 1")
	}
	fingerprint := os.Args[1]

	db, err := dbm.OpenSQLite()
	if err != nil { return err }
	defer dbm.CloseDB(db, &err)

	keys, err := dbm.KeyGetKeysByFingerprint(db, fingerprint)
	if err != nil {
		return fmt.Errorf("getting keys: %w", err)
	}
	for _, k := range keys {
		fmt.Print(k)
	}

	return nil
}
