package dbm

import (
	"database/sql"
	"fmt"
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

func RoleAdd(db *sql.DB, name string) error {
	if name == "" {
		return fmt.Errorf("role name can not be empty")
	}
	SQL := `INSERT INTO role (role_name) VALUES (?)`
	if _, err := db.Exec(SQL, name); err != nil {
		return fmt.Errorf("inserting new role: %v", err)
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
