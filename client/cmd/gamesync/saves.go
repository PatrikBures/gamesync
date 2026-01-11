package main

import (
	"fmt"
	"gamesync/internal/config"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var savesCmd = &cobra.Command{
	Use: "saves",
	Short: "Manage saves",
}

var savesAddCmd = &cobra.Command{
	Use: "add GAME_ID SAVE_PATH",
	Short: "Add save dir of a game to sync",
	Example: "gamesync saves add openttd -d ~/.local/share/openttd/save",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		gameID := args[0]

		savePath, err := filepath.Abs(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

		
		if err := config.AddSave(gameID, savePath, configFile); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	},
}

var savesLsCmd = &cobra.Command{
	Use: "ls",
	Short: "List all games in local config",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		games := &config.Current.Games

		if games == nil {
			fmt.Println("no games found in config")
		} else {
			for _, game := range *games {
				fmt.Println(game.ID)
			}
		}
	},
}

var savesRmCmd = &cobra.Command{
	Use: "rm GAME_ID...",
	Short: "Remove a game ids from config",
	Long: `Does not remove the save dir`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := config.RemoveGames(args, configFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(23)
		}
	},
}

func init() {
	savesCmd.AddCommand(savesAddCmd)

	savesCmd.AddCommand(savesLsCmd)
	savesCmd.AddCommand(savesRmCmd)

	rootCmd.AddCommand(savesCmd)
}
