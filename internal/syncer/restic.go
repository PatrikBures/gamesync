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

func initRepo(server config.ServerConfig, verbose bool) error {
	_, err := RunCmd(server, verbose, "restic", "cat", "config")

	if err == nil {
		return nil
	}

	output, err := RunCmd(server, verbose, "restic", "init")

	if err != nil {
		return fmt.Errorf("%s\n%s", err, output)
	}

	return nil
}

func CreateSnapshot(server config.ServerConfig, verbose bool, gameID string, skipUnchanged bool) error {
	saveGame := fmt.Sprintf("/data/saves/%s/%s", config.Current.Server.User, gameID)

	host, err := os.Hostname()
	if err != nil {
		return err
	}

	if err := initRepo(server, verbose); err != nil {
		return err
	}

	args := []string{"restic", "backup", saveGame, "--host", host}

	if skipUnchanged {
		args = append(args, "--skip-if-unchanged")
	}

	output, err := RunCmd(server, verbose, args...)

	if err != nil {
		return fmt.Errorf("%s\n%s", err, output)
	}

	return nil
}

func ListSnapshots(server config.ServerConfig, verbose bool, gameID string) ([]Snapshot, error) {
	saveGame := fmt.Sprintf("/data/saves/%s/%s", config.Current.Server.User, gameID)

	args := []string{"restic", "snapshots", "--json"}

	if gameID != "" {
		args = append(args, "--path", saveGame)
	}

	output, err := RunCmd(server, verbose, args...)
	if err != nil {
		return nil, fmt.Errorf("%s\n%s", err, output)
	}

	var snapshots []Snapshot

	if err := json.Unmarshal([]byte(output), &snapshots); err != nil {
		return nil, err
	}

	return snapshots, nil
}
