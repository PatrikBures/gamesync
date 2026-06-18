package main

import (
	"context"
	"fmt"
	"gamesync/internal/client"
	"gamesync/internal/client/config"
	api "gamesync/internal/ogen"

	"github.com/spf13/cobra"
)


type createCmd struct {
	cmd *cobra.Command
}

func newCreateCmd(config *config.Config) *createCmd {
	root := createCmd{}

	cmd := &cobra.Command{
		Use: "create",
		Short: "Create resource",
	}
	cmd.AddCommand(
		newCreateRepoCmd(config).cmd,
	)

	root.cmd = cmd
	return &root
}


type createRepoCmd struct {
	cmd *cobra.Command
	opts createRepoOpts
}
type createRepoOpts struct {
	repoName string
}

func newCreateRepoCmd(config *config.Config) *createRepoCmd {
	root := createRepoCmd{}

	cmd := &cobra.Command{
		Use: "repo",
		Short: "Create a new repo",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := client.Client(*config)
			if err != nil {
				return err
			}
			if err := populateCreateRepoOpts(&root.opts, args); err != nil {
				return err
			}
			if err := runCreateRepoCmd(client, root.opts, *config); err != nil {
				return err
			}
			return nil
		},
	}

	root.cmd = cmd
	return &root
}

func populateCreateRepoOpts(opts *createRepoOpts, args []string) error {
	opts.repoName = args[0]
	if opts.repoName == "" {
		return fmt.Errorf("Repo name can not be empty")
	}
	return nil
}

func runCreateRepoCmd(client *api.Client, opts createRepoOpts, config config.Config) error {
	result, err := client.PutUserRepo(context.Background(), api.PutUserRepoParams{
		UserID: config.UserID,
		RepoName: opts.repoName,
	})
	if err != nil {
		return err
	}
	switch r := result.(type) {
	case *api.PutUserRepoCreated:
		fmt.Printf("Repo named '%s' created\n", opts.repoName)
	case *api.PutUserRepoConflict:
		return fmt.Errorf("Repo named '%s' already exists", opts.repoName)
	case *api.PutUserRepoUnauthorized:
		return fmt.Errorf("unauthorized")
	case *api.PutUserRepoInternalServerError:
		return fmt.Errorf("server error: %v", r)
	default:
		return fmt.Errorf("unrecognized type %T with result: %v", r, r)
	} 
	return nil
}
