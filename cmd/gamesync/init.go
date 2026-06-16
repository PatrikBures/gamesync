package main

import (
	"fmt"
	"gamesync/internal/client/config"
	"gamesync/internal/snapshoter"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type initCmd struct {
	cmd  *cobra.Command
	opts initCmdOpts
}

type initCmdOpts struct {
	name    string
	repoDir string
}

func newInitCmd(configOpts *config.Config) *initCmd {
	root := initCmd{}

	cmd := &cobra.Command{
		Use:   "init NAME [DIR]",
		Short: "Initialize repo",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := populateInitOpts(&root.opts, args); err != nil {
				return fmt.Errorf("populating init opts: %w", err)
			}

			if err := runInitCmd(root.opts, configOpts.ChunkDir); err != nil {
				return fmt.Errorf("initializing repo: %w", err)
			}

			return nil
		},
	}

	root.cmd = cmd
	return &root
}

func populateInitOpts(opts *initCmdOpts, args []string) error {
	opts.name = args[0]
	if len(args) > 1 {
		abs, err := filepath.Abs(args[1])
		if err != nil {
			return err
		}
		opts.repoDir = abs
	} else {
		var err error
		opts.repoDir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	s, err := os.Stat(opts.repoDir)
	if err != nil {
		return err
	}
	if !s.IsDir() {
		return fmt.Errorf("repo not dir")
	}

	return nil
}

func runInitCmd(opts initCmdOpts, chunkDir string) error {
	cg := snapshoter.NewChunkGen(chunkDir)
	result, err := cg.ChunkFilesInDir(opts.repoDir)
	if err != nil {
		if result != nil {
			snapshoter.PrintFileResultsErrors(result)
		}
		return err
	}

	for _, r := range result {
		fmt.Printf("%d chunks from %s\n", len(r.Hashes), r.Path)
	}

	fmt.Println("files chunked:", cg.Info.FilesChunked)
	fmt.Println("files error:", cg.Info.FilesErr)
	fmt.Println("chunks created:", cg.Info.ChunksCreated)
	fmt.Println("chunks skipped:", cg.Info.ChunksSkipped)

	return nil
}
