package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "gamesync",
		Short: "A CLI to sync save games to a server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("GameSync!")
		},
	}

	rootCmd.AddCommand(playCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
