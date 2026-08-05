package get

import (
	"context"
	"fmt"

	"go.pabu.dev/gamesync/internal/client"
	"go.pabu.dev/gamesync/internal/client/config"
	api "go.pabu.dev/gamesync/internal/ogen"

	"github.com/spf13/cobra"
)

type repoCmd struct {
	cmd *cobra.Command
	opts repoOpts
}

type repoOpts struct {}

func newRepoCmd(conf *config.Config) *repoCmd {
	root := repoCmd{}

	cmd := &cobra.Command{
		Use: "repo",
		Short: "Get all repos",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(conf)
			if err != nil {
				return err
			}

			return runRepoCmd(c, conf, &root.opts)
		},
	}

	root.cmd = cmd
	return &root
}

func populateRepoOpts(opts *repoOpts, args []string) error {
	return nil
}

func runRepoCmd(c *api.Client, conf *config.Config, opts *repoOpts) error {
	repos, err := c.GetUserRepos(context.Background(), api.GetUserReposParams{ UserID: conf.Server.UserID })
	if err != nil {
		return fmt.Errorf("getting repos: %w", err)
	}

	if len(repos) == 0 {
		fmt.Println("no repos found")
		return nil
	}

	for _, r := range repos {
		fmt.Println(r)
	}
	
	return nil
}
