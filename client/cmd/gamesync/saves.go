package main

import (
	"fmt"
	"gamesync/internal/config"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var savesCmd = &cobra.Command{
	Use: "saves",
	Short: "Manage saves",
}

var savesAddUpdate bool

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

		
		if err := config.AddSave(gameID, savePath, configFile, savesAddUpdate); err != nil {
			fmt.Println("Error adding save:", err)
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

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)

		for _, game := range *games {
			if _, err := fmt.Fprintf(w, "%s\t%s\n", game.ID, game.SavePath); err != nil {
				fmt.Println("Error printing:", err)
			}
		}

		if err := w.Flush(); err != nil {
			fmt.Println("Error flushing:", err)
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
	savesAddCmd.Flags().BoolVarP(&savesAddUpdate, "update", "u", false, "Skips checking if gameID already exists")

	savesCmd.AddCommand(savesLsCmd)
	savesCmd.AddCommand(savesRmCmd)

	rootCmd.AddCommand(savesCmd)
}
