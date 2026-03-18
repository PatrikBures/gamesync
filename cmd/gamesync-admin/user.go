package main

import (
	"fmt"
	"gamesync/internal/dbm"

	"github.com/spf13/cobra"
)

type userCmd struct {
	cmd *cobra.Command
}
func newUserCmd() *userCmd {
	root := userCmd{}
	cmd := &cobra.Command{
		Use: "user",
		Short: "Add stuff to server",
	}
	cmd.AddCommand(
		newUserAddCmd().cmd,
		newUserChangeRoleCmd().cmd,
	)
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
				RoleID: 1,
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
