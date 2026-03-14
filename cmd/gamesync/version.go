package main

import (
	"gamesync/internal/config"
	"gamesync/internal/ui"

	"github.com/spf13/cobra"
)

type versionCmd struct {
	cmd *cobra.Command
}

func newVersionCmd() *versionCmd {
	root := versionCmd{}

	cmd := &cobra.Command{
		Use: "version",
		Short: "Print current version number",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ui.Info("gamesync version: %s\n", version)
			ui.Info("gamesync api version: %s\n", config.ApiVersion)

			return nil
		},
	}
	root.cmd = cmd
	return &root
}
