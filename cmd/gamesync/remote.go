package main

import (
	"os"

	"gamesync/internal/ui"
	"gamesync/internal/syncer"

	"github.com/spf13/cobra"
)

type remoteCmd struct {
	cmd *cobra.Command
}

func newRemoteCmd() *remoteCmd {
	root := remoteCmd{}

	cmd := &cobra.Command{
		Use: "remote",
		Short: "Manage remote saves",
	}

	cmd.AddCommand(newRemoteLsCmd().cmd)
	cmd.AddCommand(newRemoteRmCmd().cmd)

	root.cmd = cmd

	return &root
}



type remoteLsCmd struct {
	cmd *cobra.Command
}

func newRemoteLsCmd() *remoteLsCmd {
	root := remoteLsCmd{}

	cmd := &cobra.Command{
		Use: "ls",
		Short: "List remote saves",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			remoteSaves, err := syncer.RunCmd(current, "list-saves")

			if err != nil {
				ui.Error("Error listing remote: %v\n", err)
				os.Exit(3)
			}

			ui.Info("%s\n", remoteSaves)
		},
	}

	root.cmd = cmd

	return &root
}



type remoteRmCmd struct {
	cmd *cobra.Command
}

func newRemoteRmCmd() *remoteRmCmd {
	root := remoteRmCmd{}

	cmd := &cobra.Command{
		Use: "rm GAME_ID",
		Short: "Remove a remote save for a game",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			gameID := args[0]

			output, err := syncer.RemoveSaveGame(current, gameID)

			if err != nil {
				ui.Error("Error removing save: %v\n%s", err, output)
				os.Exit(3)
			}

			ui.Info("%s\n", gameID)
		},
	}

	root.cmd = cmd

	return &root
}



func init() {
	rootCmd.AddCommand(newRemoteCmd().cmd)
}
