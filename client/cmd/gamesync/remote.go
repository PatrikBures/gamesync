package main

import (
	"fmt"
	"gamesync/internal/config"
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
		remoteSaves, err := syncer.RunCmd(config.Current.Server, verbose, "ls")

		if err != nil {
			fmt.Println("Error listing remote:", err)
		}

		fmt.Print(remoteSaves)
	},
}

var remoteRmCmd = &cobra.Command{
	Use: "rm GAME_ID",
	Short: "Remove a remote save for a game",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gameID := args[0]

		output, err := syncer.RemoveSaveGame(gameID, config.Current.Server, verbose)

		if err != nil {
			fmt.Println("Error removing save:", err)
			fmt.Println(output)
		}
	},
}

var remoteBackupCmd = &cobra.Command{
	Use: "backup GAME_ID",
	Short: "Uses restic to create a backup a game save",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gameID := args[0]

		if err := syncer.BackupGame(config.Current.Server, verbose, gameID); err != nil {
			fmt.Println("Error backing up save:", err)
		}
	},
}

func init() {
	remoteCmd.AddCommand(remoteLsCmd)
	remoteCmd.AddCommand(remoteRmCmd)
	remoteCmd.AddCommand(remoteBackupCmd)

	rootCmd.AddCommand(remoteCmd)
}
