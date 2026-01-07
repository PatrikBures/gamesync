package main

import (
	"fmt"
	"gamesync/internal/config"
	"gamesync/internal/syncer"

	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use: "snapshot",
	Short: "manage snapshots made using restic on the remote",
}

var snapshotCreateCmd = &cobra.Command{
	Use: "create GAME_ID",
	Short: "Uses restic to create a snapshot of a game save",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gameID := args[0]

		if err := syncer.BackupGame(config.Current.Server, verbose, gameID); err != nil {
			fmt.Println("Error creating a snapshot save:", err)
		}
	},
}


func init() {
	snapshotCmd.AddCommand(snapshotCreateCmd)

	rootCmd.AddCommand(snapshotCmd)
}
