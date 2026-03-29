package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gamesync/internal/config"
	"gamesync/internal/ui"
	"gamesync/internal/vars"

	"github.com/spf13/cobra"
)

var current config.Current

type rootCmd struct {
	cmd *cobra.Command
	opts rootOpts
}

type rootOpts struct {
	verbose bool
	debug bool
	ignoreApiVersion bool
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
			if root.opts.debug {
				ui.SetLevel(ui.LevelDebug)
			} else if root.opts.verbose {
				ui.SetLevel(ui.LevelVerbose)
			}
			
			if root.opts.ignoreApiVersion {
				vars.ApiVersion = "0"
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&current.ConfigPath, "config", "", "config file (default ~/.config/gamesync/config.yml)")
	cmd.PersistentFlags().BoolVarP(&root.opts.verbose, "verbose", "v", false, "more verbose output")
	cmd.PersistentFlags().BoolVarP(&root.opts.debug, "debug", "", false, "debug ouput")
	cmd.PersistentFlags().BoolVarP(&root.opts.ignoreApiVersion, "ignore-api-version", "", false, "Do not check if remote has same api version. Use with caution!")

	cmd.AddCommand(
		newGenDocCmd().cmd,
		newPullCmd().cmd,
		newPushCmd().cmd,
		newSyncCmd().cmd,
		newSaveCmd().cmd,
		newRemoteCmd().cmd,
		newWrapCmd().cmd,
		newVersionCmd().cmd,
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
