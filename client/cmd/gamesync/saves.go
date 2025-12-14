package main

import (
	"fmt"
	"gamesync/internal/config"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	poolID		string
	savePath	string
	saveUpdate	bool
	saveMove	bool
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
		err := config.AddSave(gameID, poolID, savePath, configPath, saveMove)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	},
}

func init() {
	savesCmd.AddCommand(savesAddCmd)
	savesAddCmd.Flags().StringVarP(&poolID, "pool", "p", "", "ID of pool (required)")
	savesAddCmd.Flags().StringVarP(&savePath, "path", "d", "", "Directory path of the save (required)")
	savesAddCmd.Flags().BoolVarP(&saveUpdate, "update", "u", false, "Updates global config")
	savesAddCmd.Flags().BoolVarP(&saveMove, "move", "m", false, "Moves save to pool and symlinks back")

	savesAddCmd.MarkFlagRequired("pool")
	savesAddCmd.MarkFlagRequired("path")
	savesAddCmd.MarkFlagDirname("path")

	rootCmd.AddCommand(savesCmd)
}
