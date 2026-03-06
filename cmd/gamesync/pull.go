package main

import (
	"fmt"

	"gamesync/internal/syncer"

	"github.com/spf13/cobra"
)


type pullCmd struct {
	cmd *cobra.Command
	opts pullOpts
}

type pullOpts struct {
	force bool
}

func newPullCmd() *pullCmd {
	root := pullCmd{}

	cmd := &cobra.Command{
		Use: "pull GAME_ID",
		Short: "Pull the save if remote is newer",
		Example: "gamesync pull openttd",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			gameID := args[0]

			if err := syncer.HandleSync(current, gameID, syncer.ModePull, root.opts.force, true, false); err != nil {
				return fmt.Errorf("error pulling: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&root.opts.force, "force", "f", false, "Overwrite local save with remote")

	root.cmd = cmd

	return &root
}
