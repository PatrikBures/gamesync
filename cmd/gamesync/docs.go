package main

import (
	"os"

	"gamesync/internal/ui"

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
		Run: func(cmd *cobra.Command, args []string) {
			dir := "./manpages"
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				ui.Error("Error dir does not exist: %s", dir)
				os.Exit(1)
			}

			header := &doc.GenManHeader{
				Title: "gamesync",
				Section: "1",
				Source: "Auto Generated",
			}

			if err := doc.GenManTree(rootCmd, header, dir); err != nil {
				ui.Error("Error generating man-pages: %v\n", err)
				os.Exit(1)
			}
		},
	}

	root.cmd = cmd

	return &root
}

func init() {
	rootCmd.AddCommand(newGenDocCmd().cmd)
	rootCmd.DisableAutoGenTag = true
}
