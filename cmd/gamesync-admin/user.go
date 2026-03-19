package main

import (
	"fmt"
	"gamesync/internal/dbm"
	"gamesync/internal/ui"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type userCmd struct {
	cmd *cobra.Command
}
func newUserCmd(role *dbm.Role) *userCmd {
	root := userCmd{}
	cmd := &cobra.Command{
		Use: "user",
		Short: "Add stuff to server",
	}
	if role.HasPermission(dbm.PermUserAdd)        { cmd.AddCommand(newUserAddCmd().cmd) }
	if role.HasPermission(dbm.PermUserChangeRole) { cmd.AddCommand(newUserChangeRoleCmd().cmd) }
	if role.HasPermission(dbm.PermUserList)       { cmd.AddCommand(newUserLsCmd().cmd) }
	root.cmd = cmd
	return &root
}


type userAddCmd struct {
	cmd *cobra.Command
}
func newUserAddCmd() *userAddCmd {
	root := userAddCmd{}
	cmd := &cobra.Command{
		Use: "add USERNAME ROLENAME",
		Short: "Add a new user",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := dbm.UserAddSimple(dbm.User{
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
		Use: "change-role USERNAME ROLENAME",
		Short: "Changes role of a user",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			userName := args[0]
			roleName := args[1]
			if err := dbm.UserChangeRoleSimple(userName, roleName); err != nil {
				return fmt.Errorf("chaning role for %s to %s: %w", userName, roleName, err)
			}
			return nil
		},
	}
	root.cmd = cmd
	return &root
}

type userLsCmd struct {
	cmd *cobra.Command
}
func newUserLsCmd() *userLsCmd {
	root := userLsCmd{}
	cmd := &cobra.Command{
		Use: "ls",
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
