package main

import (
	"fmt"
	"log"
	"path/filepath"

	"gamesync/internal/config"

	"github.com/spf13/cobra"
)

var current config.Current

var rootCmd = &cobra.Command{
	Use: "gamesync",
	Short: "Syncs save games to a server",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if current.ConfigPath == "" {
			configDir, err := config.GetConfigDir()
			if err != nil {
				return fmt.Errorf("error getting config dir: %w", err)
			}
			current.ConfigPath = filepath.Join(configDir, "config.yml")
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
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&current.ConfigPath, "config", "", "config file (default ~/.config/gamesync/config.yml)")
	rootCmd.PersistentFlags().BoolVarP(&current.Verbose, "verbose", "v", false, "more verbose output")
}
