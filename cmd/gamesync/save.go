package main

import (
	"fmt"
	"gamesync/internal/config"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var saveCmd = &cobra.Command{
	Use: "save",
	Short: "Manage game saves in local config",
}

var saveAddUpdate bool
var saveLsQuiet bool

var saveAddCmd = &cobra.Command{
	Use: "add GAME_ID SAVE_PATH",
	Short: "Add save dir of a game to sync",
	Example: "gamesync save add openttd -d ~/.local/share/openttd/save",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		gameID := args[0]

		savePath, err := filepath.Abs(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting absolute path: %v\n", err)
			os.Exit(2)
		}

		
		if err := config.AddSave(&current, gameID, savePath, saveAddUpdate); err != nil {
			fmt.Fprintf(os.Stderr, "Error adding save: %v\n", err)
			os.Exit(2)
		}
	},
}

var saveLsCmd = &cobra.Command{
	Use: "ls",
	Short: "List all games in local config",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)

		for _, game := range current.Config.Games {
			var err error

			if saveLsQuiet {
				_, err = fmt.Fprintf(w, "%s\n", game.ID)
			} else {
				_, err = fmt.Fprintf(w, "%s\t%s\n", game.ID, game.SavePath)
			}

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error printing: %v\n", err)
			}
		}

		if err := w.Flush(); err != nil {
			fmt.Fprintf(os.Stderr, "Error flushing: %v\n", err)
		}
	},
}

var saveRmCmd = &cobra.Command{
	Use: "rm GAME_ID...",
	Short: "Remove a game ids from config",
	Long: `Does not remove the save dir`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := config.RemoveGames(&current, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error removing games: %v\n", err)
			os.Exit(23)
		}
	},
}

func init() {
	saveCmd.AddCommand(saveAddCmd)
	saveAddCmd.Flags().BoolVarP(&saveAddUpdate, "update", "u", false, "Skips checking if gameID already exists")

	saveCmd.AddCommand(saveLsCmd)
	saveLsCmd.Flags().BoolVarP(&saveLsQuiet, "quiet", "q", false, "Only prints gameIDs")
	saveCmd.AddCommand(saveRmCmd)

	rootCmd.AddCommand(saveCmd)
}
