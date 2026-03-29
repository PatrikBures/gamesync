package cmdInit

import (
	"database/sql"
	"errors"
	"fmt"
	"gamesync/internal/dbm"
	"gamesync/internal/vars"
	"os"
)

func Execute() (err error) {
	if err = dirs(); err != nil {
		return fmt.Errorf("dirs: %w", err)
	}

	db, err := dbm.OpenSQLite()
	if err != nil {
		return err
	}
	defer func(){
		if cerr := db.Close(); cerr != nil {
			err = errors.Join(err, cerr)
		}
	}()

	if err = migrate(db); err != nil {
		return fmt.Errorf("migrating schema: %w", err)
	}
	if err = permissions(db); err != nil {
		return fmt.Errorf("permissions: %w", err)
	}
	if err = roles(db); err != nil {
		return fmt.Errorf("roles: %w", err)
	}

	return nil
}

func dirs() error {
	dirs := []string{
		vars.RemoteSaveDir,
		vars.RemoteBackupDir,
		vars.RemoteDbDir,
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		if err := os.Chmod(dir, 0775); err != nil {
			return err
		}
		if err := os.Chown(dir, vars.RemoteUID, -1); err != nil {
			return err
		}
	}
	return nil
}

func migrate(db *sql.DB) error {
	if err := dbm.Migrate(db); err != nil {
		return fmt.Errorf("migrating: %v", err)
	}
	if err := os.Chown(vars.RemoteSQLiteDb, vars.RemoteUID, -1); err != nil {
		return fmt.Errorf("changing owner of db: %w", err)
	}
	return nil
}

func permissions(db *sql.DB) error {
	if err := dbm.PermsSet(db); err != nil {
		return fmt.Errorf("setting perms: %w", err)
	}
	return nil
}

func roles(db *sql.DB) error {
	roles := []dbm.Role{
		{
			ID: 0,
			Name: "none",
			Permissions: []dbm.Permission{},
		},
		{
			ID: 1,
			Name: "admin", 
			Permissions: dbm.RoleAllPerms(),
		},
	}
	for _, role := range roles {
		dbm.RoleAddWithID(db, role.ID, role.Name)
		dbm.RoleAddPermsWithIDs(db, role.ID, role.Permissions)
	}
	
	return nil
}
