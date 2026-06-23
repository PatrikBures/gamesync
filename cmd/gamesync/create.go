package main

import (
	"context"
	"fmt"
	"gamesync/internal/client"
	"gamesync/internal/client/config"
	api "gamesync/internal/ogen"
	"gamesync/internal/snapshoter"

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
		newCreateSnapshotCmd(config).cmd,
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
		return fmt.Errorf("repo name can not be empty")
	}
	return nil
}

func runCreateRepoCmd(client *api.Client, opts createRepoOpts, config config.Config) error {
	result, err := client.PutUserRepo(context.Background(), api.PutUserRepoParams{
		UserID: config.Server.UserID,
		RepoName: opts.repoName,
	})
	if err != nil {
		return err
	}
	switch r := result.(type) {
	case *api.PutUserRepoCreated:
		fmt.Printf("repo named '%s' created\n", opts.repoName)
	case *api.PutUserRepoConflict:
		return fmt.Errorf("repo named '%s' already exists", opts.repoName)
	case *api.PutUserRepoUnauthorized:
		return fmt.Errorf("unauthorized")
	case *api.PutUserRepoInternalServerError:
		return fmt.Errorf("server error: %v", r)
	default:
		return fmt.Errorf("unrecognized type %T with result: %v", r, r)
	} 
	return nil
}



type createSnapshotCmd struct {
	cmd *cobra.Command
	opts createSnapshotOpts
}
type createSnapshotOpts struct {
	repoName string
	branchName string
	dir string
}

func newCreateSnapshotCmd(config *config.Config) *createSnapshotCmd {
	root := createSnapshotCmd{}

	cmd := &cobra.Command{
		Use: "snapshot",
		Short: "Create a new snapshot",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := client.Client(*config)
			if err != nil {
				return err
			}
			if err := populateCreateSnapshotOpts(&root.opts, args); err != nil {
				return err
			}
			if err := runCreateSnapshotCmd(client, root.opts, *config); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&root.opts.repoName, "repo", "r", "", "Repo where snapshot will be created at")
	cmd.Flags().StringVarP(&root.opts.branchName, "branch", "b", "", "Branch for the snapshot")
	cmd.Flags().StringVarP(&root.opts.dir, "dir", "d", "", "Directory which will be snapshoted")

	cmd.MarkFlagRequired("repo")
	cmd.MarkFlagRequired("branch")
	cmd.MarkFlagRequired("dir")

	root.cmd = cmd
	return &root
}

func populateCreateSnapshotOpts(opts *createSnapshotOpts, args []string) error {
	return nil
}

func runCreateSnapshotCmd(client *api.Client, opts createSnapshotOpts, config config.Config) error {

	chunkGen := snapshoter.NewChunkGen(config.Global.ChunkDir)
	files, err := chunkGen.ChunkFilesInDir(opts.dir)
	if err != nil {
		snapshoter.PrintFileResultsErrors(files)
		return fmt.Errorf("generating chunks: %w", err)
	}

	params := api.PostUserRepoBranchSnapshotParams{
		UserID: config.Server.UserID,
		RepoName: opts.repoName,
		BranchName: opts.branchName,
	}

	uploader := snapshoter.NewUploader(client, params, config.Global.ChunkDir, 2)

	if err := uploader.CreateSnapshot(files); err != nil {
		return fmt.Errorf("creating chunk: %w", err)
	}

	return nil
}

