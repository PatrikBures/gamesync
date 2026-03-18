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

	const selectPermsSQL = `SELECT permission_id FROM role_permission WHERE role_id = ?`
	rows, err := db.Query(selectPermsSQL, roleID)
	if err != nil {
		return role, fmt.Errorf("selecting permissions: %w", err)
	}
	role.Permissions = make(map[Permission]bool)
	for rows.Next() {
		var id Permission
		if err := rows.Scan(&id); err != nil {
			return role, fmt.Errorf("scanning row: %w", err)
		}
		role.Permissions[id] = true
	}

	const selectRoleNameSQL = `SELECT role_name FROM role WHERE role_id = ?`
	row := db.QueryRow(selectRoleNameSQL, roleID)
	if row.Err() != nil {
		return role, fmt.Errorf("selecting role name: %w", err)
	}
	if err := row.Scan(&role.Name); err != nil {
		return role, fmt.Errorf("scanning row: %w", err)
	}
	
	return role, nil
}
