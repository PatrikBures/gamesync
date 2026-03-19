package main

import (
	"fmt"
	"gamesync/internal/dbm"

	"github.com/spf13/cobra"
)

type keyCmd struct {
	cmd *cobra.Command
}
func newKeyCmd(user *dbm.UserWithRole) *keyCmd {
	root := keyCmd{}
	cmd := &cobra.Command{
		Use: "key",
		Short: "Manage ssh keys",
	}
	if user.Role.HasPermission(dbm.PermKeyAdd) || user.Role.HasPermission(dbm.PermKeyAddSelf) { cmd.AddCommand(newKeyAddCmd(user).cmd) }
	root.cmd = cmd
	return &root
}

type keyAddCmd struct {
	cmd *cobra.Command
	opts keyAddOpts
}
type keyAddOpts struct {
	userName string
}
func newKeyAddCmd(user *dbm.UserWithRole) *keyAddCmd {
	root := keyAddCmd{}
	cmd := &cobra.Command{
		Use: "add PUBLIC_KEY",
		Short: "Add public ssh key to user (default self)",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := dbm.OpenSQLite()
			if err != nil { return err }
			defer dbm.CloseDB(db, &err)

			u := dbm.User{}
			if root.opts.userName == "" {
				if user.ID >= 0 {
					return fmt.Errorf("not running as a logged in user")
				}
				u.Name = user.Name
				u.ID = user.ID
			} else {
				s, err := dbm.UserGet(db, root.opts.userName)
				if err != nil {
					return fmt.Errorf("getting user %s: %w", root.opts.userName, err)
				}
				u = *s
			}

			pubKey := args[0]

			if err := dbm.KeyAdd(db, pubKey, u); err != nil {
				return fmt.Errorf("adding pub key: %w", err)
			}
			return nil
		},
	}
	if user.Role.HasPermission(dbm.PermKeyAdd) {
		cmd.Flags().StringVarP(&root.opts.userName, "user", "u", "", "Add key to specific user, otherwise add to yourself")
	}
	root.cmd = cmd
	return &root
}
