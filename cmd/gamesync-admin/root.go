package main

import (
	"fmt"
	"gamesync/internal/dbm"
	"os"

	"github.com/spf13/cobra"
)

func loadUserRole() (*dbm.Role, error) {
	userName := os.Getenv("GAMESYNC_USER")

	if userName == "" {
		role := dbm.Role{
			ID: -1,
			Name: "root",
			Permissions: dbm.RoleAllPerms(true),
		}
		return &role, nil
	}
	fmt.Println("running as user:", userName)

	db, err := dbm.OpenSQLite()
	if err != nil {
		return nil, err
	}
	defer dbm.CloseDB(db, &err)

	user, err := dbm.UserGet(db, userName)
	if err != nil {
		return nil, fmt.Errorf("getting user %s: %w", userName, err)
	}

	role, err := dbm.RoleGetWithPerms(db, user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("getting role for user: %s: %w", user.Name, err)
	}
	return &role, nil
}

type rootCmd struct {
	cmd *cobra.Command
}

func newRootCmd(role *dbm.Role) *rootCmd {
	root := rootCmd{}
	cmd := &cobra.Command{
		Use: "gamesync-admin",
		Short: "Manage the gamesync server",
	}
	cmd.DisableAutoGenTag = true
	cmd.SilenceUsage = true

	if role.HasPermission(dbm.PermUserAdd)           { cmd.AddCommand(newUserCmd(role).cmd) }
	if role.HasPermission(dbm.PermRoleChangePerms)   { cmd.AddCommand(newRoleCmd(role).cmd) }

	// uid is 0 if ran by root
	// cmds are only available when ran directly from the container
	if os.Getuid() == 0 {
		cmd.AddCommand(newInitCmd().cmd)
	}
	root.cmd = cmd
	return &root
}

func Execute() {
	role, err := loadUserRole()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	if err := newRootCmd(role).cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
