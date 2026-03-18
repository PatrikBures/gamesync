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
		Use: "gamesync-admin",
		Short: "Manage the gamesync server",
	}
	cmd.DisableAutoGenTag = true
	cmd.SilenceUsage = true
	cmd.AddCommand(
		newUserCmd().cmd,
	)

	// uid is 0 if ran by root
	// cmds are only available when ran directly from the container
	if os.Getuid() == 0 {
		cmd.AddCommand(newInitCmd().cmd)
	}
	root.cmd = cmd
	return &root
}

func Execute() {
	if err := newRootCmd().cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
