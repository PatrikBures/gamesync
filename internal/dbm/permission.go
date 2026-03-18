package dbm

import (
	"fmt"
	"strings"
)

func PermsSet() error {
	db, err := OpenSQLite()
	if err != nil {
		return err
	}
	defer CloseDB(db, &err)

	if err != nil {
		return err
	}

	currentPerms := make(map[Permission]string, len(permissionNames))

	const getRolesSQL = `SELECT permission_id, permission_name FROM permission`
	rows, err := db.Query(getRolesSQL)
	if err != nil {
		return fmt.Errorf("selecting current perms: %w", err)
	}
	for rows.Next() {
		var p Permission
		var name string
		if err := rows.Scan(&p, &name); err != nil {
			return fmt.Errorf("scanning current perms: %w", err)
		}
		currentPerms[p] = name
	}

	updatePerms := make([]Permission, 0, len(permissionNames))

	for p, expectedName := range permissionNames {
		currentName, ok := currentPerms[p]
		if ok && currentName == expectedName {
			continue
		}
		updatePerms = append(updatePerms, p)
	}

	if len(updatePerms) == 0 {
		return nil
	}

	tx, err := db.Begin()

	valueStrings := make([]string, 0, len(updatePerms))
	valueArgs := make([]any, 0, len(updatePerms)*2)
	for _, p := range updatePerms{
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, p, permissionNames[p])
		fmt.Printf("added perm:  %d, %s\n", p, permissionNames[p])
	}
	stmt := fmt.Sprintf("INSERT INTO permission (permission_id, permission_name) VALUES %s", strings.Join(valueStrings, ","))
	if _, err := tx.Exec(stmt, valueArgs...); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
