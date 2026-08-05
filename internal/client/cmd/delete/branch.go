package delete

import (
	"context"
	"fmt"

	"go.pabu.dev/gamesync/internal/client"
	util "go.pabu.dev/gamesync/internal/client/cmd/_util"
	"go.pabu.dev/gamesync/internal/client/config"
	api "go.pabu.dev/gamesync/internal/ogen"

	"github.com/spf13/cobra"
)

type branchCmd struct {
	cmd *cobra.Command
	opts branchOpts
}

type branchOpts struct {
	repoName string
	branchName string
}

func newBranchCmd(conf *config.Config) *branchCmd {
	root := branchCmd{}

	cmd := &cobra.Command{
		Use: "branch REPO BRANCH",
		Short: "Deletes a branch from a repo",
		Args: cobra.ExactArgs(2),
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

func populateBranchOpts(opts *branchOpts, args []string) {
	opts.repoName = args[0]
	opts.branchName = args[1]
}

func runBranchCmd(c *api.Client, conf *config.Config, opts *branchOpts) error {
	if err := c.DeleteBranch(context.Background(), api.DeleteBranchParams{
		UserID: conf.Server.UserID,
		RepoName: opts.repoName,
		BranchName: opts.branchName,
	}); err != nil {
		return util.ErrHandler(err)
	}
	fmt.Printf("Deleted branch '%s' from repo '%s'\n", opts.branchName, opts.repoName)
	return nil
}
