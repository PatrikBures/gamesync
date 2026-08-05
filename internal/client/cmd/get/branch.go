package get

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
	cmd  *cobra.Command
	opts branchOpts
}

type branchOpts struct {
	repoName string
}

func newBranchCmd(conf *config.Config) *branchCmd {
	root := branchCmd{}

	cmd := &cobra.Command{
		Use:   "branch REPO",
		Short: "Get branches in repo",
		Args:  cobra.ExactArgs(1),
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
}

func runBranchCmd(c *api.Client, conf *config.Config, opts *branchOpts) error {
	branches, err := c.GetBranches(context.Background(), api.GetBranchesParams{
		UserID:   conf.Server.UserID,
		RepoName: opts.repoName,
	})
	if err != nil {
		return util.ErrHandler(err)
	}

	if len(branches) == 0 {
		fmt.Printf("no branches found in repo '%s'\n", opts.repoName)
		return nil
	}

	for _, b := range branches {
		fmt.Println(b)
	}

	return nil
}
