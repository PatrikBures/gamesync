package main

import (
	"fmt"
	"gamesync/internal/config"
	"gamesync/internal/state"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use: "sync GAME_ID",
	Short: "Either pushes or pulls save for GAME_ID",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		game, _, err := config.GetGame(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting game: %v\n", err)
			os.Exit(40)
		}

		s, err := state.Get(game.SavePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting state of game with id \"%s\": %v\n", game.ID, err)
			os.Exit(41)
		}

		stateDir, err := config.GetStateDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting state dir: %v\n", err)
		}

		stateFile := filepath.Join(stateDir, game.ID+".json")
		
		if err := state.Write(s, stateFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing state for game with id \"%s\":%v\n", game.ID, err)
			os.Exit(42)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
