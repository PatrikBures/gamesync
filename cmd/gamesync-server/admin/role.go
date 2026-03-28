package cmdAdmin

import (
	"fmt"
	"gamesync/internal/dbm"
	"gamesync/internal/ui"
	"slices"

	"github.com/spf13/cobra"
)

type roleCmd struct {
	cmd *cobra.Command
}
func newRoleCmd(udb userDB) *roleCmd {
	root := roleCmd{}
	cmd := &cobra.Command{
		Use: "role",
		Short: "Manage roles",
	}
	if udb.user.Role.HasPermission(dbm.PermRoleAdd)             { cmd.AddCommand(newRoleAddCmd(udb).cmd) }
	if udb.user.Role.HasPermission(dbm.PermRoleDelete)          { cmd.AddCommand(newRoleDeleteCmd(udb).cmd) }
	if udb.user.Role.HasPermission(dbm.PermRoleList)            { cmd.AddCommand(newRoleListCmd(udb).cmd) }
	if udb.user.Role.HasPermission(dbm.PermRolePermListOwn)     { cmd.AddCommand(newRolePermCmd(udb).cmd) }
	root.cmd = cmd
	return &root
}

type roleAddCmd struct {
	cmd *cobra.Command
}
func newRoleAddCmd(udb userDB) *roleAddCmd {
	root := roleAddCmd{}
	cmd := &cobra.Command{
		Use: "add ROLENAME",
		Short: "Add a role",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := dbm.RoleAddSimple(udb.db, args[0]); err != nil {
				return err
			}
			return nil
		},
	}
	root.cmd = cmd
	return &root
}

type rolePermCmd struct {
	cmd *cobra.Command
}
func newRolePermCmd(udb userDB) *rolePermCmd {
	root := rolePermCmd{}
	cmd := &cobra.Command{
		Use: "perm",
		Short: "Manage permissions for roles",
	}
	if udb.user.Role.HasPermission(dbm.PermRolePermListOwn) { cmd.AddCommand(newRolePermListCmd(udb).cmd) }
	if udb.user.Role.HasPermission(dbm.PermRolePermMod)     { cmd.AddCommand(newRolePermAddCmd(udb).cmd) }
	root.cmd = cmd
	return &root
}

type rolePermListCmd struct {
	cmd *cobra.Command
}
func newRolePermListCmd(udb userDB) *rolePermListCmd {
	root := rolePermListCmd{}
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

type rolePermAddCmd struct {
	cmd *cobra.Command
}
func newRolePermAddCmd(udb userDB) *rolePermAddCmd {
	root := rolePermAddCmd{}
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


type roleListCmd struct {
	cmd *cobra.Command
}
func newRoleListCmd(udb userDB) *roleListCmd {
	root := roleListCmd{}
	cmd := &cobra.Command{
		Use: "ls",
		Aliases: []string{"list"},
		Short: "List all roles",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			roles, err := dbm.RoleGetAll(udb.db)
			if err != nil {
				return fmt.Errorf("getting all roles: %w", err)
			}
			if len(roles) == 0 {
				ui.Info("No roles found in database\n")
				return nil
			}
			for _, role := range roles {
				ui.Info("%s\n", role.Name)
			}
			ui.Info("%d roles found\n", len(roles))

			return nil
		},
	}
	root.cmd = cmd
	return &root
}

type roleDeleteCmd struct {
	cmd *cobra.Command
	opts roleDeleteOpts
}
type roleDeleteOpts struct {
	removeUsers bool
}
func newRoleDeleteCmd(udb userDB) *roleDeleteCmd {
	root := roleDeleteCmd{}
	cmd := &cobra.Command{
		Use: "delete ROLE",
		Short: "Delete a role with no users",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			rolename := args[0]
			roleID, err := dbm.RoleGetID(udb.db, rolename)
			if err != nil {
				return fmt.Errorf("getting role id: %w", err)
			}
			if root.opts.removeUsers {
				amountDeleted, err := dbm.UserDeleteAllInRole(udb.db, roleID)
				if err != nil {
					return fmt.Errorf("removing all users with role %s and id %d: %w", rolename, roleID, err)
				}
				ui.Info("Deleted %d users in role %s\n", amountDeleted, rolename)
			}
			if err := dbm.RoleDeleteWithID(udb.db, roleID); err != nil {
				return fmt.Errorf("removing role %s with ID %d: %w", rolename, roleID, err)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&root.opts.removeUsers, "remove-users", "", false, "Removes role along with all users in it")
	root.cmd = cmd
	return &root
}
