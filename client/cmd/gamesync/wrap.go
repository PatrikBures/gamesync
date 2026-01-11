package main

import (
	"fmt"
	"gamesync/internal/config"
	"gamesync/internal/syncer"
	"gamesync/internal/ui"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var wrapCreateSnapshot bool
var wrapNotify bool

var wrapCmd = &cobra.Command{
	Use: "wrap GAME_ID -- COMMAND...",
	Short: "Wrap a game process",
	Run: func(cmd *cobra.Command, args []string) {
		dashIdx := cmd.ArgsLenAtDash()

		if dashIdx == -1 {
			fmt.Println("Error: Found no '--'")
			os.Exit(1)
		}

		userArgs := args[:dashIdx]
		cmdArgs := args[dashIdx:]

		if len(userArgs) < 1 {
			fmt.Println("Error: Missing GAME_ID")
			os.Exit(1)
		}
		if len(cmdArgs) < 1 {
			fmt.Println("Error: Missing commands after '--'")
			os.Exit(1)
		}



		gameID := userArgs[0]

		game, err := config.GetGame(gameID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// pulls
		if err := syncer.SyncGame(*game, config.Current.Server, true, verbose); err != nil {
			fmt.Printf("WARNING: Pulling save failed: %v\n", err)
			if wrapNotify { ui.Notify("error", "pulling save") }
		} else {
			fmt.Println("Pulled game")
			if wrapNotify { ui.Notify("sucess", "pulling save") }
		}

		runCmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		runCmd.Stdin = os.Stdin

		exitCode := 0
		if err := runCmd.Run(); err != nil {
			fmt.Printf("Error ")
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
			} else {
				fmt.Printf("Error running game: %v\n", err)
				exitCode = 1
			}
		}

		// pushes
		if err := syncer.SyncGame(*game, config.Current.Server, false, verbose); err != nil {
			fmt.Printf("WARNING: Pushing save failed: %v\n", err)
			if wrapNotify { ui.Notify("error", "pushing save") }
		} else {
			fmt.Println("Pushed game")
			if wrapNotify { ui.Notify("sucess", "pushing save") }
		}

		if wrapCreateSnapshot {
			if err := syncer.CreateSnapshot(config.Current.Server, verbose, gameID); err != nil {
				fmt.Printf("Failed creating snapshot: %v\n", err)
				if wrapNotify { ui.Notify("error", "creating snapshot") }
			} else {
				fmt.Println("Snapshot created")
				if wrapNotify { ui.Notify("sucess", "creating snapshot") }
			}
		}

		fmt.Printf("Game exited with exit code: %d\n", exitCode)

		os.Exit(exitCode)
	},
}

func init() {
	wrapCmd.Flags().BoolVarP(&wrapCreateSnapshot, "snapshot", "s", false, "")
	wrapCmd.Flags().BoolVarP(&wrapNotify, "notify", "n", false, "")
	rootCmd.AddCommand(wrapCmd)
}
