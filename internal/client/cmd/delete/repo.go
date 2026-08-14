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

type repoCmd struct {
	cmd  *cobra.Command
	opts repoOpts
}

type repoOpts struct {
	repoName string
}

func newRepoCmd(conf *config.Config) *repoCmd {
	root := repoCmd{}

	cmd := &cobra.Command{
		Use:   "repo REPO",
		Short: "Deletes a repo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(conf)
			if err != nil {
				return err
			}
			populateRepoOpts(&root.opts, args)

			return runRepoCmd(c, conf, &root.opts)
		},
	}

	root.cmd = cmd
	return &root
}

func populateRepoOpts(opts *repoOpts, args []string) {
	opts.repoName = args[0]
}

func runRepoCmd(c *api.Client, conf *config.Config, opts *repoOpts) error {
	if err := c.DeleteRepo(context.Background(), api.DeleteRepoParams{
		UserID:   conf.Server.UserID,
		RepoName: opts.repoName,
	}); err != nil {
		return util.ErrHandler(err)
	}
	fmt.Printf("Deleted repo '%s'\n", opts.repoName)
	return nil
}
