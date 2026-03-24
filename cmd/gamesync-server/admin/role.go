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
func newRoleCmd(user *dbm.UserWithRole) *roleCmd {
	root := roleCmd{}
	cmd := &cobra.Command{
		Use: "role",
		Short: "Manage roles",
	}
	if user.Role.HasPermission(dbm.PermRoleAdd)             { cmd.AddCommand(newRoleAddCmd().cmd) }
	if user.Role.HasPermission(dbm.PermRoleListPermsOwn)    { cmd.AddCommand(newRoleListPermsCmd(user).cmd) }
	if user.Role.HasPermission(dbm.PermRoleDelete)          { cmd.AddCommand(newRoleDeleteCmd().cmd) }
	if user.Role.HasPermission(dbm.PermRoleList)            { cmd.AddCommand(newRoleListCmd().cmd) }
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

type roleListCmd struct {
	cmd *cobra.Command
}
func newRoleListCmd() *roleListCmd {
	root := roleListCmd{}
	cmd := &cobra.Command{
		Use: "ls",
		Aliases: []string{"list"},
		Short: "List all roles",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := dbm.OpenSQLite()
			if err != nil { return err }
			defer dbm.CloseDB(db, &err)

			roles, err := dbm.RoleGetAll(db)
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
func newRoleDeleteCmd() *roleDeleteCmd {
	root := roleDeleteCmd{}
	cmd := &cobra.Command{
		Use: "delete ROLE",
		Short: "Delete a role with no users",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := dbm.OpenSQLite()
			if err != nil { return err }
			defer dbm.CloseDB(db, &err)

			rolename := args[0]
			roleID, err := dbm.RoleGetID(db, rolename)
			if err != nil {
				return fmt.Errorf("getting role id: %w", err)
			}
			if root.opts.removeUsers {
				amountDeleted, err := dbm.UserDeleteAllInRole(db, roleID)
				if err != nil {
					return fmt.Errorf("removing all users with role %s and id %d: %w", rolename, roleID, err)
				}
				ui.Info("Deleted %d users in role %s\n", amountDeleted, rolename)
			}
			if err := dbm.RoleDeleteWithID(db, roleID); err != nil {
				return fmt.Errorf("removing role %s with ID %d: %w", rolename, roleID, err)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&root.opts.removeUsers, "remove-users", "", false, "Removes role along with all users in it")
	root.cmd = cmd
	return &root
}
