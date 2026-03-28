package cmdAdmin

import (
	"fmt"
	"gamesync/internal/dbm"
	"gamesync/internal/ui"
	"slices"

	"github.com/spf13/cobra"
)

type permCmd struct {
	cmd *cobra.Command
}
func newPermCmd(udb userDB) *permCmd {
	root := permCmd{}
	cmd := &cobra.Command{
		Use: "perm",
		Short: "Manage permissions",
	}
	if udb.user.Role.HasPermission(dbm.PermRolePermListOwn) { cmd.AddCommand(newPermListCmd(udb).cmd) }
	if udb.user.Role.HasPermission(dbm.PermRolePermMod) {
		cmd.AddCommand(newPermAddCmd(udb).cmd)
		cmd.AddCommand(newPermRemoveCmd(udb).cmd)
	}

	root.cmd = cmd
	return &root
}

type permListCmd struct {
	cmd *cobra.Command
}
func newPermListCmd(udb userDB) *permListCmd {
	root := permListCmd{}
	cmd := &cobra.Command{
		Use: "ls [ROLE]",
		Short: "List permissions for a role, if no ROLE present, use current user role",
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var role dbm.Role
			if len(args) == 0 {
				role = udb.user.Role
			} else {
				roleName := args[0]
				roleID, err := dbm.RoleGetID(udb.db, roleName)
				if err != nil {
					return fmt.Errorf("could not get id for role with name %s: %w", roleName, err)
				}
				role, err = dbm.RoleGetWithPerms(udb.db, roleID)
				if err != nil {
					return fmt.Errorf("getting role with id %d: %w", roleID, err)
				}
			}

			if len(role.Permissions) == 0 {
				ui.Info("Role '%s' has no permissions\n", role.Name)
			}

			permSlice := make([]dbm.Permission, 0, len(role.Permissions))
			for perm, enabled := range role.Permissions {
				if !enabled { continue }
				permSlice = append(permSlice, perm)
			}
			slices.Sort(permSlice)
			for _, perm := range permSlice {
				ui.Info("%s\n", perm)
			}

			return nil
		},
	}
	root.cmd = cmd
	return &root
}

type permAddCmd struct {
	cmd *cobra.Command
}
func newPermAddCmd(udb userDB) *permAddCmd {
	root := permAddCmd{}
	cmd := &cobra.Command{
		Use: "add ROLE PERMS...",
		Short: "Add permissions to ROLE",
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			rolename := args[0]
			perms := args[1:]

			roleID, err := dbm.RoleGetID(udb.db, rolename)
			if err != nil {
				return err
			}
			if err := dbm.RoleAddPerms(udb.db, roleID, perms); err != nil {
				return err
			}
			return nil
		},
	}
	root.cmd = cmd
	return &root
}

type permRemoveCmd struct {
	cmd *cobra.Command
}
func newPermRemoveCmd(udb userDB) *permRemoveCmd {
	root := permRemoveCmd{}
	cmd := &cobra.Command{
		Use: "rm ROLE PERMS...",
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			rolename := args[0]
			perms := args[1:]

			roleID, err := dbm.RoleGetID(udb.db, rolename)
			if err != nil {
				return err
			}
			if err := dbm.RoleRemovePerms(udb.db, roleID, perms); err != nil {
				return fmt.Errorf("removing perms from role: %w", err)
			}
			return nil
		},
	}
	root.cmd = cmd
	return &root
}
