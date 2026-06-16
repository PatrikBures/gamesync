package main

import (
	"fmt"
	"gamesync/internal/client/config"
	"os"

	"github.com/spf13/cobra"
)

type rootCmd struct {
	cmd    *cobra.Command
	config config.Config
}

func newRootCmd() *rootCmd {
	root := rootCmd{}

	cmd := &cobra.Command{
		Use:   "gamesync",
		Short: "Syncs save games to a server",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := config.LoadConfig(&root.config); err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&root.config.Token, "token", "", "Token used to authenticate with server")
	cmd.PersistentFlags().StringVar(&root.config.Server, "server", "", "Server url")

	cmd.AddCommand(
		newGenDocCmd().cmd,
		newInitCmd(&root.config).cmd,
		newGetCmd(&root.config).cmd,
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
