package syncer

import (
	"fmt"
	"gamesync/internal/config"
)

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

func CreateSnapshot(server config.ServerConfig, verbose bool, gameID string) error {
	saveGame := fmt.Sprintf("/data/saves/%s/%s", config.Current.Server.User, gameID)

	if err := initRepo(server, verbose); err != nil {
		return err
	}

	output, err := RunCmd(server, verbose, "restic", "backup", saveGame)

	if err != nil {
		return fmt.Errorf("%s\n%s", err, output)
	}

	return nil
}
