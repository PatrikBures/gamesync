package main

import (
	"fmt"
	"gamesync/internal/config"
	"gamesync/internal/syncer"
	"os"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use: "sync GAME_ID",
	Short: "Either pushes or pulls save for GAME_ID",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := syncer.HandleSync(config.Current.Server, args[0], syncer.ModeAuto, false,  verbose); err != nil {
			fmt.Fprintf(os.Stderr, "Error syncing game: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
