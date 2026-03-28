package dbm

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Permission int
const (
	PermSync Permission = iota
	PermGameDelete
	PermGameDeleteOwn
	PermGameRename
	PermGameRenameOwn
	PermGameAdd
	PermUserAdd
	PermUserDelete
	PermUserRename
	PermUserRenameSelf
	PermUserChangeRole
	PermUserList
	PermRoleAdd
	PermRoleDelete
	PermRoleChangePerms
	PermRoleList
	PermRolePermList
	PermRolePermListOwn
	PermRolePermMod
	PermKeyAdd
	PermKeyAddSelf
	PermKeyList
	PermKeyListOwn
)

var permissionNames = map[Permission]string{
	PermSync:              "sync",
	PermGameDelete:        "game_delete",
	PermGameDeleteOwn:     "game_delete_own",
	PermGameRename:        "game_rename",
	PermGameRenameOwn:     "game_rename_own",
	PermGameAdd:           "game_add",
	PermUserAdd:           "user_add",
	PermUserDelete:        "user_delete",
	PermUserRename:        "user_rename",
	PermUserRenameSelf:    "user_rename_self",
	PermUserChangeRole:    "user_change_role",
	PermUserList:          "user_list",
	PermRoleAdd:           "role_add",
	PermRoleDelete:        "role_delete",
	PermRoleChangePerms:   "role_change_perms",
	PermRoleList:          "role_list",
	PermRolePermList:      "role_perm_list",
	PermRolePermListOwn:   "role_perm_list_own",
	PermRolePermMod:       "role_perm_mod",
	PermKeyAdd:            "key_add",
	PermKeyAddSelf:        "key_add_self",
	PermKeyList:           "key_list",
	PermKeyListOwn:        "key_list_own",
}

func (p Permission) String() string {
	if name, ok := permissionNames[p]; ok {
		return name
	}
	panic(fmt.Errorf("permission with id %d not mapped", p))
}

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
		fmt.Println("added perm:", p,)
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

func PermNamesToIDs(db *sql.DB, perms []string) (permIDs []int, err error) {
	SQL := `SELECT permission_id FROM permission WHERE permission_name IN ` + listSQL(len(perms))

	args := make([]any, len(perms))
	for i, p := range perms {
		args[i] = p
	}

	rows, err := db.Query(SQL, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			err = errors.Join(err, cerr)
		}
	}()

	permIDs = make([]int, 0, len(perms))
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		permIDs = append(permIDs, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return permIDs, nil
}
