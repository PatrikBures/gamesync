package dbm

import (
	"database/sql"
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
	PermRoleCreate
	PermRoleRemove
	PermRoleChangePerms
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
	PermRoleCreate:        "role_create",
	PermRoleRemove:        "role_remove",
	PermRoleChangePerms:   "role_change_perms",
}
func (p Permission) String() string {
	if name, ok := permissionNames[p]; ok {
		return name
	}
	panic(fmt.Errorf("permission with id %d not mapped", p))
}

type Role struct {
	ID int
	Name string
	Permissions map[Permission]bool
}
func (r *Role) HasPermission(perm Permission) bool {
	return r.Permissions[perm]
}

type User struct {
	ID int
	RoleID int
	Name string
}

func UserAdd(db *sql.DB, user User) error {
	if user.Name == "" {
		return fmt.Errorf("user name can not be empty")
	}
	SQL := `INSERT INTO user (user_name, role_id) VALUES (?, ?)`
	if _, err := db.Exec(SQL, user.Name, user.RoleID); err != nil {
		return fmt.Errorf("inserting new user: %v", err)
	}
	return nil
}

func UserAddSimple(user User) error {
	db, err := OpenSQLite()
	if err != nil {
		return err
	}
	defer CloseDB(db, &err)

	if err := UserAdd(db, user); err != nil {
		return fmt.Errorf("creating user: %v", err)
	}
	return nil
}


func UserGet(db *sql.DB, name string) (*User, error) {
	SQL := `SELECT user_id, role_id, user_name FROM user WHERE user_name = ?`
	row := db.QueryRow(SQL, name)
	if row.Err() != nil {
		return nil, row.Err()
	}
	user := User{}
	row.Scan(&user.ID, &user.RoleID, &user.Name)
	return &user, nil
}

func UserChangeRole(db *sql.DB, userID int, roleID int) error {
	SQL := `UPDATE OR FAIL user SET role_id = ? WHERE user_id = ?`
	if _, err := db.Exec(SQL, roleID, userID); err != nil {
		return err
	}
	return nil
}
func UserChangeRoleSimple(userName string, roleName string) error {
	db, err := OpenSQLite()
	if err != nil {
		return err
	}
	defer CloseDB(db, &err)

	user, err := UserGet(db, userName)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}

	roleID, err := RoleGetID(db, roleName)
	if err != nil {
		return fmt.Errorf("getting role id: %w", err)
	}

	if err := UserChangeRole(db, user.ID, roleID); err != nil {
		return fmt.Errorf("changing role for user: %w", err)
	}

	return nil
}


func RoleGetID(db *sql.DB, name string) (int, error) {
	SQL := `SELECT role_id FROM role WHERE role_name = ?`
	row := db.QueryRow(SQL, name)
	if row.Err() != nil {
		return -1, row.Err()
	}
	var id int
	row.Scan(&id)
	return id, nil
}

func RoleAddSimple(db *sql.DB, name string) error {
	if name == "" {
		return fmt.Errorf("role name can not be empty")
	}
	SQL := `INSERT INTO role (role_name) VALUES (?)`
	if _, err := db.Exec(SQL, name); err != nil {
		return fmt.Errorf("inserting new role: %v", err)
	}
	return nil
}

func RoleAddWithPerms(db *sql.DB, role Role) error {
	if role.Name == "" {
		return fmt.Errorf("role name can not be empty")
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	const addRoleSQL = `INSERT INTO role (role_id, role_name) VALUES (?, ?)`
	if _, err := tx.Exec(addRoleSQL, role.ID, role.Name); err != nil {
		return fmt.Errorf("inserting new role: %v", err)
	}

	if len(role.Permissions) > 0 {
		valueStrings := make([]string, 0, len(role.Permissions))
		valueArgs := make([]any, 0, len(role.Permissions)*2)
		for p, enabled := range role.Permissions {
			if !enabled {
				continue
			}
			valueStrings = append(valueStrings, "(?, ?)")
			valueArgs = append(valueArgs, p, role.ID)
		}
		if len(valueArgs) > 0 {
			stmt := fmt.Sprintf("INSERT INTO role_permission (permission_id, role_id) VALUES %s", strings.Join(valueStrings, ","))
			if _, err := tx.Exec(stmt, valueArgs...); err != nil {
				return fmt.Errorf("inserting roles: %s: %w", stmt, err)
			}
		}
	}
	tx.Commit()
	return nil
}

func RoleAllPerms(enabled bool) map[Permission]bool {
	perms := make(map[Permission]bool, len(permissionNames))
	for perm := range permissionNames {
		perms[perm] = enabled
	}
	return perms
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
	for true {
		if !rows.Next() {
			break
		}
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
