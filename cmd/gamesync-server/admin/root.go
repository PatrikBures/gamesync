package cmdAdmin

import (
	"database/sql"
	"fmt"
	"gamesync/internal/dbm"
	"os"

	"github.com/spf13/cobra"
)

type userDB struct {
	user *dbm.UserWithRole
	db *sql.DB
}

func loadUserRole(db *sql.DB) (*dbm.UserWithRole, error) {
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

	userWithRole, err := dbm.UserGetWithRole(db, username)
	if err != nil {
		return nil, fmt.Errorf("getting user with their role: %w", err)
	}
	return userWithRole, nil
}

type rootCmd struct {
	cmd *cobra.Command
}

func newRootCmd(udb userDB) *rootCmd {
	root := rootCmd{}
	cmd := &cobra.Command{
		Use: "gamesync-admin",
		Short: "Manage the gamesync server",
	}
	cmd.DisableAutoGenTag = true
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	if udb.user.Role.HasPermission(dbm.PermUserAdd)           { cmd.AddCommand(newUserCmd(udb).cmd) }
	if udb.user.Role.HasPermission(dbm.PermRoleChangePerms)   { cmd.AddCommand(newRoleCmd(udb).cmd) }
	if udb.user.Role.HasPermission(dbm.PermKeyAddSelf)        { cmd.AddCommand(newKeyCmd(udb).cmd) }

	// uid is 0 if ran by root
	// cmds are only available when ran directly from the container
	if os.Getuid() == 0 {
		cmd.AddCommand(newInitCmd(udb).cmd)
	}
	root.cmd = cmd
	return &root
}

func Execute() error {
	db, err := dbm.OpenSQLite()
	if err != nil { return err }
	defer dbm.CloseDB(db, &err)

	user, err := loadUserRole(db)
	if err != nil {
		return err
	}
	udb := userDB{user: user, db: db}
	if err := newRootCmd(udb).cmd.Execute(); err != nil {
		return err
	}
	return nil
}
