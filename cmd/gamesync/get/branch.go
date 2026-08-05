package get

import (
	"context"
	"fmt"

	"go.pabu.dev/gamesync/internal/client"
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
}

func newBranchCmd(conf *config.Config) *branchCmd {
	root := branchCmd{}

	cmd := &cobra.Command{
		Use: "branch",
		Short: "SUMMARY_PLACEHOLDER",
		Args: cobra.ExactArgs(1),
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
	return nil
}

func runBranchCmd(c *api.Client, conf *config.Config, opts *branchOpts) error {
	branches, err := c.GetUserRepoBranches(context.Background(), api.GetUserRepoBranchesParams{
		UserID: conf.Server.UserID,
		RepoName: opts.repoName,
	})
	if err != nil {
		return fmt.Errorf("getting branches: %w", err)
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
