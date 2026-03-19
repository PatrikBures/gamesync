package main

import (
	"gamesync/internal/dbm"

	"github.com/spf13/cobra"
)

type roleCmd struct {
	cmd *cobra.Command
}
func newRoleCmd(role *dbm.Role) *roleCmd {
	root := roleCmd{}
	cmd := &cobra.Command{
		Use: "role",
		Short: "Manage roles",
	}
	if role.HasPermission(dbm.PermRoleAdd) { cmd.AddCommand(newRoleAddCmd().cmd) }
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
