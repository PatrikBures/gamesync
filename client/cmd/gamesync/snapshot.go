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

		if err := syncer.CreateSnapshot(config.Current.Server, verbose, gameID); err != nil {
			fmt.Println("Error creating a snapshot save:", err)
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
			fmt.Printf("Error listing snapshots: %v\n", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)


		var printErr error
		_, printErr = fmt.Fprintf(w, "Name:\tHostname:\tTime:\n")

		for _, snapshot := range snapshots {
			_, printErr = fmt.Fprintf(w, "%s\t%s\t%s\n", filepath.Base(snapshot.Paths[0]), snapshot.Hostname, snapshot.Time.Format(time.DateTime))
		}

		if printErr != nil {
			fmt.Println("Error printing:", printErr)
		}

		if err := w.Flush(); err != nil {
			fmt.Println("Error flushing:", err)
		}
	},
}


func init() {
	snapshotCmd.AddCommand(snapshotCreateCmd)
	snapshotCmd.AddCommand(snapshotLsCmd)

	rootCmd.AddCommand(snapshotCmd)
}
