package main

import (
	"fmt"
	"log"

	"gamesync/internal/config"

	"github.com/spf13/cobra"
)

var configFile string

var rootCmd = &cobra.Command{
	Use: "gamesync",
	Short: "A CLI to sync save games to a server",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Load(configFile); err != nil {
			return fmt.Errorf("Failed loading config: %w", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GameSync!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default ~/.config/gamesync/config.yml)")
}
