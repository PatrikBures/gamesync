package create

import (
	"context"
	"fmt"
	util "gamesync/cmd/gamesync/_util"
	"gamesync/internal/client"
	"gamesync/internal/client/config"
	api "gamesync/internal/ogen"

	"github.com/spf13/cobra"
)

type repoCmd struct {
	cmd *cobra.Command
	opts repoOpts
}
type repoOpts struct {
	repoName string
}

func newRepoCmd(conf *config.Config) *repoCmd {
	root := repoCmd{}

	cmd := &cobra.Command{
		Use: "repo",
		Short: "Create a new repo",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := client.New(conf)
			if err != nil {
				return err
			}
			if err := populateRepoOpts(&root.opts, args); err != nil {
				return err
			}
			if err := runRepoCmd(client, root.opts, *conf); err != nil {
				return err
			}
			return nil
		},
	}

	root.cmd = cmd
	return &root
}

func populateRepoOpts(opts *repoOpts, args []string) error {
	opts.repoName = args[0]
	if opts.repoName == "" {
		return fmt.Errorf("repo name can not be empty")
	}
	return nil
}

func runRepoCmd(client *api.Client, opts repoOpts, conf config.Config) error {
	err := client.PutUserRepo(context.Background(), api.PutUserRepoParams{
		UserID: conf.Server.UserID,
		RepoName: opts.repoName,
	})
	if err := util.ErrHandler(err); err != nil {
		return err
	}
	fmt.Printf("repo named '%s' created\n", opts.repoName)
	return nil
}
