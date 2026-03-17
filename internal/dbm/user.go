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
		return fmt.Errorf("UserName can not be empty")
	}

	SQL := `INSERT INTO user (user_name, user_role_id) VALUES (?, ?)`
	_, err := db.Exec(SQL, user.UserName, user.ID)
	if err != nil {
		return fmt.Errorf("inserting new user: %v", err)
	}

	return nil
}
