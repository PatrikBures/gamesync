package syncer

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"gamesync/internal/config"
	"gamesync/internal/state"
	"gamesync/internal/ui"
)

func SyncGame(current config.Current, game config.GameConfig, pull bool) error {
	remoteDir := path.Join("data", "saves", current.Config.Server.User, game.ID)
	remotePath := fmt.Sprintf("%s@%s:/%s", current.Config.Server.User, current.Config.Server.Host, path.Clean(remoteDir))
	localPath := filepath.Clean(game.SavePath)

	sshCmd := fmt.Sprintf("ssh -p %s -i %s", current.Config.Server.Port, current.Config.Server.IdentityFile)
	var cmd *exec.Cmd

	var src, dest string

	if pull {
		src = remotePath+"/"
		dest = localPath
	} else {
		if _, err := RunCmd(current, "mkdir", "-p", "/"+remoteDir); err != nil {
			return fmt.Errorf("failed to create remote dir: %w", err)
		}
		src = localPath+"/"
		dest = remotePath
	}

	cmd = exec.Command("rsync", "-azhPv", "--delete", "-e", sshCmd, src, dest)

	ui.Debug("running rsync command:\n%s\n", cmd.String())
	if ui.GetLevel() == ui.LevelDebug {
		cmd.Stdout = os.Stdout
	}

	if ui.GetLevel() <= ui.LevelError {
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error syncing %s, from %s, to %s", game.ID, localPath, remotePath)
	}

	return nil
}

type SyncMode int
const (
	ModeAuto SyncMode = iota
	ModePush
	ModePull
)

func HandleSync(current config.Current, gameID string, mode SyncMode, force bool) error {
	game, _, err := config.GetGame(current, gameID)
	if err != nil {
		return fmt.Errorf("getting game: %v", err)
	}

	stateDir, err := config.GetStateDir()
	if err != nil {
		return fmt.Errorf("getting state dir: %v", err)
	}

	localState, err := state.Get(game.SavePath)
	if err != nil {
		return fmt.Errorf("getting local state: %v", err)
	}

	oldLocalState, err := state.GetOld(current, game.ID)
	if err != nil {
		return fmt.Errorf("getting old local state: %v", err)
	}

	remoteState, err := GetRemoteState(current, game.ID)
	if err != nil {
		return fmt.Errorf("getting remote state: %v", err)
	}

	compareResult, err := state.Compare(localState, oldLocalState, remoteState, false)
	if err != nil {
		return fmt.Errorf("comparing states: %v", err)
	}

	var newStateToSave map[string]state.FileMeta = nil

	switch compareResult {
	case state.SyncStateUnchanged:
		ui.Info("Already in sync, nothing to do.\n")
		return nil
	case state.SyncStateError:
		return fmt.Errorf("error during state comparison: %v", err)
	}

	switch mode {
	case ModeAuto:
		switch compareResult {
		case state.SyncStateConflict:
			return fmt.Errorf("conflict, remote and local changes")
		case state.SyncStatePush:
			ui.Info("Pushing...\n")
			if err := push(current, game); err != nil { return err }
			newStateToSave = localState
		case state.SyncStatePull:
			ui.Info("Pulling...\n")
			if err := pull(current, game); err != nil { return err }
			newStateToSave = remoteState
		}
	case ModePush:
		switch compareResult {
		case state.SyncStateConflict:
			if !force {
				return fmt.Errorf("conflict, remote and local changes")
			}
			ui.Info("Force pushing over conflict...\n")
		case state.SyncStatePush:
			ui.Info("Pushing...\n")
		case state.SyncStatePull:
			if !force {
				return fmt.Errorf("aborted: unsynced remote changes")
			}
			ui.Info("Force pushing, overwriting newer remote changes...\n")
		}
		if err := push(current, game); err != nil { return err }
		newStateToSave = localState
	case ModePull:
		switch compareResult {
		case state.SyncStateConflict:
			if !force {
				return fmt.Errorf("conflict, remote and local changes")
			}
			ui.Info("Force pulling over conflict...\n")
		case state.SyncStatePush:
			if !force {
				return fmt.Errorf("aborted: unsynced local changes")
			}
			ui.Info("Force pulling, overwriting newer local changes...\n")
		case state.SyncStatePull:
			ui.Info("Pulling...\n")
		}
		if err := pull(current, game); err != nil { return err }
		newStateToSave = remoteState
	}

	stateFile := filepath.Join(stateDir, game.ID+".json")
	if newStateToSave != nil {
		if err := state.Write(newStateToSave, stateFile); err != nil {
			return fmt.Errorf("writing state to file: %s: %v",stateFile, err)
		}
		ui.Info("Updated state file: %s\n", stateFile)
	}
	return nil
}

func pull(current config.Current, game config.GameConfig) error {
	return SyncGame(current, game, true)
}
func push(current config.Current, game config.GameConfig) error {
	return SyncGame(current, game, false)
}
