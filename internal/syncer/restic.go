package syncer

import (
	"encoding/json"
	"fmt"
	"gamesync/internal/config"
	"os"
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
