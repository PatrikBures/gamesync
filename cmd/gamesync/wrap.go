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
	cmd *cobra.Command
	opts wrapOpts
}

type wrapOpts struct {
	notify			bool
	noPull			bool
	forcePull		bool
	noPush			bool
	forcePush		bool
	exitOnError 	bool
	handleConflict  bool
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
			var pullHandled syncer.HandledChoice
			if !root.opts.noPull {
				var err error
				pullHandled, err = syncer.HandleSync(current, gameID, syncer.ModePull, root.opts.forcePull, true, root.opts.handleConflict)
				if err != nil {
					ui.Info("WARNING: Pulling save failed: %v\n", err)
					if root.opts.notify { ui.Notify("error", "pulling save") }
					if root.opts.exitOnError { return err }
				} else {
					if root.opts.notify { ui.Notify("sucess", "pulling save") }
					ui.Info("Synced save for %s\n", gameID)
				}
				if pullHandled == syncer.HandledCancel {
					return fmt.Errorf("cancelled launch of %s", gameID)
				}
			}

			// runs command
			runCmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

			runCmd.Stdout = os.Stdout
			runCmd.Stderr = os.Stderr
			runCmd.Stdin = os.Stdin

			exitCode := 0
			if err := runCmd.Run(); err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					exitCode = exitError.ExitCode()
				} else {
					ui.Error("Error running game: %v\n", err)
					if root.opts.exitOnError { return err }
					exitCode = 1
				}
			}

			// pushes
			if !root.opts.noPush {

				pullHandleConflict := root.opts.handleConflict
				if pullHandled == syncer.HandledIgnore {
					pullHandleConflict = false
				}

				pushHandled, err := syncer.HandleSync(current, gameID, syncer.ModePush, root.opts.forcePush, true, pullHandleConflict)
				if ; err != nil {
					ui.Info("WARNING: Pushing save failed: %v\n", err)
					if root.opts.notify { ui.Notify("error", "pushing save") }
					if root.opts.exitOnError { return err }
				} else {
					if root.opts.notify { ui.Notify("sucess", "pushing save") }
					ui.Info("Synced save for %s\n", gameID)
				}
				if pushHandled == syncer.HandledCancel {
					return fmt.Errorf("cancelled push")
				}
			}

			ui.Info("Game exited with exit code: %d\n", exitCode)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&root.opts.exitOnError, "exit-on-error", "e", false, "")

	cmd.Flags().BoolVarP(&root.opts.noPull, "no-pull", "", false, "")
	cmd.Flags().BoolVarP(&root.opts.forcePull, "force-pull", "", false, "")

	cmd.Flags().BoolVarP(&root.opts.noPush, "no-push", "", false, "")
	cmd.Flags().BoolVarP(&root.opts.forcePull, "force-push","" , false, "")

	cmd.Flags().BoolVarP(&root.opts.handleConflict, "handle-conflict", "", false, "Opens a ui to let user pick handle method")

	root.cmd = cmd

	return &root
}
