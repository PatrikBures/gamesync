package main

import (
	"fmt"
	"gamesync/internal/config"
	"gamesync/internal/syncer"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var snapshotSkipUnchanged bool

var snapshotCmd = &cobra.Command{
	Use: "snapshot",
	Short: "manage snapshots made using restic on the remote",
}

var snapshotCreateCmd = &cobra.Command{
	Use: "create GAME_ID",
	Short: "Uses restic to create a snapshot of a game save",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gameID := args[0]

		if err := syncer.CreateSnapshot(config.Current.Server, verbose, gameID, snapshotSkipUnchanged); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating a snapshot: %v\n", err)
			os.Exit(4)
		}
	},
}

var snapshotLsCmd = &cobra.Command{
	Use: "ls [GAME_ID]",
	Short: "List snapshots of GAME_ID or all",
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gameID := ""
		if len(args) > 0 {
			gameID = args[0]
		}

		snapshots, err := syncer.ListSnapshots(config.Current.Server, verbose, gameID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting list of snapshots: %v\n", err)
			os.Exit(4)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)


		if _, err = fmt.Fprintf(w, "Name:\tHostname:\tTime:\n"); err != nil {
			fmt.Fprintf(os.Stderr, "Error adding table header to writer: %v\n", err)
			os.Exit(4)
		}

		for _, snapshot := range snapshots {
			if _, err = fmt.Fprintf(w, "%s\t%s\t%s\n", filepath.Base(snapshot.Paths[0]), snapshot.Hostname, snapshot.Time.Format(time.DateTime)); err != nil {
				fmt.Fprintf(os.Stderr, "Error adding row to writer: %v\n", err)
				os.Exit(4)
			}
		}

		if err := w.Flush(); err != nil {
			fmt.Fprintf(os.Stderr, "Error flushing: %v\n", err)
			os.Exit(4)
		}
	},
}


func init() {
	snapshotCmd.AddCommand(snapshotCreateCmd)
	snapshotCreateCmd.Flags().BoolVarP(&snapshotSkipUnchanged, "skip-unchanged", "S", false, "Skips creating snapshot if there were no changes from the last one")

	snapshotCmd.AddCommand(snapshotLsCmd)

	rootCmd.AddCommand(snapshotCmd)
}
