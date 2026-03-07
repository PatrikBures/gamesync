package main

import (
	"fmt"

	"gamesync/internal/syncer"
	"gamesync/internal/ui"

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
		RunE: func(cmd *cobra.Command, args []string) error {
			output, err := syncer.RunCmd(current.Config.Server, "list-saves")

			if err != nil {
				return fmt.Errorf("error listing remote: %v\n%s", err, output)
			}

			ui.Info("%s", output)

			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			gameID := args[0]

			output, err := syncer.RemoveSaveGame(current, gameID)

			if err != nil {
				return fmt.Errorf("error removing save: %v\n%s", err, output)
			}

			ui.Info("%s\n", gameID)

			return nil
		},
	}

	root.cmd = cmd

	return &root
}
