package main

import (
	"fmt"
	"gamesync/internal/dbm"
	"os"

	"github.com/spf13/cobra"
)

func loadUserRole() (*dbm.UserWithRole, error) {
	username := os.Getenv("GAMESYNC_USER")

	if username == "" {
		u := dbm.UserWithRole{
			Name: "root",
			ID: -1,
			Role: dbm.Role{
				ID: -1,
				Name: "root",
				Permissions: dbm.RoleAllPerms(true),
			},
		}
		return &u, nil
	}

	db, err := dbm.OpenSQLite()
	if err != nil {
		return nil, err
	}
	defer dbm.CloseDB(db, &err)

	userWithRole, err := dbm.UserGetWithRole(db, username)
	if err != nil {
		return nil, fmt.Errorf("getting user with their role: %w", err)
	}
	return userWithRole, nil
}

type rootCmd struct {
	cmd *cobra.Command
}

func newRootCmd(user *dbm.UserWithRole) *rootCmd {
	root := rootCmd{}
	cmd := &cobra.Command{
		Use: "gamesync-admin",
		Short: "Manage the gamesync server",
	}
	cmd.DisableAutoGenTag = true
	cmd.SilenceUsage = true

	if user.Role.HasPermission(dbm.PermUserAdd)           { cmd.AddCommand(newUserCmd(user).cmd) }
	if user.Role.HasPermission(dbm.PermRoleChangePerms)   { cmd.AddCommand(newRoleCmd(user).cmd) }
	if user.Role.HasPermission(dbm.PermKeyAddSelf)        { cmd.AddCommand(newKeyCmd(user).cmd) }

	// uid is 0 if ran by root
	// cmds are only available when ran directly from the container
	if os.Getuid() == 0 {
		cmd.AddCommand(newInitCmd().cmd)
	}
	root.cmd = cmd
	return &root
}

func Execute() {
	user, err := loadUserRole()
	if err != nil {
		os.Exit(2)
	}
	if err := newRootCmd(user).cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
