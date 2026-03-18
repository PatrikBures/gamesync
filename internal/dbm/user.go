package dbm

import (
	"database/sql"
	"fmt"
)


type UserRole struct {
	ID int
	Name string

	PermGameDelete bool
	PermGameDeleteOwn bool
	PermGameRename bool
	PermGameRenameOwn bool
	PermGameAdd bool
	PermUserAdd bool
	PermUserDelete bool
	PermUserRename bool
	PermUserRenameOwn bool
	PermUserChangeRole bool
	PermUserList bool
	PermRoleCreate bool
	PermRoleRemove bool
	PermRoleChangePerms bool
	PermSync bool
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

func RoleAdd(db *sql.DB, role UserRole) error {
	if role.Name == "" {
		return fmt.Errorf("role name can not be empty")
	}
	SQL := `INSERT INTO user_role (user_role_id, role_name) VALUES (?, ?)`
	if _, err := db.Exec(SQL, role.ID, role.Name); err != nil {
		return fmt.Errorf("inserting new role: %v", err)
	}
	return nil
}
