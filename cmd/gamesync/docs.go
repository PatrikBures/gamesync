package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

type genDocCmd struct {
	cmd *cobra.Command
}

func newGenDocCmd() *genDocCmd {
	root := genDocCmd{}

	cmd := &cobra.Command{
		Use: "gen-man",
		Short: "Generate man-pages for program",
		Hidden: true,
		// this prerun is here so that the root prerun does not run
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "./manpages"
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				return fmt.Errorf("error dir does not exist: %s", dir)
			}

			header := &doc.GenManHeader{
				Title: "gamesync",
				Section: "1",
				Source: "Auto Generated",
			}

			if err := doc.GenManTree(cmd.Root(), header, dir); err != nil {
				return fmt.Errorf("error generating man-pages: %v", err)
			}

			return nil
		},
	}

	root.cmd = cmd

	return &root
}
