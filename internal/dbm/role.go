package dbm

import (
	"database/sql"
	"fmt"
	"strings"
)

type Role struct {
	ID int
	Name string
	Permissions map[Permission]bool
}
func (r *Role) HasPermission(perm Permission) bool {
	return r.Permissions[perm]
}

func RoleGetID(db *sql.DB, name string) (int, error) {
	SQL := `SELECT role_id FROM role WHERE role_name = ?`
	row := db.QueryRow(SQL, name)
	if row.Err() != nil {
		return -1, row.Err()
	}
	var id int
	if err := row.Scan(&id); err != nil {
		return -1, fmt.Errorf("scannig row: %w", err)
	}
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
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commiting to db: %w", err)
	}
	return nil
}

func RoleAllPerms(enabled bool) map[Permission]bool {
	perms := make(map[Permission]bool, len(permissionNames))
	for perm := range permissionNames {
		perms[perm] = enabled
	}
	return perms
}

func RoleGetWithPerms(db *sql.DB, roleID int) (Role, error) {
	role := Role{ID: roleID}

	p, err := RoleGetPerms(db, roleID)
	if err != nil {
		return role, fmt.Errorf("getting perms for role with id %d: %w", roleID, err)
	}
	role.Permissions = p

	const selectRoleNameSQL = `SELECT role_name FROM role WHERE role_id = ?`
	row := db.QueryRow(selectRoleNameSQL, roleID)
	if err := row.Scan(&role.Name); err != nil {
		return role, fmt.Errorf("scanning row: %w", err)
	}
	
	return role, nil
}

func RoleGetPerms(db *sql.DB, roleID int) (map[Permission]bool, error) {
	const SQL = `SELECT permission_id FROM role_permission WHERE role_id = ?`
	rows, err := db.Query(SQL, roleID)
	if err != nil {
		return nil, fmt.Errorf("selecting permissions: %w", err)
	}
	perms := make(map[Permission]bool)
	for rows.Next() {
		var id Permission
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		perms[id] = true
	}
	return perms, nil
}

func RoleGetAll(db *sql.DB) ([]Role, error) {
	rows, err := db.Query(`SELECT role_id FROM role`)
	if err != nil {
		return nil, fmt.Errorf("selecting all roles: %w", err)
	}

	var roles []Role
	for rows.Next() {
		var roleID int
		if err := rows.Scan(&roleID); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		role, err := RoleGetWithPerms(db, roleID)
		if err != nil {
			return nil, fmt.Errorf("getting perms for role with ID %d: %w", roleID, err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func RoleDeleteWithID(db *sql.DB, roleID int) error {
	const SQL = `DELETE FROM role WHERE role_id = ?`
	if _, err := db.Exec(SQL, roleID); err != nil {
		return err
	}
	return nil
}
