package main

import "github.com/spf13/cobra"

var initCmd = &cobra.Command{
	Use: "init",
	Short: "initializes current dir as sync dir",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// check if there is a .gitsync dir in current dir
		// otherwise create that dir
		// create data.json
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
