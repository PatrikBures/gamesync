package delete

import (
	"context"
	"fmt"
	"strconv"

	"go.pabu.dev/gamesync/internal/client"
	util "go.pabu.dev/gamesync/internal/client/cmd/_util"
	"go.pabu.dev/gamesync/internal/client/config"
	api "go.pabu.dev/gamesync/internal/ogen"

	"github.com/spf13/cobra"
)

type snapshotCmd struct {
	cmd *cobra.Command
	opts snapshotOpts
}

type snapshotOpts struct {
	repoName string
	snapshotID int64
}

func newSnapshotCmd(conf *config.Config) *snapshotCmd {
	root := snapshotCmd{}

	cmd := &cobra.Command{
		Use: "snapshot REPO SNAPSHOT",
		Short: "Delete snapshot",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := populateSnapshotOpts(&root.opts, args); err != nil {
				return err
			}

			c, err := client.New(conf)
			if err != nil {
				return err
			}

			return runSnapshotCmd(c, conf, &root.opts)
		},
	}

	root.cmd = cmd
	return &root
}

func populateSnapshotOpts(opts *snapshotOpts, args []string) error {
	opts.repoName = args[0]

	snapshotID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("parsing snapshot id: %w", err)
	}
	opts.snapshotID = snapshotID

	return nil
}

func runSnapshotCmd(c *api.Client, conf *config.Config, opts *snapshotOpts) error {
	if err := c.DeleteSnapshot(context.Background(), api.DeleteSnapshotParams{
		UserID: conf.Server.UserID,
		RepoName: opts.repoName,
		SnapshotID: opts.snapshotID,
	}); err != nil {
		return util.ErrHandler(err)
	}
	fmt.Printf("Deleted snapshot %d from repo '%s'\n", opts.snapshotID, opts.repoName)
	return nil
}
