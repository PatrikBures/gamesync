package syncer

import (
	"fmt"
	"gamesync/internal/config"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	statusChange = iota
	statusSyncedOrOlder
	statusMissing
)

func SyncGame(game config.GameConfig, server config.ServerConfig, pull bool) error {
	remotePath := filepath.Join(fmt.Sprintf("%s@%s:", server.User, server.Host), "data", "saves", server.User, game.ID)
	localPath := game.SavePath

	sshCmd := fmt.Sprintf("ssh -p %s -i %s", server.Port, server.IdentityFile)
	var cmd *exec.Cmd


	if pull {
		state, err := getRemoteState(remotePath, localPath, sshCmd)
		if err != nil {
			return fmt.Errorf("checking remote state: %w", err)
		}

		switch state {
		case statusMissing:
			return fmt.Errorf("cannot pull: save does not exist on remote")
		case statusSyncedOrOlder:
			fmt.Println("Already up to date")
			return nil
		}

		fmt.Println("Pulling save...")
		cmd = exec.Command("rsync", "-avzhP", "-e", sshCmd, remotePath, localPath)
	} else {
		remoteState, err := getRemoteState(remotePath, localPath, sshCmd)
		if err != nil {
			return fmt.Errorf("checking remote state: %w", err)
		}

		if remoteState == statusChange {
			return  fmt.Errorf("can not push remote as remote is newer than local save")
		}

		localState, err := getRemoteState(localPath, remotePath, sshCmd)
		if err != nil {
			return fmt.Errorf("checking local state: %w", err)
		}

		if localState == statusSyncedOrOlder {
			fmt.Println("Remote save is already up to date.")
			return nil
		}

		fmt.Println("Pushing save...")
		cmd = exec.Command("rsync", "-avzhP", "-e", sshCmd, localPath, remotePath)
	}


	fmt.Printf("running this rsync command:\n%s\n", cmd.String())

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("--- rsync output start ---")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error syncing %s, from %s, to %s", game.ID, localPath, remotePath)
	}
	fmt.Println("--- rsync output end ---")

	return nil
}

func getRemoteState(src string, dest string, sshCmd string, ) (int, error) {
	cmd := exec.Command("rsync", "-naui", "-e", sshCmd, src, dest)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			code := exitError.ExitCode()
			switch code {
			case 23:
				return statusMissing, nil
			default:
				fmt.Println(string(output))
				return statusSyncedOrOlder, err
			}
		}
	}

	if len(output) > 0 {
		return statusChange, nil
	}

	return statusSyncedOrOlder, nil
}
