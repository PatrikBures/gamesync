package syncer

import (
	"fmt"
	"gamesync/internal/config"
	"gamesync/internal/ui"
	"os"
	"os/exec"
)

func RunRsync(server config.ServerConfig, localPath string, remotePath string, toRemote bool, flags ...string) error {

	args := []string{}

	if len(flags) > 0 {
		args = append(args, flags...)
	}

	if server.SshHost == "" {
		ssh := fmt.Sprintf("ssh -p %s -i %s", server.Port, server.IdentityFile)
		args = append(args, "-e", ssh)

		remotePath = fmt.Sprintf("%s@%s:%s", server.User, server.Host, remotePath)
	} else {
		remotePath = fmt.Sprintf("%s:%s", server.SshHost, remotePath)
	}

	if toRemote {
		args = append(args, localPath, remotePath)
	} else {
		args = append(args, remotePath, localPath)
	}

	cmd := exec.Command("rsync", args...)

	ui.Debug("running rsync command:\n%s\n", cmd.String())
	if ui.GetLevel() == ui.LevelDebug {
		cmd.Stdout = os.Stdout
	}

	if ui.GetLevel() <= ui.LevelError {
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("running rsync command: %s: %w", cmd.String(), err)
	}

	return nil
}
