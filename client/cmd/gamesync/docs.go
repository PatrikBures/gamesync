package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var genDocCmd = &cobra.Command{
	Use: "gen-man",
	Short: "Generate man-pages for program",
	Hidden: true,
	// this prerun is here so that the root prerun does not run
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		dir := "./manpages"
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.Mkdir(dir, 0755)
		}

		header := &doc.GenManHeader{
			Title: "gamesync",
			Section: "1",
			Source: "Auto Generated",
		}

		if err := doc.GenManTree(rootCmd, header, dir); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(genDocCmd)
	rootCmd.DisableAutoGenTag = true
}
