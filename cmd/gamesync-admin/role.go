package main

import (
	"fmt"
	"gamesync/internal/dbm"
	"slices"

	"github.com/spf13/cobra"
)

type roleCmd struct {
	cmd *cobra.Command
}
func newRoleCmd(user *dbm.UserWithRole) *roleCmd {
	root := roleCmd{}
	cmd := &cobra.Command{
		Use: "role",
		Short: "Manage roles",
	}
	if user.Role.HasPermission(dbm.PermRoleAdd) { cmd.AddCommand(newRoleAddCmd().cmd) }
	if user.Role.HasPermission(dbm.PermRoleListPermsOwn) { cmd.AddCommand(newRoleListPermsCmd(user).cmd) }
	root.cmd = cmd
	return &root
}

type roleAddCmd struct {
	cmd *cobra.Command
}
func newRoleAddCmd() *roleAddCmd {
	root := roleAddCmd{}
	cmd := &cobra.Command{
		Use: "add ROLENAME",
		Short: "Add a role",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := dbm.OpenSQLite()
			if err != nil {
				return err
			}
			defer dbm.CloseDB(db, &err)

			if err := dbm.RoleAddSimple(db, args[0]); err != nil {
				return err
			}
			return nil
		},
	}
	root.cmd = cmd
	return &root
}

type roleListPermsCmd struct {
	cmd *cobra.Command
}
func newRoleListPermsCmd(user *dbm.UserWithRole) *roleListPermsCmd {
	root := roleListPermsCmd{}
	cmd := &cobra.Command{
		Use: "perms [ROLE]",
		Short: "List permissions for a role, if no ROLE present, use current user role",
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := dbm.OpenSQLite()
			if err != nil { return err }
			defer dbm.CloseDB(db, &err)

			var role dbm.Role
			if len(args) == 0 {
				role = user.Role
			} else {
				roleName := args[0]
				roleID, err := dbm.RoleGetID(db, roleName)
				if err != nil {
					return fmt.Errorf("could not get id for role with name %s: %w", roleName, err)
				}
				role, err = dbm.RoleGetWithPerms(db, roleID)
				if err != nil {
					return fmt.Errorf("getting role with id %d: %w", roleID, err)
				}
			}

			permSlice := make([]dbm.Permission, 0, len(role.Permissions))
			for perm, enabled := range role.Permissions {
				if !enabled { continue }
				permSlice = append(permSlice, perm)
			}
			slices.Sort(permSlice)
			for _, perm := range permSlice {
				fmt.Printf("%s\n", perm)
			}

			return nil
		},
	}
	root.cmd = cmd
	return &root
}
