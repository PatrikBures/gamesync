package main

import (
	"os"

	"gamesync/internal/ui"
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
		Run: func(cmd *cobra.Command, args []string) {
			gameID := args[0]

			if err := syncer.HandleSync(current, gameID, syncer.ModePull, root.opts.force); err != nil {
				ui.Error("Error pulling: %v\n", err)
				os.Exit(20)
			}
		},
	}

	cmd.Flags().BoolVarP(&root.opts.force, "force", "f", false, "Overwrite local save with remote")

	root.cmd = cmd

	return &root
}



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
		Run: func(cmd *cobra.Command, args []string) {
			gameID := args[0]

			if err := syncer.HandleSync(current, gameID, syncer.ModePush, root.opts.force); err != nil {
				ui.Error("Error pushing: %v\n", err)
				os.Exit(20)
			}
		},
	}

	cmd.Flags().BoolVarP(&root.opts.force, "force", "f", false, "Overwrite remote save with local")

	root.cmd = cmd

	return &root
}



func init() {
	rootCmd.AddCommand(newPullCmd().cmd)
	rootCmd.AddCommand(newPushCmd().cmd)
}
