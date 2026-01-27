package main

import (
	"fmt"
	"os"
	"os/exec"

	"gamesync/internal/syncer"
	"gamesync/internal/ui"

	"github.com/spf13/cobra"
)

type wrapCmd struct {
	cmd			*cobra.Command
	createSnapshot bool
	createSnapshotUnchanged bool
	notify		bool
	noPull		bool
	forcePull	bool
	noPush		bool
	forcePush	bool
}

func newWrapCmd() *wrapCmd {
	root := wrapCmd{}

	cmd := &cobra.Command{
		Use: "wrap GAME_ID -- COMMAND...",
		Short: "Wrap a game process",
		RunE: func(cmd *cobra.Command, args []string) error {
			dashIdx := cmd.ArgsLenAtDash()

			if dashIdx == -1 {
				return fmt.Errorf("found no '--'")
			}

			userArgs := args[:dashIdx]
			cmdArgs := args[dashIdx:]

			if len(userArgs) < 1 {
				return fmt.Errorf("missing GAME_ID")
			}
			if len(cmdArgs) < 1 {
				return fmt.Errorf("missing command after '--'")
			}



			gameID := userArgs[0]

			// pulls
			if !root.noPull {
				if err := syncer.HandleSync(current, gameID, syncer.ModePull, root.forcePull); err != nil {
					fmt.Printf("WARNING: Pulling save failed: %v\n", err)
					if root.notify { ui.Notify("error", "pulling save") }
				} else {
					fmt.Println("Pulled game")
					if root.notify { ui.Notify("sucess", "pulling save") }
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
			if !root.noPush {
				if err := syncer.HandleSync(current, gameID, syncer.ModePush, root.forcePush); err != nil {
					fmt.Printf("WARNING: Pushing save failed: %v\n", err)
					if root.notify { ui.Notify("error", "pushing save") }
				} else {
					fmt.Println("Pushed game")
					if root.notify { ui.Notify("sucess", "pushing save") }
				}
			}

			// snapshot
			if root.createSnapshot || root.createSnapshotUnchanged {
				if err := syncer.CreateSnapshot(current, gameID, root.createSnapshotUnchanged); err != nil {
					fmt.Printf("Failed creating snapshot: %v\n", err)
					if root.notify { ui.Notify("error", "creating snapshot") }
				} else {
					fmt.Println("Snapshot created")
					if root.notify { ui.Notify("sucess", "creating snapshot") }
				}
			}

			fmt.Printf("Game exited with exit code: %d\n", exitCode)

			return nil
		},
	}


	cmd.Flags().BoolVarP(&root.createSnapshot, "snapshot", "s", false, "Creates snapshot on remote after push")
	cmd.Flags().BoolVarP(&root.createSnapshotUnchanged, "skip-unchanged", "S", false, "Creates snapshot if there were changes from the previous snapshot on remote after push")
	cmd.Flags().BoolVarP(&root.notify, "notify", "n", false, "Sends a notification when pulled, pushed and created a snapshot and if succeeded")

	cmd.Flags().BoolVarP(&root.noPull, "no-pull", "", false, "")
	cmd.Flags().BoolVarP(&root.forcePull, "force-pull", "", false, "")

	cmd.Flags().BoolVarP(&root.noPush, "no-push", "", false, "")
	cmd.Flags().BoolVarP(&root.forcePull, "force-push","" , false, "")

	root.cmd = cmd

	return &root
}

func init() {
	rootCmd.AddCommand(newWrapCmd().cmd)
}
