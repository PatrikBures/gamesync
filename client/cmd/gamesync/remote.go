package main

import (
	"fmt"
	"gamesync/internal/config"
	"gamesync/internal/syncer"

	"github.com/spf13/cobra"
)

var remoteCmd = &cobra.Command{
	Use: "remote <cmd>",
	Short: "Manage remote",
}

var remoteLsCmd = &cobra.Command{
	Use: "ls",
	Short: "List remote saves",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		remoteSaves, err := syncer.RunCmd(config.Current.Server, verbose, "ls")

		if err != nil {
			fmt.Println("Error listing remote:", err)
		}

		fmt.Print(remoteSaves)
	},
}

func init() {
	remoteCmd.AddCommand(remoteLsCmd)
	rootCmd.AddCommand(remoteCmd)
}
