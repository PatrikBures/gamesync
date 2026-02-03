package main

import (
	"os"

	"gamesync/internal/ui"
	"gamesync/internal/syncer"

	"github.com/spf13/cobra"
)

var remoteCmd = &cobra.Command{
	Use: "remote",
	Short: "Manage remote saves",
}

var remoteLsCmd = &cobra.Command{
	Use: "ls",
	Short: "List remote saves",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		remoteSaves, err := syncer.RunCmd(current, "list-saves")

		if err != nil {
			ui.Error("Error listing remote: %v\n", err)
			os.Exit(3)
		}

		ui.Info("%s\n", remoteSaves)
	},
}

var remoteRmCmd = &cobra.Command{
	Use: "rm GAME_ID",
	Short: "Remove a remote save for a game",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gameID := args[0]

		output, err := syncer.RemoveSaveGame(current, gameID)

		if err != nil {
			ui.Error("Error removing save: %v\n%s", err, output)
			os.Exit(3)
		}

		ui.Info("%s\n", gameID)
	},
}

func init() {
	remoteCmd.AddCommand(remoteLsCmd)
	remoteCmd.AddCommand(remoteRmCmd)

	rootCmd.AddCommand(remoteCmd)
}
