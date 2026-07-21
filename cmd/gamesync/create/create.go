package create

import (
	"context"
	"errors"
	"fmt"
	util "gamesync/cmd/gamesync/_util"
	"gamesync/internal/client"
	"gamesync/internal/client/config"
	"gamesync/internal/client/profiler"
	"gamesync/internal/client/syncer"
	api "gamesync/internal/ogen"
	"gamesync/internal/snapshoter"

	"github.com/spf13/cobra"
)


type createCmd struct {
	Cmd *cobra.Command
}

func New(config *config.Config) *createCmd {
	root := createCmd{}

	cmd := &cobra.Command{
		Use: "create",
		Short: "Create resource",
	}
	cmd.AddCommand(
		newCreateRepoCmd(config).cmd,
		newCreateSnapshotCmd(config).cmd,
		newCreateProfileCmd(config).cmd,
	)

	root.Cmd = cmd
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
			client, err := client.New(config)
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
	err := client.PutUserRepo(context.Background(), api.PutUserRepoParams{
		UserID: config.Server.UserID,
		RepoName: opts.repoName,
	})
	if err := util.ErrHandler(err); err != nil {
		return err
	}
	fmt.Printf("repo named '%s' created\n", opts.repoName)
	return nil
}



type createSnapshotCmd struct {
	cmd *cobra.Command
	opts createSnapshotOpts
}
type createSnapshotOpts struct {
	profile profiler.Profile
}

func newCreateSnapshotCmd(conf *config.Config) *createSnapshotCmd {
	root := createSnapshotCmd{}

	cmd := &cobra.Command{
		Use: "snapshot PROFILE",
		Short: "Create a new snapshot",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := client.New(conf)
			if err != nil {
				return err
			}
			if err := populateCreateSnapshotOpts(conf, &root.opts, args); err != nil {
				return err
			}
			if err := runCreateSnapshotCmd(client, &root.opts, conf); err != nil {
				return err
			}
			return nil
		},
	}

	root.cmd = cmd
	return &root
}

func populateCreateSnapshotOpts(conf *config.Config, opts *createSnapshotOpts, args []string) error {
	profile, ok, err := profiler.Get(args[0], conf.Global.ProfilesFile)
	if err != nil {
		return fmt.Errorf("initializing profiler: %s", err)
	}
	if !ok {
		return fmt.Errorf("profile '%s' does not exist", args[0])
	}
	opts.profile = profile

	return nil
}

func runCreateSnapshotCmd(client *api.Client, opts *createSnapshotOpts, conf *config.Config) error {

	chunkGen := snapshoter.NewChunkGen(conf.Global.ChunkDir)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	files, err := chunkGen.ChunkFilesInDirApiFile(ctx, opts.profile.Dir)
	if err != nil {
		return err
	}

	syncer := syncer.New(conf, client, opts.profile)

	if err := syncer.CreateSnapshot(files); err != nil {
		return fmt.Errorf("creating snapshot: %w", err)
	}

	return nil
}



type createProfileCmd struct {
	cmd *cobra.Command
	opts createProfileOpts
}
type createProfileOpts struct {
	slug    string
	force   bool
	profile profiler.Profile
}

func newCreateProfileCmd(conf *config.Config) *createProfileCmd {
	root := createProfileCmd{}

	cmd := &cobra.Command{
		Use: "profile PROFILE REPO BRANCH DIR",
		Short: "Create new profile to sync",
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			populateCreateProfileOpts(&root.opts, args)

			if err := runCreateProfileCmd(root.opts, conf); err != nil { return err }

			return nil
		},
	}
	cmd.Flags().BoolVarP(&root.opts.force, "force", "f", false, "Overwrites any exising profile")
	root.cmd = cmd
	return &root
}

func populateCreateProfileOpts(opts *createProfileOpts, args []string) {
	opts.slug = args[0]
	opts.profile = profiler.Profile{
		RepoName: args[1],
		BranchName: args[2],
		Dir: args[3],
	}
}

func runCreateProfileCmd(opts createProfileOpts, conf *config.Config) (err error) {
	pr, err := profiler.New(conf.Global.ProfilesFile)
	if err != nil {
		return fmt.Errorf("initializing profiler: %w", err)
	}
	defer func() {
		err = errors.Join(err, pr.Close())
	}()

	if opts.force {
		pr.AddOverwrite(opts.slug, opts.profile)
	} else {
		if err := pr.Add(opts.slug, opts.profile); err != nil { return err }
	}

	if opts.force && pr.Exists(opts.slug) {
		fmt.Printf("Force created profile '%s'\n", opts.slug)
	} else {
		fmt.Printf("Created profile '%s'\n", opts.slug)
	}

	if err := pr.Save(); err != nil { return err }

	return nil
}
