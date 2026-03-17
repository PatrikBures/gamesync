package main

import (
	"fmt"
	"gamesync/internal/dbm"

	"github.com/spf13/cobra"
)

type addCmd struct {
	cmd *cobra.Command
}
func newAddCmd() *addCmd {
	root := addCmd{}
	cmd := &cobra.Command{
		Use: "add",
		Short: "Add stuff to server",
	}
	cmd.AddCommand(
		newAddUserCmd().cmd,
	)
	root.cmd = cmd
	return &root
}


type addUserCmd struct {
	cmd *cobra.Command
}
func newAddUserCmd() *addUserCmd {
	root := addUserCmd{}
	cmd := &cobra.Command{
		Use: "user USERNAME",
		Short: "Add a new user",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := dbm.AddUserSimple(dbm.User{
				UserName: args[0], 
				UserRoleId: dbm.UserRoleUser,
			}); err != nil {
				return fmt.Errorf("adding new user: %v", err)
			}
			return nil
		},
	}
	root.cmd = cmd
	return &root
}

type addAdminCmd struct {
	cmd *cobra.Command
}
func newAddAdminCmd() *addAdminCmd {
	root := addAdminCmd{}
	cmd := &cobra.Command{
		Use: "admin USERNAME",
		Short: "Add a new admin user which can not syns. Used only for admin commands",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := dbm.AddUserSimple(dbm.User{
				UserName: args[0], 
				UserRoleId: dbm.UserRoleAdmin,
			}); err != nil {
				return fmt.Errorf("adding new admin: %v", err)
			}
			return nil
		},
	}
	root.cmd = cmd
	return &root
}
