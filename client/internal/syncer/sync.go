package syncer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"gamesync/internal/config"
)

func SyncGame(game config.GameConfig, server config.ServerConfig, pull bool) error {
	remoteDest := filepath.Join(fmt.Sprintf("%s@%s:", server.User, server.Host), "data", "saves", server.User, game.ID)

	sshCmd := fmt.Sprintf("ssh -p %s -i %s", server.Port, server.IdentityFile)
	var cmd *exec.Cmd

	isRemoteNewer, err := isRemoteNewer(remoteDest, sshCmd, game.SavePath)
	if err != nil {
		return fmt.Errorf("getting remote, %w", err)
	}

	if pull {
		if ! isRemoteNewer {
			return fmt.Errorf("Can not pull as remote is older than local save")
		}
		fmt.Println("Pulling save")
		cmd = exec.Command("rsync", "-avzhP", "-e", sshCmd, remoteDest, game.SavePath)
	} else {
		if isRemoteNewer {
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

func isRemoteNewer(remoteDest string, sshCmd string, localDir string) (bool, error) {
	cmd := exec.Command("rsync", "-naui", "-e", sshCmd, remoteDest, localDir)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return false, err
	}

	if len(output) > 0 {
		return true, nil
	}

	return false, nil
}
