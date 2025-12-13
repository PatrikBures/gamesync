package main

import (
	"fmt"
	"gamesync/internal/config"
	"os"

	"github.com/spf13/cobra"
)

var poolCmd = &cobra.Command{
	Use: "pool <cmd>",
	Short: "Manage pools",
}

var poolInitCmd = &cobra.Command{
	Use: "init <pool_id> [pool_dir]",
	Short: "initalize a pool",
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {

		numArgs := len(args)
		poolId := args[0]
		var poolDir string

		if numArgs == 1 {
			var err error
			poolDir, err = os.Getwd()
			if err != nil {
				fmt.Printf("Failed getting working directory, %v", err)
				os.Exit(1)
			}
		} else {
			poolDir = args[1]
		}

		update, _ := cmd.Flags().GetBool("update")
		configPath := ""
		if update {
			configPath = configFile
		}
		err := config.InitPool(poolId, poolDir, configPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	poolCmd.AddCommand(poolInitCmd)
	rootCmd.AddCommand(poolCmd)
	poolInitCmd.Flags().BoolP("update", "u", false, "Updates global config")
}
