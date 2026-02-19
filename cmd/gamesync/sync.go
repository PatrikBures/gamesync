package main

import (
	"fmt"
	"gamesync/internal/syncer"

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
		RunE: func(cmd *cobra.Command, args []string) error {
			gameID := args[0]

			if err := syncer.HandleSync(current, gameID, syncer.ModeAuto, false, true); err != nil {
				return fmt.Errorf("error syncing game: %v", err)
			}

			return nil
		},
	}

	root.cmd = cmd

	return &root
}
