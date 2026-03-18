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
	Name string
	RoleID int
}

func UserAdd(db *sql.DB, user User) error {
	if user.Name == "" {
		return fmt.Errorf("user name can not be empty")
	}
	SQL := `INSERT INTO user (user_name, user_role_id) VALUES (?, ?)`
	if _, err := db.Exec(SQL, user.Name, user.ID); err != nil {
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

func RoleAdd(db *sql.DB, role Role) error {
	if role.Name == "" {
		return fmt.Errorf("role name can not be empty")
	}
	SQL := `INSERT INTO user_role (user_role_id, role_name) VALUES (?, ?)`
	if _, err := db.Exec(SQL, role.ID, role.Name); err != nil {
		return fmt.Errorf("inserting new role: %v", err)
	}
	return nil
}
