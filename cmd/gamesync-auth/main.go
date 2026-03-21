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


	// write fingerprint to tmp, for testing
	f, err := os.CreateTemp("", "fingerprint_")
	if err != nil {
		log.Fatalln("create tmp file")
	}
	defer f.Close()
	_, _ = fmt.Fprintln(f, "1")


	db, err := dbm.OpenSQLite()
	if err != nil { 
		_, _ = fmt.Fprintln(f, err)
		return err 
	}
	defer func() {
		_, _ = fmt.Fprintln(f, err)
		dbm.CloseDB(db, &err)
	}()


	keys, err := dbm.KeyGetKeysByFingerprint(db, fingerprint)
	if err != nil {
		_, _ = fmt.Fprintln(f, "3")
		return fmt.Errorf("getting keys: %w", err)
	}
	for _, k := range keys {
		fmt.Print(k)
		_, _ = fmt.Fprint(f, k)
	}



	return nil
}
