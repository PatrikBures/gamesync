package cmdAdmin

import (
	"fmt"
	"gamesync/internal/dbm"
	"gamesync/internal/ui"
	"gamesync/internal/vars"
	"os"

	"github.com/spf13/cobra"
)

type initCmd struct {
	cmd * cobra.Command
}
func newInitCmd(udb userDB) *initCmd {
	root := initCmd{}
	cmd := &cobra.Command{
		Use: "init",
		Short: "Commands used to initialize the server",
	}
	cmd.AddCommand(
		newInitDirsCmd().cmd,
		newInitMigrateCmd().cmd,
		newInitRolesCmd(udb).cmd,
		newInitPermsCmd().cmd,
	)
	root.cmd = cmd
	return &root
}


type initMigrateCmd struct {
	cmd *cobra.Command
}
func newInitMigrateCmd() *initMigrateCmd {
	root := initMigrateCmd{}
	cmd := &cobra.Command{
		Use: "migrate",
		Short: "Migrates db to newer schema, initializes if it doesn't exist",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := dbm.Migrate(); err != nil {
				return fmt.Errorf("migrating: %v", err)
			}
			if err := os.Chown(vars.RemoteSQLiteDb, vars.RemoteUID, -1); err != nil {
				return fmt.Errorf("changing owner of db: %w", err)
			}
			return nil
		},
	}
	root.cmd = cmd
	return &root
}


type initDirsCmd struct {
	cmd *cobra.Command
}
func newInitDirsCmd() *initDirsCmd {
	root := initDirsCmd{}
	cmd := &cobra.Command{
		Use: "dirs",
		Short: "Creates the necessary dirs for the server, moves old paths",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, path := range vars.RemoteSaveDirOld {
				f, err := os.Stat(path)
				if err != nil || !f.IsDir(){
					continue
				}
				if err := os.Rename(path, vars.RemoteSaveDir); err != nil {
					return fmt.Errorf("moving %s to %s: %w", path, vars.RemoteSaveDir, err)
				}
				ui.Info("moved save dir from %s to %s\n", path, vars.RemoteSaveDir)
				break
			}

			dirs := []string{
				vars.RemoteSaveDir,
				vars.RemoteBackupDir,
				vars.RemoteDbDir,
			}
			for _, dir := range dirs {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return err
				}
				if err := os.Chmod(dir, 0775); err != nil {
					return err
				}
				if err := os.Chown(dir, 1000, -1); err != nil {
					return err
				}
				ui.Info("Ensured dir exists: %s\n", dir)
			}
			return nil
		},
	}
	root.cmd = cmd
	return &root
}

type initPermsCmd struct {
	cmd *cobra.Command
}
func newInitPermsCmd() *initPermsCmd {
	root := initPermsCmd{}
	cmd := &cobra.Command{
		Use: "perms",
		Short: "Ensures all permissions exist",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := dbm.PermsSet(); err != nil {
				return fmt.Errorf("setting perms: %w", err)
			}
			return nil
		},
	}
	root.cmd = cmd
	return &root
}

type initRolesCmd struct {
	cmd *cobra.Command
}
func newInitRolesCmd(udb userDB) *initRolesCmd {
	root := initRolesCmd{}
	cmd := &cobra.Command{
		Use: "roles",
		Short: "Ensures user and admin roles exist in the db",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			roles := []dbm.Role{
				{
					ID: 0,
					Name: "none",
					Permissions: dbm.RoleAllPerms(false),
				},
				{
					ID: 1,
					Name: "admin", 
					Permissions: dbm.RoleAllPerms(true),
				},
			}
			for _, role := range roles {
				_ = dbm.RoleDeleteWithID(udb.db, role.ID)
				if err := dbm.RoleAddWithPerms(udb.db, role); err != nil {
					return fmt.Errorf("adding role with perms: %w", err)
				}
				ui.Info("added role: %s\n", role.Name)
			}
			
			return nil
		},
	}
	root.cmd = cmd
	return &root
}
