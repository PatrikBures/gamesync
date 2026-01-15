package main

import (
	"fmt"
	"gamesync/internal/config"
	"gamesync/internal/state"
	"gamesync/internal/syncer"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use: "sync GAME_ID",
	Short: "Either pushes or pulls save for GAME_ID",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		game, _, err := config.GetGame(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting game: %v\n", err)
			os.Exit(40)
		}

		stateDir, err := config.GetStateDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting state dir: %v\n", err)
			os.Exit(41)
		}

		localState, err := state.Get(game.SavePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting local state of game with id \"%s\": %v\n", game.ID, err)
			os.Exit(42)
		}

		oldLocalState, err := state.GetOld(game.ID, verbose)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting old local state for game with id \"%s\": %v\n", game.ID, err)
			os.Exit(45)
		}

		remoteState, err := syncer.GetRemoteState(game.ID, config.Current.Server, verbose)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting remote state of game with id \"%s\": %v\n", game.ID, err)
			os.Exit(44)
		}

		updateStateFile := false

		compareResult, err := state.Compare(localState, oldLocalState, remoteState, false, verbose)

		switch compareResult {
		case state.SyncStateConflict:
			fmt.Println("Conflict! Local and remote changes, aborting.")
			os.Exit(100)
		case state.SyncStateUnchanged:
			fmt.Println("Already in sync, nothing to do.")
		case state.SyncStatePush:
			if err := syncer.SyncGame(*game, config.Current.Server, false, verbose); err != nil {
				fmt.Fprintf(os.Stderr, "Error pushing game: %v\n", err)
			}
			fmt.Println("Local changes, pushing.")

			updateStateFile = true
		case state.SyncStatePull:
			fmt.Println("Remote changes, pulling.")

			if err := syncer.SyncGame(*game, config.Current.Server, true, verbose); err != nil {
				fmt.Fprintf(os.Stderr, "Error pulling game: %v\n", err)
			}

			localState, err = state.Get(game.SavePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting local state of game with id \"%s\": %v\n", game.ID, err)
				os.Exit(42)
			}
			updateStateFile = true
		case state.SyncStateError:
			fmt.Fprintf(os.Stderr, "Error comparing states: %v\n", err)
		default:
			panic(fmt.Errorf("unknown state from state.Compare(): %d", compareResult))
		}

		if updateStateFile {
			stateFile := filepath.Join(stateDir, game.ID+".json")
			if err := state.Write(localState, stateFile); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing state for game with id \"%s\":%v\n", game.ID, err)
				os.Exit(43)
			}
			fmt.Println("Updated state file:", stateFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
