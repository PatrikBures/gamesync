package main

import (
	"gamesync/internal/ui"
	"gamesync/internal/vars"

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
			ui.Info("gamesync version: %s\n", vars.Version)
			ui.Info("gamesync api version: %s\n", vars.ApiVersion)

			return nil
		},
	}
	root.cmd = cmd
	return &root
}
