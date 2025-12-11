package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "gamesync",
	Short: "A CLI to sync save games to a server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GameSync!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
