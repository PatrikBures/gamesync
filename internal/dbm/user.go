package dbm

import (
	"database/sql"
	"fmt"
)

type UserRoles int
const (
	UserRoleAdmin UserRoles = iota
	UserRoleUser
)

type UserRole struct {
	ID int
	RoleName string
}

type User struct {
	ID int
	UserName string
	UserRoleId UserRoles
}

func AddUser(db *sql.DB, user User) error {
	if user.UserName == "" {
		return fmt.Errorf("user name can not be empty")
	}
	SQL := `INSERT INTO user (user_name, user_role_id) VALUES (?, ?)`
	if _, err := db.Exec(SQL, user.UserName, user.ID); err != nil {
		return fmt.Errorf("inserting new user: %v", err)
	}
	return nil
}

func AddUserSimple(user User) error {
	db, err := OpenSQLite()
	if err != nil {
		return err
	}
	defer CloseDB(db, &err)

	if err := AddUser(db, user); err != nil {
		return fmt.Errorf("creating user: %v", err)
	}
	return nil
}

func AddUserRole(db *sql.DB, role UserRole) error {
	if role.RoleName == "" {
		return fmt.Errorf("role name can not be empty")
	}
	SQL := `INSERT INTO user_role (user_role_id, role_name) VALUES (?, ?)`
	if _, err := db.Exec(SQL, role.ID, role.RoleName); err != nil {
		return fmt.Errorf("inserting new role: %v", err)
	}
	return nil
}
