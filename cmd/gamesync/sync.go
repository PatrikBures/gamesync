package main

import (
	"gamesync/internal/syncer"
	"gamesync/internal/ui"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use: "sync GAME_ID",
	Short: "Either pushes or pulls save for GAME_ID",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gameID := args[0]
		if err := syncer.HandleSync(current, gameID, syncer.ModeAuto, false); err != nil {
			ui.Error("Error syncing game: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
