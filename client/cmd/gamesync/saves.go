package main

import (
	"fmt"
	"gamesync/internal/config"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	savePath	string
	saveUpdate	bool
)

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

		configPath := ""
		if saveUpdate {
			configPath = configFile
		}
		
		if err := config.AddSave(gameID, savePath, configPath); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	},
}

func init() {
	savesCmd.AddCommand(savesAddCmd)
	savesAddCmd.Flags().StringVarP(&savePath, "path", "d", "", "Directory path of the save (required)")
	savesAddCmd.Flags().BoolVarP(&saveUpdate, "update", "u", false, "Updates global config")

	if err := savesAddCmd.MarkFlagRequired("path"); err != nil {
		fmt.Println("Error marking flag as required")
		os.Exit(30)
	}
	if err := savesAddCmd.MarkFlagDirname("path"); err != nil {
		fmt.Println("Error marking flag as dirname")
	}

	rootCmd.AddCommand(savesCmd)
}
