package cmdAdmin

import (
	"fmt"
	"gamesync/internal/dbm"
	"gamesync/internal/ui"

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
