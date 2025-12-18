package syncer

import (
	"fmt"
	"gamesync/internal/config"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	remoteNewer = iota
	remoteOlder
	remoteNo
)

func SyncGame(game config.GameConfig, server config.ServerConfig, pull bool) error {
	remoteDest := filepath.Join(fmt.Sprintf("%s@%s:", server.User, server.Host), "data", "saves", server.User, game.ID)

	sshCmd := fmt.Sprintf("ssh -p %s -i %s", server.Port, server.IdentityFile)
	var cmd *exec.Cmd

	remoteState, err := isRemoteNewer(remoteDest, sshCmd, game.SavePath)
	if err != nil {
		return fmt.Errorf("getting remote, %w", err)
	}

	if pull {
		switch remoteState {
		case remoteOlder:
			return fmt.Errorf("Can not pull as remote is older than local save")
		case remoteNo:
			return fmt.Errorf("Can not pull as there is no save for %s", game.ID)
		}

		fmt.Println("Pulling save")
		cmd = exec.Command("rsync", "-avzhP", "-e", sshCmd, remoteDest, game.SavePath)
	} else {
		if remoteState == remoteNewer {
			return  fmt.Errorf("Can not push remote as remote is newer than local save")
		}
		fmt.Println("Pushing save")
		cmd = exec.Command("rsync", "-avzhP", "-e", sshCmd, game.SavePath, remoteDest)
	}


	fmt.Printf("running this rsync command:\n%s\n", cmd.String())

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("--- rsync output start ---")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error syncing %s, from %s, to %s", game.ID, game.SavePath, remoteDest)
	}
	fmt.Println("--- rsync output end ---")

	return nil
}

func isRemoteNewer(remoteDest string, sshCmd string, localDir string) (int, error) {
	cmd := exec.Command("rsync", "-naui", "-e", sshCmd, remoteDest, localDir)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			code := exitError.ExitCode()
			switch code {
			case 23:
				return remoteNo, nil
			default:
				fmt.Println(string(output))
				return remoteOlder, err
			}
		}
	}

	if len(output) > 0 {
		return remoteNewer, nil
	}

	return remoteOlder, nil
}
