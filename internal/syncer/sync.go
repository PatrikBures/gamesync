package syncer

import (
	"fmt"
	"gamesync/internal/config"
	"gamesync/internal/state"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func SyncGame(game config.GameConfig, server config.ServerConfig, pull bool, verbose bool) error {
	remoteDir := path.Join("data", "saves", server.User, game.ID)
	remotePath := fmt.Sprintf("%s@%s:/%s", server.User, server.Host, path.Clean(remoteDir))
	localPath := filepath.Clean(game.SavePath)

	sshCmd := fmt.Sprintf("ssh -p %s -i %s", server.Port, server.IdentityFile)
	var cmd *exec.Cmd

	var src, dest string

	if pull {
		src = remotePath+"/"
		dest = localPath
	} else {
		if _, err := RunCmd(server, verbose, "mkdir", "-p", "/"+remoteDir); err != nil {
			return fmt.Errorf("failed to create remote dir: %w", err)
		}
		src = localPath+"/"
		dest = remotePath
	}

	cmd = exec.Command("rsync", "-azhPv", "--delete", "-e", sshCmd, src, dest)

	if verbose {
		fmt.Printf("running this rsync command:\n%s\n", cmd.String())
		cmd.Stdout = os.Stdout
	}

	cmd.Stderr = os.Stderr

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

func HandleSync(server config.ServerConfig, gameID string, mode SyncMode, force bool, verbose bool) error {
	game, _, err := config.GetGame(gameID)
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

	oldLocalState, err := state.GetOld(game.ID, verbose)
	if err != nil {
		return fmt.Errorf("getting old local state: %v", err)
	}

	remoteState, err := GetRemoteState(game.ID, config.Current.Server, verbose)
	if err != nil {
		return fmt.Errorf("getting remote state: %v", err)
	}

	compareResult, err := state.Compare(localState, oldLocalState, remoteState, false, verbose)
	if err != nil {
		return fmt.Errorf("comparing states: %v", err)
	}

	var newStateToSave map[string]state.FileMeta = nil

	switch compareResult {
	case state.SyncStateUnchanged:
		fmt.Println("Already in sync, nothing to do.")
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
			fmt.Println("Pushing...")
			if err := push(game, server, verbose); err != nil { return err }
			newStateToSave = localState
		case state.SyncStatePull:
			fmt.Println("Pulling...")
			if err := pull(game, server, verbose); err != nil { return err }
			newStateToSave = remoteState
		}
	case ModePush:
		switch compareResult {
		case state.SyncStateConflict:
			if !force {
				return fmt.Errorf("conflict, remote and local changes")
			}
			fmt.Println("Force pushing over conflict...")
		case state.SyncStatePush:
			fmt.Println("Pushing...")
		case state.SyncStatePull:
			if !force {
				return fmt.Errorf("aborted: unsynced remote changes")
			}
			fmt.Println("Force pushing, overwriting newer remote changes...")
		}
		if err := push(game, server, verbose); err != nil { return err }
		newStateToSave = localState
	case ModePull:
		switch compareResult {
		case state.SyncStateConflict:
			if !force {
				return fmt.Errorf("conflict, remote and local changes")
			}
			fmt.Println("Force pulling over conflict...")
		case state.SyncStatePush:
			if !force {
				return fmt.Errorf("aborted: unsynced local changes")
			}
			fmt.Println("Force pulling, overwriting newer local changes...")
		case state.SyncStatePull:
			fmt.Println("Pulling...")
		}
		if err := pull(game, server, verbose); err != nil { return err }
		newStateToSave = remoteState
	}

	stateFile := filepath.Join(stateDir, game.ID+".json")
	if newStateToSave != nil {
		if err := state.Write(newStateToSave, stateFile); err != nil {
			return fmt.Errorf("writing state to file: %s: %v",stateFile, err)
		}
		fmt.Println("Updated state file:", stateFile)
	}
	return nil
}

func pull(game config.GameConfig, server config.ServerConfig, verbose bool) error {
	return SyncGame(game, server, true, verbose)
}
func push(game config.GameConfig, server config.ServerConfig, verbose bool) error {
	return SyncGame(game, server, false, verbose)
}
