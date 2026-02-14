package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"gamesync/internal/config"

	"github.com/spf13/cobra"
)

type saveCmd struct {
	cmd *cobra.Command
}

func newSaveCmd() *saveCmd {
	root := saveCmd{}

	cmd := &cobra.Command{
		Use: "save",
		Short: "Manage game saves in local config",
	}

	cmd.AddCommand(newSaveAddCmd().cmd)
	cmd.AddCommand(newSaveLsCmd().cmd)
	cmd.AddCommand(newSaveRmCmd().cmd)

	root.cmd = cmd

	return &root
}




type saveAddCmd struct {
	cmd *cobra.Command
	opts saveAddOpts
}

type saveAddOpts struct {
	update bool
}

func newSaveAddCmd() *saveAddCmd {
	root := saveAddCmd{}

	cmd := &cobra.Command{
		Use: "add GAME_ID SAVE_PATH",
		Short: "Add save dir of a game to sync",
		Example: "gamesync save add openttd -d ~/.local/share/openttd/save",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			gameID := args[0]

			savePath, err := filepath.Abs(args[1])
			if err != nil {
				return fmt.Errorf("Error getting absolute path: %v\n", err)
			}

			
			if err := config.AddSave(&current, gameID, savePath, root.opts.update); err != nil {
				return fmt.Errorf("Error adding save: %v\n", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&root.opts.update, "update", "u", false, "Skips checking if gameID already exists")

	root.cmd = cmd

	return &root
}



type saveLsCmd struct {
	cmd *cobra.Command
	opts saveLsOpts
}

type saveLsOpts struct {
	quiet bool
}

func newSaveLsCmd() *saveLsCmd {
	root := saveLsCmd{}
	
	cmd := &cobra.Command{
		Use: "ls",
		Short: "List all games in local config",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)

			for _, game := range current.Config.Games {
				var err error

				if root.opts.quiet {
					_, err = fmt.Fprintf(w, "%s\n", game.ID)
				} else {
					_, err = fmt.Fprintf(w, "%s\t%s\n", game.ID, game.SavePath)
				}

				if err != nil {
					return fmt.Errorf("Error printing: %v\n", err)
				}
			}

			if err := w.Flush(); err != nil {
				return fmt.Errorf("Error flushing: %v\n", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&root.opts.quiet, "quiet", "q", false, "Only prints gameIDs")

	root.cmd = cmd

	return &root
}


type saveRmCmd struct {
	cmd *cobra.Command
}

func newSaveRmCmd() *saveRmCmd {
	root := saveRmCmd{}

	cmd := &cobra.Command{
		Use: "rm GAME_ID...",
		Short: "Remove a game ids from config",
		Long: `Does not remove the save dir`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := config.RemoveGames(&current, args)

			if err != nil {
				return fmt.Errorf("Error removing games: %v\n", err)
			}

			return nil
		},
	}

	root.cmd = cmd

	return &root
}

func init() {
	rootCmd.AddCommand(newSaveCmd().cmd)
}
