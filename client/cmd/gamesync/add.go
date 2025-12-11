package main

import "github.com/spf13/cobra"


var addCmd = &cobra.Command{
	Use: "add [game of game]",
	Short: "Adds game to sync and creates symlink",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
