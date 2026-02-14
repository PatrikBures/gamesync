package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	"gamesync/internal/syncer"
	"gamesync/internal/ui"

	"github.com/spf13/cobra"
)

var snapshotSkipUnchanged bool

type snapshotCmd struct {
	cmd *cobra.Command
}

func newSnapshotCmd() *snapshotCmd {
	root := snapshotCmd{}

	cmd := &cobra.Command{
		Use: "snapshot",
		Short: "manage snapshots made using restic on the remote",
	}

	cmd.AddCommand(newSnapshotLsCmd().cmd)
	cmd.AddCommand(newSnapshotCreateCmd().cmd)

	root.cmd = cmd

	return &root
}



type snapshotCreateCmd struct {
	cmd *cobra.Command
	opts snapshotCreateOpts
}

type snapshotCreateOpts struct {
	skipUnchanged bool
}

func newSnapshotCreateCmd() *snapshotCreateCmd {
	root := snapshotCreateCmd{}

	cmd := &cobra.Command{
		Use: "create GAME_ID",
		Short: "Uses restic to create a snapshot of a game save",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			gameID := args[0]

			if err := syncer.CreateSnapshot(current, gameID, snapshotSkipUnchanged); err != nil {
				ui.Error("Error creating a snapshot: %v\n", err)
				os.Exit(4)
			}
		},
	}

	cmd.Flags().BoolVarP(&root.opts.skipUnchanged, "skip-unchanged", "S", false, "Skips creating snapshot if there were no changes from the last one")

	root.cmd = cmd

	return &root
}



type snapshotLsCmd struct {
	cmd *cobra.Command
}

func newSnapshotLsCmd() *snapshotLsCmd {
	root := snapshotLsCmd{}

	cmd := &cobra.Command{
		Use: "ls [GAME_ID]",
		Short: "List snapshots of GAME_ID or all",
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			gameID := ""
			if len(args) > 0 {
				gameID = args[0]
			}

			snapshots, err := syncer.ListSnapshots(current, gameID)
			if err != nil {
				ui.Error("Error getting list of snapshots: %v\n", err)
				os.Exit(4)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)


			if _, err = fmt.Fprintf(w, "Name:\tHostname:\tTime:\n"); err != nil {
				ui.Error("Error adding table header to writer: %v\n", err)
				os.Exit(4)
			}

			for _, snapshot := range snapshots {
				if _, err = fmt.Fprintf(w, "%s\t%s\t%s\n", filepath.Base(snapshot.Paths[0]), snapshot.Hostname, snapshot.Time.Format(time.DateTime)); err != nil {
					ui.Error("Error adding row to writer: %v\n", err)
					os.Exit(4)
				}
			}

			if err := w.Flush(); err != nil {
				ui.Error("Error flushing: %v\n", err)
				os.Exit(4)
			}
		},
	}

	root.cmd = cmd

	return &root
}



func init() {
	rootCmd.AddCommand(newSnapshotCmd().cmd)
}
