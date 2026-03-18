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
		Use: "add",
		Short: "Add stuff to server",
	}
	cmd.AddCommand(
		newUserAddCmd().cmd,
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
