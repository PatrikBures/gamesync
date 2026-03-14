package syncer

import (
	"fmt"
	"path"
	"path/filepath"

	"gamesync/internal/config"
	"gamesync/internal/state"
	"gamesync/internal/ui"
)

func SyncGame(current config.Current, game config.GameConfig, pull bool) error {
	remotePath := "/"+path.Join("data", "saves", current.Config.Server.User, game.ID)
	localPath := filepath.Clean(game.SavePath)

	if pull {
		remotePath = remotePath+"/"
	} else {
		if _, err := RunCmd(current.Config.Server, false, "mkdir", "-p", remotePath); err != nil {
			return fmt.Errorf("failed to create remote dir: %w", err)
		}
		localPath = localPath+"/"
	}

	if err := RunRsync(current.Config.Server, localPath, remotePath, !pull, "-azhPv", "--delete"); err != nil {
		return fmt.Errorf("syncing %s: %w", game.ID, err)
	}

	return nil
}

type SyncMode int
const (
	ModeAuto SyncMode = iota
	ModePush
	ModePull
)

type HandledChoice int
const (
	HandledIgnore HandledChoice = iota
	HandledCancel
	HandledResolve
	HandledError
)

func HandleSync(current config.Current, gameID string, mode SyncMode, force bool, updateState bool, handleConflict bool) (HandledChoice, error) {
	game, _, err := config.GetGame(current, gameID)
	if err != nil {
		return HandledError, fmt.Errorf("getting game: %v", err)
	}

	stateDir, err := config.GetStateDir()
	if err != nil {
		return HandledError, fmt.Errorf("getting state dir: %v", err)
	}

	if err := SameApiVersion(current.Config.Server); err != nil {
		return HandledError, err
	}

	localState, err := state.Get(game.SavePath)
	if err != nil {
		return HandledError, fmt.Errorf("getting local state: %v", err)
	}

	oldLocalState, err := state.GetOld(current, game.ID)
	if err != nil {
		return HandledError, fmt.Errorf("getting old local state: %v", err)
	}

	remoteState, err := GetRemoteState(current, game.ID)
	if err != nil {
		return HandledError, fmt.Errorf("getting remote state: %v", err)
	}

	compareResult, err := state.Compare(localState, oldLocalState, remoteState, false)
	if err != nil {
		return HandledError, fmt.Errorf("comparing states: %v", err)
	}

	var newStateToSave map[string]state.FileMeta = nil

	switch compareResult {
	case state.SyncStateUnchanged:
		if ! force {
			ui.Info("Already in sync, nothing to do.\n")
			return HandledResolve, nil
		}
	case state.SyncStateError:
		return HandledError, fmt.Errorf("error during state comparison: %v", err)
	}

	var conflictErr error

	switch mode {
	case ModeAuto:
		switch compareResult {
		case state.SyncStateConflict:
			conflictErr = fmt.Errorf("conflict, remote and local changes")
		case state.SyncStatePush:
			ui.Info("Pushing...\n")
			if err := push(current, game); err != nil { return HandledError, err }
			newStateToSave = localState
		case state.SyncStatePull:
			ui.Info("Pulling...\n")
			if err := pull(current, game); err != nil { return HandledError, err }
			newStateToSave = remoteState
		}
	case ModePush:
		switch compareResult {
		case state.SyncStateConflict:
			if !force {
				conflictErr = fmt.Errorf("conflict, remote and local changes")
			} else {
				ui.Info("Force pushing over conflict...\n")
			}
		case state.SyncStatePush:
			ui.Info("Pushing...\n")
		case state.SyncStatePull:
			if !force {
				return HandledError, fmt.Errorf("aborted: unsynced remote changes")
			}
			ui.Info("Force pushing, overwriting newer remote changes...\n")
		}
		if conflictErr == nil {
			if err := push(current, game); err != nil { return HandledError, err }
			newStateToSave = localState
		}
	case ModePull:
		switch compareResult {
		case state.SyncStateConflict:
			if !force {
				conflictErr = fmt.Errorf("conflict, remote and local changes")
			} else {
				ui.Info("Force pulling over conflict...\n")
			}
		case state.SyncStatePush:
			if !force {
				return HandledError, fmt.Errorf("aborted: unsynced local changes")
			}
			ui.Info("Force pulling, overwriting newer local changes...\n")
		case state.SyncStatePull:
			ui.Info("Pulling...\n")
		}
		if conflictErr == nil {
			if err := pull(current, game); err != nil { return HandledError, err }
			newStateToSave = remoteState
		}
	}

	if conflictErr != nil && ! handleConflict {
		return HandledError, conflictErr
	} else if conflictErr != nil && handleConflict {
		latestLocalModTime  := state.LatestModTime(localState)
		latestRemoteModTime := state.LatestModTime(remoteState)

		solution, err := ui.DialogConflict(gameID, latestLocalModTime, latestRemoteModTime)
		if err != nil {
			return HandledError, err
		}
		switch solution {
		case ui.ConflictPush:
			ui.Info("Resolved conflict by force pushing\n")
			if err := push(current, game); err != nil { return HandledResolve, err }
			newStateToSave = localState
		case ui.ConflictPull:
			ui.Info("Resolved conflict by force pulling\n")
			if err := pull(current, game); err != nil { return HandledResolve, err }
			newStateToSave = remoteState
		case ui.ConflictCancel:
			ui.Verbose("Cancelled conflict\n")
			return HandledCancel, nil
		case ui.ConflictIgnore:
			ui.Verbose("Ignored conflict\n")
			return HandledIgnore, nil
		case ui.ConflictError:
			return HandledError, fmt.Errorf("somehow returned an error from dialog conflict")
		default:
			return HandledError, fmt.Errorf("unhandled Conflict enum")
		}
	}

	if newStateToSave != nil && updateState {
		stateFile := filepath.Join(stateDir, game.ID+".json")
		if err := state.Write(newStateToSave, stateFile); err != nil {
			return HandledError, fmt.Errorf("writing state to file: %s: %v",stateFile, err)
		}
		ui.Info("Updated state file: %s\n", stateFile)
	}
	return HandledResolve, nil
}

func pull(current config.Current, game config.GameConfig) error {
	return SyncGame(current, game, true)
}
func push(current config.Current, game config.GameConfig) error {
	return SyncGame(current, game, false)
}
