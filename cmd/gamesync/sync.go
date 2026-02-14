package main

import (
	"gamesync/internal/syncer"
	"gamesync/internal/ui"

	"github.com/spf13/cobra"
)

type syncCmd struct {
	cmd *cobra.Command
}

func newSyncCmd() *syncCmd {
	root := syncCmd{}

	cmd := &cobra.Command{
		Use: "sync GAME_ID",
		Short: "Either pushes or pulls save for GAME_ID",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			gameID := args[0]
			if err := syncer.HandleSync(current, gameID, syncer.ModeAuto, false); err != nil {
				ui.Error("urror syncing game: %v", err)
			}
		},
	}

	root.cmd = cmd

	return &root
}

func init() {
	rootCmd.AddCommand(newSyncCmd().cmd)
}
