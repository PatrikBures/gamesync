package main

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"gamesync/internal/config"
	"gamesync/internal/syncer"
)

var pullCmd = &cobra.Command{
	Use: "pull <game_id>",
	Short: "Pull the save if remote is newer",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := pushOrPull(args[0], true); err != nil {
			fmt.Printf("Error pulling, %v\n", err)
			os.Exit(20)
		}
	},
}

var pushCmd = &cobra.Command{
	Use: "push <game_id>",
	Short: "Push the save if remote is older",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := pushOrPull(args[0], false); err != nil {
			fmt.Printf("Error pushing, %v\n", err)
			os.Exit(21)
		}
	},
}

func pushOrPull(gameID string, pull bool) error {
	game, err := config.GetGame(gameID)
	if err != nil {
		return err
	}

	if err := syncer.SyncGame(*game, config.Current.Server, pull); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(pushCmd)
}
