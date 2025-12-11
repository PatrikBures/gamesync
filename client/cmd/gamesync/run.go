package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use: "run [game_id]",
	Short: "Sync saves and run game",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// check if game with that id exists
		game_ID := args[0]
		fmt.Printf("Syncing game %s...\n", game_ID)
		fmt.Printf("Starting game %s...\n", game_ID)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
