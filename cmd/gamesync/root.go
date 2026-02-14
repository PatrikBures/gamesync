package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gamesync/internal/config"
	"gamesync/internal/ui"

	"github.com/spf13/cobra"
)

var current config.Current

type rootCmd struct {
	cmd *cobra.Command
}

func newRootCmd() *rootCmd {
	root := rootCmd{}

	cmd := &cobra.Command{
		Use: "gamesync",
		Short: "Syncs save games to a server",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if current.ConfigPath == "" {
				configEnv := os.Getenv("GAMESYNC_CONFIG")

				if configEnv == "" {
					configDir, err := config.GetConfigDir()
					if err != nil {
						return fmt.Errorf("error getting config dir: %w", err)
					}
					current.ConfigPath = filepath.Join(configDir, "config.yml")
				} else {
					current.ConfigPath = configEnv
				}
			} else {
				var err error
				current.ConfigPath, err = filepath.Abs(current.ConfigPath)
				if err != nil {
					return err
				}
			}
			if err := config.Load(&current); err != nil {
				return fmt.Errorf("loading config: %w", err)
			}
			if current.Verbose {
				ui.SetLevel(ui.LevelVerbose)
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&current.ConfigPath, "config", "", "config file (default ~/.config/gamesync/config.yml)")
	cmd.PersistentFlags().BoolVarP(&current.Verbose, "verbose", "v", false, "more verbose output")

	cmd.AddCommand(
		newGenDocCmd().cmd,
		newPullCmd().cmd,
		newPushCmd().cmd,
		newSyncCmd().cmd,
		newSaveCmd().cmd,
		newRemoteCmd().cmd,
		newSnapshotCmd().cmd,
		newWrapCmd().cmd,
	)

	cmd.DisableAutoGenTag = true

	root.cmd = cmd

	return &root
}

func Execute() {
	if err := newRootCmd().cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
