package main

import (
	"os"

	"github.com/spf13/cobra"
)

type rootCmd struct {
	cmd *cobra.Command
}

func newRootCmd() *rootCmd {
	root := rootCmd{}

	cmd := &cobra.Command{
		Use: "gamesync",
		Short: "Syncs save games to a server",
	}

	cmd.AddCommand(
		newGenDocCmd().cmd,
	)

	cmd.DisableAutoGenTag = true

	root.cmd = cmd

	return &root
}

func Execute() {
	if err := newRootCmd().cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
