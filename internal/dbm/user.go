package dbm

import (
	"database/sql"
	"fmt"
	"gamesync/internal/vars"
	"os"
	"path"
	"strconv"
)

type User struct {
	ID int
	RoleID int
	Name string
}
type UserWithRoleName struct {
	ID int
	RoleID int
	RoleName string
	Name string
}
type UserWithRole struct {
	ID int
	Name string
	Role Role
}

/*
Add new user, needs role id and name
*/
func UserAdd(db *sql.DB, user User) error {
	if user.Name == "" {
		return fmt.Errorf("user name can not be empty")
	}
	const SQL = `INSERT INTO user (user_name, role_id) VALUES (?, ?) RETURNING user_id`
	var userID int
	err := db.QueryRow(SQL, user.Name, user.RoleID).Scan(&userID)
	if err != nil {
		return fmt.Errorf("inserting new user: %v", err)
	}
	
	d := path.Join(vars.RemoteSaveDir, strconv.Itoa(userID))
	if err := os.Mkdir(d, 0775); err != nil {
		return fmt.Errorf("creating user save dir: %w", err)
	}
	if err := os.Chown(d, vars.RemoteUID, -1); err != nil {
		return fmt.Errorf("chaning owner of user save dir: %w", err)
	}
	return nil
}

/*
Gets user id, name, and role id from name.
*/
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
/*
Gets user using their name. 
Returns user with their role and perms.
*/
func UserGetWithRole(db *sql.DB, name string) (*UserWithRole, error) {
	user, err := UserGet(db, name)
	if err != nil {
		return nil, fmt.Errorf("getting user %s: %w", name, err)
	}

	role, err := RoleGetWithPerms(db, user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("getting role for user: %s: %w", user.Name, err)
	}

	userWithRole := UserWithRole{
		Name: user.Name,
		ID: user.ID,
		Role: role,
	}
	return &userWithRole, nil
}

/*
Gets user using an id. Includes id, name and role id.
*/
func UserGetFromID(db *sql.DB, id int) (*User, error) {
	const SQL = `SELECT user_id, role_id, user_name FROM user WHERE user_id = ?`
	row := db.QueryRow(SQL, id)
	user := User{}
	if err := row.Scan(&user.ID, &user.RoleID, &user.Name); err != nil {
		return nil, fmt.Errorf("scanning row: %w", err)
	}
	return &user, nil
}
/*
Sets role for a user by their ids
*/
func UserChangeRole(db *sql.DB, userID int, roleID int) error {
	SQL := `UPDATE OR FAIL user SET role_id = ? WHERE user_id = ?`
	if _, err := db.Exec(SQL, roleID, userID); err != nil {
		return err
	}
	return nil
}
/*
Sets role for user, takes in the users name and the new role name.
*/
func UserChangeRoleSimple(db *sql.DB, username string, roleName string) error {
	user, err := UserGet(db, username)
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

/*
Returns all users with their id, name, role_id and role_name
*/
func UserGetAll(db *sql.DB) ([]UserWithRoleName, error) {
	const SQL = `SELECT user_id, role_id, role_name, user_name FROM user JOIN role USING (role_id)`
	rows, err := db.Query(SQL)
	if err != nil {
		return nil, fmt.Errorf("selecting all users: %w", err)
	}
	var users []UserWithRoleName
	for rows.Next() {
		user := UserWithRoleName{}
		if err := rows.Scan(&user.ID, &user.RoleID, &user.RoleName, &user.Name); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}


/*
Deletes all users with a specific roleID and returns the amount of users deleted.
*/
func UserDeleteAllInRole(db *sql.DB, roleID int) (int, error) {
	var qty int
	row := db.QueryRow(`SELECT COUNT(*) FROM user WHERE role_id = ?`, roleID)
	if err := row.Scan(&qty); err != nil {
		return 0, fmt.Errorf("getting use count of users with role: %w", err)
	}

	if _, err := db.Exec(`DELETE FROM user where role_id = ?`, roleID); err != nil {
		return 0, err
	}
	return qty, nil
}

func UserDelete(db *sql.DB, userID int) error {
	if _, err := db.Exec(`DELETE FROM user WHERE user_id = ?`, userID); err != nil {
		return fmt.Errorf("deleting user with id %d: %w", userID, err)
	}
	return nil
}
