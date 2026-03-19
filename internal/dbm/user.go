package dbm

import (
	"database/sql"
	"fmt"
)

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
	if err := row.Scan(&user.ID, &user.RoleID, &user.Name); err != nil {
		return nil, fmt.Errorf("scanning row: %w", err)
	}
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

func UserGetAll(db *sql.DB) ([]User, error) {
	const SQL = `SELECT user_id, role_id, user_name FROM user`
	rows, err := db.Query(SQL)
	if err != nil {
		return nil, fmt.Errorf("selecting all users: %w", err)
	}
	var users []User
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.RoleID, &user.Name); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

