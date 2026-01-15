package syncer

import (
	"fmt"
	"gamesync/internal/config"
	"gamesync/internal/state"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const (
	stateSynced = iota
	stateSourceNewer
	stateSourceHasExtra
	stateError
)

func SyncGame(game config.GameConfig, server config.ServerConfig, pull bool, verbose bool) error {
	remoteDir := path.Join("data", "saves", server.User, game.ID)
	remotePath := fmt.Sprintf("%s@%s:/%s/", server.User, server.Host, remoteDir)
	localPath := ensureTrailingSep(game.SavePath)

	sshCmd := fmt.Sprintf("ssh -p %s -i %s", server.Port, server.IdentityFile)
	var cmd *exec.Cmd

	if pull {
		state, err := checkSyncState(remotePath, localPath[:len(localPath)-1], sshCmd, verbose)
		if err != nil {
			return fmt.Errorf("checking remote state: %w", err)
		}

		switch state {
		case stateSynced:
			fmt.Println("Already up to date")
			return nil
		}

		flags := "-azhP"
		if verbose {
			flags = flags+"v"
			fmt.Println("Pulling save...")
		}

		cmd = exec.Command("rsync", flags, "--delete", "-e", sshCmd, remotePath, localPath[:len(localPath)-1])
	} else {
		_, err := RunCmd(server, verbose, "mkdir", "-p", "/"+remoteDir)

		if err != nil {
			return fmt.Errorf("failed to create remote dir: %w", err)
		}

		remoteState, err := checkSyncState(remotePath, localPath[:len(localPath)-1], sshCmd, verbose)
		if err != nil {
			return fmt.Errorf("checking remote state: %w", err)
		}

		if remoteState == stateSourceNewer {
			return  fmt.Errorf("can not push remote as remote is newer than local save")
		}

		localState, err := checkSyncState(localPath, remotePath[:len(remotePath)-1], sshCmd, verbose)
		if err != nil {
			return fmt.Errorf("checking local state: %w", err)
		}

		if localState == stateSynced && remoteState == stateSynced {
			fmt.Println("Remote save is already up to date.")
			return nil
		}

		flags := "-azhP"
		if verbose {
			flags = flags+"v"
			fmt.Println("Pushing save...")
		}

		cmd = exec.Command("rsync", flags, "--delete", "-e", sshCmd, localPath, remotePath[:len(remotePath)-1])
	}


	if verbose {
		fmt.Printf("running this rsync command:\n%s\n", cmd.String())
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error syncing %s, from %s, to %s", game.ID, localPath, remotePath)
	}

	return nil
}

func checkSyncState(src string, dest string, sshCmd string, verbose bool) (int, error) {
	cmd := exec.Command("rsync", "-naui", "-e", sshCmd, src, dest)

	outputBytes, err := cmd.CombinedOutput()
	output := string(outputBytes)

	if verbose {
		fmt.Printf("checking sync state with this command:\n%s\n\n", cmd.String())
	}

	if err != nil {
		fmt.Println(output)
		return stateError, err
	}

	lines := strings.Split(output, "\n")
	hasNewerFiles := false
	hasExtraFiles := false

	for _, line := range lines {
		if len(line) < 12 { continue }

		parts := strings.Fields(line)
		code := parts[0]

		if len(code) < 2 || code[1] != 'f' {
			continue
		}

		if code[0] != '>' && code[0] != '<' {
			continue
		}

		if code[2] == '+' {
			hasExtraFiles = true
		} else {
			hasNewerFiles = true
		}
	}

	if hasNewerFiles {
		return stateSourceNewer, nil
	}
	if hasExtraFiles {
		return stateSourceHasExtra, nil
	}

	return stateSynced, nil
}

func ensureTrailingSep(p string) string {
	if !strings.HasSuffix(p, string(filepath.Separator)) {
		return p + string(filepath.Separator)
	}
	return p
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
		return fmt.Errorf("Error getting game: %v\n", err)
	}

	stateDir, err := config.GetStateDir()
	if err != nil {
		return fmt.Errorf("Error getting state dir: %v\n", err)
	}

	localState, err := state.Get(game.SavePath)
	if err != nil {
		return fmt.Errorf("Error getting local state of game with id \"%s\": %v\n", game.ID, err)
	}

	oldLocalState, err := state.GetOld(game.ID, verbose)
	if err != nil {
		return fmt.Errorf("Error getting old local state for game with id \"%s\": %v\n", game.ID, err)
	}

	remoteState, err := GetRemoteState(game.ID, config.Current.Server, verbose)
	if err != nil {
		return fmt.Errorf("Error getting remote state of game with id \"%s\": %v\n", game.ID, err)
	}

	compareResult, err := state.Compare(localState, oldLocalState, remoteState, false, verbose)
	if err != nil {
		return fmt.Errorf("comparing states: %v\n", err)
	}

	var newStateToSave map[string]state.FileMeta = nil

	switch compareResult {
	case state.SyncStateUnchanged:
		fmt.Println("Already in sync, nothing to do.")
		return nil
	case state.SyncStateError:
		return fmt.Errorf("error during state comparison: %v\n", err)
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
			return fmt.Errorf("writing state to file: %s: %v\n",stateFile, err)
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
