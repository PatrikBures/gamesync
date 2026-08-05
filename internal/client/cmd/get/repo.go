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

type repoCmd struct {
	cmd *cobra.Command
}

func newRepoCmd(conf *config.Config) *repoCmd {
	root := repoCmd{}

	cmd := &cobra.Command{
		Use:   "repo",
		Short: "Get all repos",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := client.New(conf)
			if err != nil {
				return err
			}

			return runRepoCmd(c, conf)
		},
	}

	root.cmd = cmd
	return &root
}

func runRepoCmd(c *api.Client, conf *config.Config) error {
	repos, err := c.GetRepos(context.Background(), api.GetReposParams{UserID: conf.Server.UserID})
	if err != nil {
		return util.ErrHandler(err)
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
