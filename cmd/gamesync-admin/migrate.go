package main

import (
	"fmt"
	"gamesync/internal/dbm"

	"github.com/spf13/cobra"
)

type migrateCmd struct {
	cmd *cobra.Command
}

func newMigrateCmd() *migrateCmd {
	root := migrateCmd{}
	cmd := &cobra.Command{
		Use: "migrate",
		Short: "migrates db to newer schema, initializes if it doesn't exist",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := dbm.Migrate(); err != nil {
				return fmt.Errorf("migrating: %v", err)
			}

			return nil
		},
	}
	root.cmd = cmd
	return &root
}
