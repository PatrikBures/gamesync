package create

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

type branchCmd struct {
	cmd  *cobra.Command
	opts branchOpts
}

type branchOpts struct {
	repoName   string
	branchName string
	snapshotID int64
}

func newBranchCmd(conf *config.Config) *branchCmd {
	root := branchCmd{}

	cmd := &cobra.Command{
		Use:   "branch REPO BRANCH SNAPSHOT",
		Short: "Create new branch in repo",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(conf)
			if err != nil {
				return err
			}
			populateBranchOpts(&root.opts, args)

			return runBranchCmd(c, conf, &root.opts)
		},
	}

	root.cmd = cmd
	return &root
}

func populateBranchOpts(opts *branchOpts, args []string) error {
	opts.repoName = args[0]
	opts.branchName = args[1]
	
	snapshotID, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return fmt.Errorf("parsing snapshot: %w", err)
	}

	opts.snapshotID = snapshotID

	return nil
}

func runBranchCmd(c *api.Client, conf *config.Config, opts *branchOpts) error {
	if err := c.PutBranch(context.Background(), &api.SnapshotID{
		SnapshotID: opts.snapshotID,
	}, api.PutBranchParams{
		UserID:     conf.Server.UserID,
		RepoName:   opts.repoName,
		BranchName: opts.branchName,
		
		
	}); err != nil {
		return util.ErrHandler(err)
	}
	fmt.Printf("branch '%s' created in repo '%s'\n", opts.branchName, opts.repoName)
	return nil
}
