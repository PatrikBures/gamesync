package dbm

import (
	"database/sql"
	"fmt"
)

type UserTypes int
const (
	UserTypeAdmin UserTypes = iota
	UserTypeUser
)

type UserType struct {
	ID int
	UserTypeName string
}

type User struct {
	ID int
	UserName string
	UserTypeId UserTypes
}

func AddUser(db *sql.DB, user User) error {
	if user.UserName == "" {
		return fmt.Errorf("UserName can not be empty")
	}

	SQL := `INSERT INTO user (user_name, user_type_id) VALUES (?, ?)`
	_, err := db.Exec(SQL, user.UserName, user.ID)
	if err != nil {
		return fmt.Errorf("inserting new user: %v", err)
	}

	return nil
}
