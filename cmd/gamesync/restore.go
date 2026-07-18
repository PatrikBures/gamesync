package main

import (
	"fmt"
	"gamesync/internal/client"
	"gamesync/internal/client/config"
	api "gamesync/internal/ogen"
	"gamesync/internal/snapshoter"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)


type restoreCmd struct {
	cmd *cobra.Command
	opts restoreOpts
}
type restoreOpts struct {
	branch string
	repo string
	dir string
	snapshotID int64
}

func newRestoreCmd(conf *config.Config) *restoreCmd {
	root := restoreCmd{}

	cmd := &cobra.Command{
		Use:   "restore DIR SNAPSHOT_ID",
		Short: "Restore a sync to a specific snapshot",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := populateRestoreOpts(&root.opts, args); err != nil { return err }

			c, err := client.Client(conf)
			if err != nil { return err }

			if err := runRestoreCmd(c, root.opts, conf); err != nil { return err }

			return nil
		},
	}

	cmd.Flags().StringVarP(&root.opts.repo,   "repo", "r", "", "")
	cmd.Flags().StringVarP(&root.opts.branch, "branch", "b", "", "")

	if err := markFlagsRequired(cmd, []string{"repo", "branch"}); err != nil {
		panic(err)
	}

	root.cmd = cmd
	return &root
}

func populateRestoreOpts(opts *restoreOpts, args []string) error {
	opts.dir = args[0]

	d, err := os.Stat(opts.dir)
	if err != nil { return err }
	if !d.IsDir() {
		return fmt.Errorf("provided dir is not a dir: %s", opts.dir)
	}

	opts.snapshotID, err = strconv.ParseInt(args[1], 10, 64)
	if err != nil { return err }

	return nil
}

func runRestoreCmd(c *api.Client, opts restoreOpts, conf *config.Config) (err error) {
	syncer, err := snapshoter.NewSyncer(conf, c, opts.repo, opts.branch, opts.dir, opts.snapshotID)
	if err != nil {
		return fmt.Errorf("initializing syncer: %w", err)
	}

	if err := syncer.Pull(); err != nil {
		return fmt.Errorf("pulling: %w", err)
	}

	return nil
}
