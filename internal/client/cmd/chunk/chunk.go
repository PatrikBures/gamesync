package chunk

import (
	"context"
	"fmt"

	"go.pabu.dev/gamesync/internal/client/config"
	"go.pabu.dev/gamesync/internal/client/profiler"
	"go.pabu.dev/gamesync/internal/snapshoter"

	"github.com/spf13/cobra"
)

type Cmd struct {
	Cmd *cobra.Command
	opts Opts
}

type Opts struct {
	dir string
	argIsDir bool
}

func New(conf *config.Config) *Cmd {
	root := Cmd{}

	cmd := &cobra.Command{
		Use: "chunk PROFILE",
		Short: "Chunks all files in PROFILEs dir",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := populateOpts(conf, &root.opts, args); err != nil {
				return err
			}

			return runCmd(conf, &root.opts)
		},
	}

	cmd.Flags().BoolVarP(&root.opts.argIsDir, "dir", "d", false, "Interperets PROFILE arg as a directory which it will chunk")

	root.Cmd = cmd
	return &root
}

func populateOpts(conf *config.Config, opts *Opts, args []string) error {
	if opts.argIsDir {
		opts.dir = args[0]
		return nil
	}

	profileName := args[0]
	profile, ok, err := profiler.Get(profileName, conf.Global.ProfilesFile)
	if err != nil {
		return fmt.Errorf("getting profile '%s': %w", profileName, err)
	}
	if !ok {
		return fmt.Errorf("profile '%s' not found", profileName)
	}
	opts.dir = profile.Dir
	return nil
}

func runCmd(conf *config.Config, opts *Opts) error {

	fmt.Println("chunking dir:", opts.dir)

	chunkGen := snapshoter.NewChunkGen(conf.Global.ChunkDir)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := chunkGen.ChunkFilesInDir(ctx, opts.dir)
	if err != nil {
		return err
	}
	for fr := range stream.Ch {
		if fr.Err != nil {
			return fr.Err
		}
		chunkGen.ProcessedFile()
	}
	if stream.Err() != nil {
		return stream.Err()
	}

	chunkGen.Info.Print()
	
	return nil
}
