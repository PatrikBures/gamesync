package syncer

import (
	"encoding/json"
	"fmt"
	"gamesync/internal/config"
	"os"
	"path"
	"time"
)

type Snapshot struct {
	Paths		[]string	`json:"paths"`
	Hostname	string		`json:"hostname"`
	Time		time.Time	`json:"time"`
}

func initRepo(current config.Current) error {
	_, err := RunCmd(current.Config.Server, false, "restic", "cat", "config")

	if err == nil {
		return nil
	}

	output, err := RunCmd(current.Config.Server, false, "restic", "init")

	if err != nil {
		return fmt.Errorf("%s\n%s", err, output)
	}

	return nil
}

func CreateSnapshot(current config.Current, gameID string, skipUnchanged bool) error {
	if err := SameApiVersion(current.Config.Server); err != nil {
		return fmt.Errorf("could not create snapshot: %w", err)
	}

	saveGame := fmt.Sprintf("%s/%s/%s", config.RemoteSavesDir, current.Config.Server.User, gameID)

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

	output, err := RunCmd(current.Config.Server, false, args...)

	if err != nil {
		return fmt.Errorf("%s\n%s", err, output)
	}

	return nil
}

func ListSnapshots(current config.Current, gameID string) ([]Snapshot, error) {
	saveGame := fmt.Sprintf("%s/%s/%s", config.RemoteSavesDir, current.Config.Server.User, gameID)

	args := []string{"restic", "snapshots", "--json"}

	if gameID != "" {
		args = append(args, "--path", saveGame)
	}

	output, err := RunCmd(current.Config.Server, true, args...)
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
	output, err := RunCmd(current.Config.Server, true, "get-restic-password")
	if err != nil {
		return "", fmt.Errorf("getting restic password: %w", err)
	}

	return output, err
}

func SetResticPassword(current config.Current, newPassword string) (error) {
	if err := SameApiVersion(current.Config.Server); err != nil {
		return err
	}
	if newPassword == "" {
		return fmt.Errorf("password can not be empty")
	}

	tempPassFile, err := os.CreateTemp("", "gamesync_secret_")
	if err != nil {
		return fmt.Errorf("error creating temp file: %w", err)
	}

	_, err = tempPassFile.WriteString(newPassword)
	if err != nil {
		return fmt.Errorf("error writing password to temp file: %s", err)
	}

	remoteNewPassFile       := path.Join(config.RemotePasswordsDir, current.Config.Server.User, "new")
	remoteCurrentPassFile   := path.Join(config.RemotePasswordsDir, current.Config.Server.User, "current")
	remoteOldPassFile       := path.Join(config.RemotePasswordsDir, current.Config.Server.User, "old")

	if err := RunRsync(current.Config.Server, tempPassFile.Name(), remoteNewPassFile, true, "-azhPv"); err != nil {
		return fmt.Errorf("error copying new password to remote: %w", err)
	}

	if _, err := RunCmd(current.Config.Server, false, "restic", "key", "passwd", "--new-password-file", remoteNewPassFile); err != nil {
		return fmt.Errorf("changing restic password (password did not change): %w", err)
	}
	if _, err := RunCmd(current.Config.Server, false, "mv", "-f", remoteCurrentPassFile, remoteOldPassFile); err != nil {
		return fmt.Errorf("moving old password from %s, to %s: %w", remoteCurrentPassFile, remoteOldPassFile, err)
	}
	if _, err := RunCmd(current.Config.Server, false, "mv", remoteNewPassFile, remoteCurrentPassFile); err != nil {
		return fmt.Errorf("moving new password from %s, to %s: %w", remoteNewPassFile, remoteCurrentPassFile, err)
	}

	if err := os.RemoveAll(tempPassFile.Name()); err != nil {
		return fmt.Errorf("removing temp password file: %s: %w", tempPassFile.Name(), err)
	}

	return nil
}
