package main

import (
	"fmt"

	"gamesync/internal/syncer"

	"github.com/spf13/cobra"
)


type pushCmd struct {
	cmd *cobra.Command
	opts pushOpts
}

type pushOpts struct {
	force bool
}

func newPushCmd() *pushCmd {
	root := pushCmd{}

	cmd := &cobra.Command{
		Use: "push GAME_ID",
		Short: "Push the save if remote is older",
		Example: "gamesync push openttd",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			gameID := args[0]

			if _, err := syncer.HandleSync(current, gameID, syncer.ModePush, root.opts.force, true, false); err != nil {
				return fmt.Errorf("error pushing: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&root.opts.force, "force", "f", false, "Overwrite remote save with local")

	root.cmd = cmd

	return &root
}
