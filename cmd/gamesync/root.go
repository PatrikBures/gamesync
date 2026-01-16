package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gamesync/internal/config"

	"github.com/spf13/cobra"
)

var configFile string
var verbose bool

var rootCmd = &cobra.Command{
	Use: "gamesync",
	Short: "Syncs save games to a server",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if configFile == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("could not find home directory: %w", err)
			}
			configFile = filepath.Join(home, ".config", "gamesync", "config.yml")
		} else {
			var err error
			configFile, err = filepath.Abs(configFile)
			if err != nil {
				return err
			}
		}
		if err := config.Load(configFile); err != nil {
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
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default ~/.config/gamesync/config.yml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "more verbose output")
}
