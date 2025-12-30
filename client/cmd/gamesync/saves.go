package main

import (
	"fmt"
	"gamesync/internal/config"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var savePath string

var savesCmd = &cobra.Command{
	Use: "saves <cmd>",
	Short: "Manage saves",
}

var savesAddCmd = &cobra.Command{
	Use: "add <game_id>",
	Short: "add save dir of game to sync",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gameID := args[0]

		savePath, err := filepath.Abs(savePath)
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
	Short: "list all games",
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
	Use: "rm <game_id>...",
	Short: "remove a game ids from config",
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
	savesAddCmd.Flags().StringVarP(&savePath, "path", "d", "", "Directory path of the save (required)")

	if err := savesAddCmd.MarkFlagRequired("path"); err != nil {
		fmt.Println("Error marking flag as required")
		os.Exit(30)
	}
	if err := savesAddCmd.MarkFlagDirname("path"); err != nil {
		fmt.Println("Error marking flag as dirname")
	}

	savesCmd.AddCommand(savesLsCmd)
	savesCmd.AddCommand(savesRmCmd)

	rootCmd.AddCommand(savesCmd)
}
