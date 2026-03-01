package syncer

import (
	"encoding/json"
	"fmt"
	"gamesync/internal/config"
	"gamesync/internal/ui"
	"os"
	"os/exec"
	"path"
	"time"
)

type Snapshot struct {
	Paths		[]string	`json:"paths"`
	Hostname	string		`json:"hostname"`
	Time		time.Time	`json:"time"`
}

func initRepo(current config.Current) error {
	_, err := RunCmd(current, "restic", "cat", "config")

	if err == nil {
		return nil
	}

	output, err := RunCmd(current, "restic", "init")

	if err != nil {
		return fmt.Errorf("%s\n%s", err, output)
	}

	return nil
}

func CreateSnapshot(current config.Current, gameID string, skipUnchanged bool) error {
	saveGame := fmt.Sprintf("/data/saves/%s/%s", current.Config.Server.User, gameID)

	host, err := os.Hostname()
	if err != nil {
		return err
	}

	if err := initRepo(current); err != nil {
		return err
	}

	args := []string{"restic", "backup", saveGame, "--host", host}

	if skipUnchanged {
		args = append(args, "--skip-if-unchanged")
	}

	output, err := RunCmd(current, args...)

	if err != nil {
		return fmt.Errorf("%s\n%s", err, output)
	}

	return nil
}

func ListSnapshots(current config.Current, gameID string) ([]Snapshot, error) {
	saveGame := fmt.Sprintf("/data/saves/%s/%s", current.Config.Server.User, gameID)

	args := []string{"restic", "snapshots", "--json"}

	if gameID != "" {
		args = append(args, "--path", saveGame)
	}

	output, err := RunCmd(current, args...)
	if err != nil {
		return nil, fmt.Errorf("%s\n%s", err, output)
	}

	var snapshots []Snapshot

	if err := json.Unmarshal([]byte(output), &snapshots); err != nil {
		return nil, err
	}

	return snapshots, nil
}

func GetResticPassword(current config.Current) (string, error) {
	output, err := RunCmd(current, "get-restic-password")
	if err != nil {
		return "", fmt.Errorf("getting restic password: %w", err)
	}

	return output, err
}

func SetResticPassword(current config.Current, newPassword string) (error) {
	if newPassword == "" {
		return fmt.Errorf("password can not be empty")
	}

	
	file, err := os.CreateTemp("", "gamesync_secret_")
	if err != nil {
		return fmt.Errorf("error creating temp file: %w", err)
	}

	_, err = file.WriteString(newPassword)
	if err != nil {
		return fmt.Errorf("error writing password to temp file: %s", err)
	}

	remoteNewPassFile       := path.Join(config.RemotePasswordsDir, current.Config.Server.User, "new")
	remoteCurrentPassFile   := path.Join(config.RemotePasswordsDir, current.Config.Server.User, "current")
	remoteOldPassFile       := path.Join(config.RemotePasswordsDir, current.Config.Server.User, "old")

	sshCmd := fmt.Sprintf("ssh -p %s -i %s", current.Config.Server.Port, current.Config.Server.IdentityFile)
	remotePath := fmt.Sprintf("%s@%s:/%s", current.Config.Server.User, current.Config.Server.Host, remoteNewPassFile)

	cmd := exec.Command("rsync", "-azhPv", "-e", sshCmd, file.Name(), remotePath)

	ui.Debug("running rsync command:\n%s\n", cmd.String())
	if ui.GetLevel() == ui.LevelDebug {
		cmd.Stdout = os.Stdout
	}
	if ui.GetLevel() <= ui.LevelError {
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error copying new password to remote: %w", err)
	}

	RunCmd(current, "restic", "key", "passwd", "--new-password-file", remoteNewPassFile)

	RunCmd(current, "mv", "-f", remoteCurrentPassFile, remoteOldPassFile)
	RunCmd(current, "mv", remoteNewPassFile, remoteCurrentPassFile)

	return nil
}
