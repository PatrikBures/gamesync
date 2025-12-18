package main

import (
	"fmt"
	"gamesync/internal/config"
	"gamesync/internal/syncer"
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

var savesSyncCmd = &cobra.Command{
	Use: "sync <game_id> <push/pull>",
	Short: "Sync game save",
}

var savesSyncPullCmd = &cobra.Command{
	Use: "pull <game_id>",
	Short: "Pull the save if remote is newer",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := syncOrPull(args[0], true); err != nil {
			fmt.Println(err)
			os.Exit(3)
		}
	},
}

var savesSyncPushCmd = &cobra.Command{
	Use: "push <game_id>",
	Short: "Push the save if remote is older",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := syncOrPull(args[0], true); err != nil {
			fmt.Println(err)
			os.Exit(3)
		}
	},

}

func syncOrPull(gameID string, pull bool) error {
	game, err := config.GetGame(gameID)
	if err != nil {
		return err
	}

	if err := syncer.SyncGame(*game, config.Current.Server, pull); err != nil {
		return err
	}

	return nil
}

func init() {
	savesCmd.AddCommand(savesAddCmd)
	savesAddCmd.Flags().StringVarP(&savePath, "path", "d", "", "Directory path of the save (required)")
	savesAddCmd.Flags().BoolVarP(&saveUpdate, "update", "u", false, "Updates global config")

	savesAddCmd.MarkFlagRequired("path")
	savesAddCmd.MarkFlagDirname("path")

	savesCmd.AddCommand(savesSyncCmd)
	savesSyncCmd.AddCommand(savesSyncPullCmd)
	savesSyncCmd.AddCommand(savesSyncPushCmd)

	rootCmd.AddCommand(savesCmd)
}
