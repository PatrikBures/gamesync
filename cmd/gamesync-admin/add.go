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
			db, err := dbm.OpenSQLite()
			if err != nil {
				return err
			}
			defer dbm.CloseDB(db, &err)

			user := dbm.User{
				UserName: args[0],
				UserTypeId: dbm.UserTypeUser,
			}
			
			if err := dbm.AddUser(db, user); err != nil {
				return fmt.Errorf("creating user: %v", err)
			}

			return nil
		},
	}
	root.cmd = cmd
	return &root
}
