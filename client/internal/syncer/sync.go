package syncer

import (
	"fmt"
	"gamesync/internal/config"
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
		preCmd := exec.Command("ssh", "-p", server.Port, "-i", server.IdentityFile, 
			fmt.Sprintf("%s@%s", server.User, server.Host),
			"mkdir -p /"+remoteDir)

		if verbose {
			fmt.Printf("Ran preCmd:\n%s\n\n", preCmd.String())
		}

		if err := preCmd.Run(); err != nil {
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
