package main

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"gamesync/internal/config"
	"gamesync/internal/syncer"
)

var pullCmd = &cobra.Command{
	Use: "pull GAME_ID",
	Short: "Pull the save if remote is newer",
	Example: "gamesync pull openttd",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gameID := args[0]

		if err := syncer.HandleSync(config.Current.Server, gameID, syncer.ModePull, false, verbose); err != nil {
			fmt.Fprintf(os.Stderr, "Error pulling: %v\n", err)
			os.Exit(20)
		}
	},
}

var pushCmd = &cobra.Command{
	Use: "push GAME_ID",
	Short: "Push the save if remote is older",
	Example: "gamesync push openttd",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gameID := args[0]

		if err := syncer.HandleSync(config.Current.Server, gameID, syncer.ModePush, false, verbose); err != nil {
			fmt.Fprintf(os.Stderr, "Error pushing: %v\n", err)
			os.Exit(20)
		}
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(pushCmd)
}
