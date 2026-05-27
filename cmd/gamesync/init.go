package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)


type initCmd struct {
	cmd *cobra.Command
	opts initCmdOpts
}

type initCmdOpts struct {
	name string
	dir string
}

func populateInitOpts(opts *initCmdOpts, args []string) error {
	opts.name = args[0]
	if len(args) > 1 {
		abs, err := filepath.Abs(args[1])
		if err != nil {
			return err
		}
		opts.dir = abs
	} else {
		var err error
		opts.dir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	s, err := os.Stat(opts.dir)
	if err != nil {
		return err
	}
	if !s.IsDir() {
		return fmt.Errorf("is not dir")
	}

	return nil
}

func newInitCmd() *initCmd {
	root := initCmd{}

	cmd := &cobra.Command{
		Use: "init NAME [DIR]",
		Short: "Initialize repo",
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := populateInitOpts(&root.opts, args); err != nil {
				return fmt.Errorf("populating opts: %w", err)
			}

			fmt.Println("name:", root.opts.name)
			fmt.Println("dir:", root.opts.dir)

			return nil
		},
	}

	root.cmd = cmd
	return &root
}
