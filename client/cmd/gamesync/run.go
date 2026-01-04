package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"gamesync/internal/config"
)

var runCmd = &cobra.Command{
	Use: "run GAME_ID",
	Short: "Sync saves and run game",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gameID := args[0]
		game, err := config.GetGame(gameID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Syncing %s\n", game.ID)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
