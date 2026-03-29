package dbm

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"
)

type Role struct {
	ID int
	Name string
	Permissions []Permission
}
func (r *Role) HasPermission(perm Permission) bool {
	return slices.Contains(r.Permissions, perm)
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

func RoleSetWithID(db *sql.DB, id int, name string) error {
	if name == "" {
		return fmt.Errorf("role name can not be empty")
	}
	const SQL = `INSERT INTO role (role_id, role_name) VALUES (?,?) ON CONFLICT (role_id) DO UPDATE SET role_name = EXCLUDED.role_name`
	if _, err := db.Exec(SQL, id, name); err != nil {
		return fmt.Errorf("inserting new role: %v", err)
	}
	return nil
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

func RoleAllPerms() []Permission {
	perms := make([]Permission, 0, len(permissionNames))
	for perm := range permissionNames {
		perms = append(perms, perm)
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

func RoleGetPerms(db *sql.DB, roleID int) ([]Permission, error) {
	const SQL = `SELECT permission_id FROM role_permission WHERE role_id = ?`
	rows, err := db.Query(SQL, roleID)
	if err != nil {
		return nil, fmt.Errorf("selecting permissions: %w", err)
	}
	var perms []Permission
	for rows.Next() {
		var id Permission
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		perms = append(perms, id)
	}
	return perms, nil
}

func RoleAddPermsWithIDs(db *sql.DB, roleID int, permIDs []Permission) error {
	if len(permIDs) == 0 {
		return nil
	}
	args := strings.Repeat(`(?,?),`, len(permIDs)-1) + `(?,?)`
	SQL := "INSERT INTO role_permission (permission_id, role_id) VALUES " + args + " ON CONFLICT (role_id, permission_id) DO NOTHING"
	roleIDs := make([]int, 0, len(permIDs))
	permIDsInt := make([]int, 0, len(permIDs))
	for i := range len(permIDs) {
		roleIDs = append(roleIDs, roleID)
		permIDsInt = append(permIDsInt, int(permIDs[i]))
	}
	if _, err := db.Exec(SQL, FlattenArgs(permIDsInt, roleIDs)...); err != nil {
		return err
	}
	return nil
}
func RoleAddPerms(db *sql.DB, roleID int, perms []string) error {
	if len(perms) == 0 {
		return fmt.Errorf("no perm names provided")
	}
	permIDs, err := PermNamesToIDs(db, perms)
	if err != nil {
		return fmt.Errorf("converting perm names to ids: %w", err)
	}
	if len(permIDs) != len(perms) {
		return fmt.Errorf("not all perms provided exist, make sure you provided the correct names and that there are no dublicates. Provided %d, returned %d", len(perms), len(permIDs))
	}
	args := strings.Repeat(`(?,?),`, len(permIDs)-1) + `(?,?)`
	SQL := "INSERT INTO role_permission (permission_id, role_id) VALUES " + args + " ON CONFLICT (role_id, permission_id) DO NOTHING"
	roleIDs := make([]int, 0, len(permIDs))
	for range len(permIDs) {
		roleIDs = append(roleIDs, roleID)
	}
	if _, err := db.Exec(SQL, FlattenArgs(permIDs, roleIDs)...); err != nil {
		return err
	}
	return nil
}

func RoleRemovePerms(db *sql.DB, roleID int, perms []string) error {
	if len(perms) == 0 {
		return fmt.Errorf("no perm names provided")
	}
	permIDs, err := PermNamesToIDs(db, perms)
	if err != nil {
		return fmt.Errorf("converting perm names to ids: %w", err)
	}
	if len(permIDs) != len(perms) {
		return fmt.Errorf("not all perms provided exist, make sure you provided the correct names and that there are no dublicates. Provided %d, returned %d", len(perms), len(permIDs))
	}
	argPlaceholders := "(" + strings.Repeat("?,", len(permIDs)-1) + "?)"
	SQL := "DELETE FROM role_permission WHERE role_id = ? AND permission_id IN " + argPlaceholders
	
	args := make([]any, 0, len(permIDs)+1)
	args = append(args, roleID)
	flattenedPermIDs := FlattenArgs(permIDs)
	args = append(args, flattenedPermIDs...)
	if _, err := db.Exec(SQL, args...); err != nil {
		return err
	}

	return nil
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
