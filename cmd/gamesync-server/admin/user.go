package cmdAdmin

import (
	"fmt"
	"gamesync/internal/dbm"
	"gamesync/internal/ui"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type userCmd struct {
	cmd *cobra.Command
}
func newUserCmd(user *dbm.UserWithRole) *userCmd {
	root := userCmd{}
	cmd := &cobra.Command{
		Use: "user",
		Short: "Add stuff to server",
	}
	if user.Role.HasPermission(dbm.PermUserAdd)        { cmd.AddCommand(newUserAddCmd().cmd) }
	if user.Role.HasPermission(dbm.PermUserChangeRole) { cmd.AddCommand(newUserChangeRoleCmd().cmd) }
	if user.Role.HasPermission(dbm.PermUserList)       { cmd.AddCommand(newUserListCmd().cmd) }
	if user.Role.HasPermission(dbm.PermUserDelete)     { cmd.AddCommand(newUserDeleteCmd().cmd) }
	root.cmd = cmd
	return &root
}


type userAddCmd struct {
	cmd *cobra.Command
}
func newUserAddCmd() *userAddCmd {
	root := userAddCmd{}
	cmd := &cobra.Command{
		Use: "add USERNAME",
		Short: "Add a new user",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := dbm.OpenSQLite()
			if err != nil {
				return err
			}
			defer dbm.CloseDB(db, &err)

			if err := dbm.UserAdd(db, dbm.User{
				Name: args[0], 
				RoleID: 0,
			}); err != nil {
				return fmt.Errorf("adding new user: %v", err)
			}
			return nil
		},
	}
	root.cmd = cmd
	return &root
}


type userChangeRoleCmd struct {
	cmd *cobra.Command
}
func newUserChangeRoleCmd() *userChangeRoleCmd {
	root := userChangeRoleCmd{}
	cmd := &cobra.Command{
		Use: "role USERNAME ROLENAME",
		Short: "Change role of a user",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]
			roleName := args[1]
			if err := dbm.UserChangeRoleSimple(username, roleName); err != nil {
				return fmt.Errorf("changing role for %s to %s: %w", username, roleName, err)
			}
			return nil
		},
	}
	root.cmd = cmd
	return &root
}

type userListCmd struct {
	cmd *cobra.Command
}
func newUserListCmd() *userListCmd {
	root := userListCmd{}
	cmd := &cobra.Command{
		Use: "ls",
		Aliases: []string{"list"},
		Short: "Lists all users",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := dbm.OpenSQLite()
			if err != nil {
				return err
			}
			defer dbm.CloseDB(db, &err)

			users, err := dbm.UserGetAll(db)
			if err != nil {
				return fmt.Errorf("getting list of all users: %w", err)
			}
			if len(users) == 0 {
				ui.Info("Found no users in database\n")
				return nil
			}

			w := tabwriter.NewWriter(ui.OutWriter, 0, 0, 2, ' ', tabwriter.TabIndent)
			if _, err := fmt.Fprintln(w, "Username:\tRole:"); err != nil {
				return fmt.Errorf("printing header: %w", err)
			}
			for _, user := range users {
				if _, err := fmt.Fprintf(w, "%s\t%s\n", user.Name, user.RoleName); err != nil {
					return fmt.Errorf("adding user to writer: %w", err)
				}
			}
			if err := w.Flush(); err != nil {
				return fmt.Errorf("flushing: %w", err)
			}

			ui.Info("%d users found.\n", len(users))
			return nil
		},
	}
	root.cmd  = cmd
	return &root
}


type userDeleteCmd struct {
	cmd *cobra.Command
}
func newUserDeleteCmd() *userDeleteCmd {
	root := userDeleteCmd{}
	cmd := &cobra.Command{
		Use: "delete USER",
		Short: "Delete a user and ALL of it's files",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := dbm.OpenSQLite()
			if err != nil { return err }
			defer dbm.CloseDB(db, &err)

			username := args[0]
			user, err := dbm.UserGet(db, username)
			if err != nil {
				return fmt.Errorf("getting user with name %s: %w", username, err)
			}
			if err := dbm.UserDelete(db, user.ID); err != nil {
				return fmt.Errorf("deleting user %s with id %d: %w", user.Name, user.ID, err)
			}
			return nil
		},
	}
	root.cmd = cmd
	return &root
}
