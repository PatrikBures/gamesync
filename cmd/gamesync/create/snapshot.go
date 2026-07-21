package create

import (
	"context"
	"fmt"
	"gamesync/internal/client"
	"gamesync/internal/client/config"
	"gamesync/internal/client/profiler"
	"gamesync/internal/client/syncer"
	api "gamesync/internal/ogen"
	"gamesync/internal/snapshoter"

	"github.com/spf13/cobra"
)

type snapshotCmd struct {
	cmd *cobra.Command
	opts snapshotOpts
}
type snapshotOpts struct {
	profile profiler.Profile
}

func newSnapshotCmd(conf *config.Config) *snapshotCmd {
	root := snapshotCmd{}

	cmd := &cobra.Command{
		Use: "snapshot PROFILE",
		Short: "Create a new snapshot",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := client.New(conf)
			if err != nil {
				return err
			}
			if err := populateSnapshotOpts(conf, &root.opts, args); err != nil {
				return err
			}
			if err := runSnapshotCmd(client, &root.opts, conf); err != nil {
				return err
			}
			return nil
		},
	}

	root.cmd = cmd
	return &root
}

func populateSnapshotOpts(conf *config.Config, opts *snapshotOpts, args []string) error {
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

func runSnapshotCmd(client *api.Client, opts *snapshotOpts, conf *config.Config) error {

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
