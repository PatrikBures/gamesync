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

var wrapCreateSnapshot, 
	wrapCreateSnapshotUnchanged,
	wrapNotify,
	wrapNoPull,
	wrapForcePull,
	wrapNoPush,
	wrapForcePush bool

var wrapCmd = &cobra.Command{
	Use: "wrap GAME_ID -- COMMAND...",
	Short: "Wrap a game process",
	Run: func(cmd *cobra.Command, args []string) {
		dashIdx := cmd.ArgsLenAtDash()

		if dashIdx == -1 {
			fmt.Fprintf(os.Stderr, "Error: Found no '--'\n")
			os.Exit(1)
		}

		userArgs := args[:dashIdx]
		cmdArgs := args[dashIdx:]

		if len(userArgs) < 1 {
			fmt.Fprintf(os.Stderr, "Error: Missing GAME_ID\n")
			os.Exit(1)
		}
		if len(cmdArgs) < 1 {
			fmt.Fprintf(os.Stderr, "Error: Missing commands after '--'\n")
			os.Exit(1)
		}



		gameID := userArgs[0]

		// pulls
		if !wrapNoPull {
			if err := syncer.HandleSync(config.Current.Server, gameID, syncer.ModePull, wrapForcePull, verbose); err != nil {
				fmt.Printf("WARNING: Pulling save failed: %v\n", err)
				if wrapNotify { ui.Notify("error", "pulling save") }
			} else {
				fmt.Println("Pulled game")
				if wrapNotify { ui.Notify("sucess", "pulling save") }
			}
		}

		// runs command
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
		if !wrapNoPush {
			if err := syncer.HandleSync(config.Current.Server, gameID, syncer.ModePush, wrapForcePush, verbose); err != nil {
				fmt.Printf("WARNING: Pushing save failed: %v\n", err)
				if wrapNotify { ui.Notify("error", "pushing save") }
			} else {
				fmt.Println("Pushed game")
				if wrapNotify { ui.Notify("sucess", "pushing save") }
			}
		}

		// snapshot
		if wrapCreateSnapshot || wrapCreateSnapshotUnchanged {
			if err := syncer.CreateSnapshot(config.Current.Server, verbose, gameID, wrapCreateSnapshotUnchanged); err != nil {
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
	wrapCmd.Flags().BoolVarP(&wrapCreateSnapshot, "snapshot", "s", false, "Creates snapshot on remote after push")
	wrapCmd.Flags().BoolVarP(&wrapCreateSnapshotUnchanged, "skip-unchanged", "S", false, "Creates snapshot if there were changes from the previous snapshot on remote after push")
	wrapCmd.Flags().BoolVarP(&wrapNotify, "notify", "n", false, "Sends a notification when pulled, pushed and created a snapshot and if succeeded")

	wrapCmd.Flags().BoolVarP(&wrapNoPull, "no-pull", "", false, "")
	wrapCmd.Flags().BoolVarP(&wrapForcePull, "force-pull", "", false, "")

	wrapCmd.Flags().BoolVarP(&wrapNoPush, "no-push", "", false, "")
	wrapCmd.Flags().BoolVarP(&wrapForcePull, "force-push","" , false, "")

	rootCmd.AddCommand(wrapCmd)
}
