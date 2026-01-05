package syncer

import (
	"fmt"
	"gamesync/internal/config"
)

func initRepo(server config.ServerConfig, verbose bool, passwordFile string, repo string) error {
	_, err := RunCmd(server, verbose, "restic", 
		"--password-file", passwordFile, 
		"--repo", repo,
		"cat", "config")

	if err != nil {
		_, err = RunCmd(server, verbose, "restic", 
			"--password-file", passwordFile, 
			"--repo", repo,
			"init")
	}

	return nil
}

func BackupGame(server config.ServerConfig, verbose bool, gameID string) error {
	passwordFile := fmt.Sprintf("/data/saves/%s/.restic_password", config.Current.Server.User)
	repo := fmt.Sprintf("/data/repos/%s", config.Current.Server.User)
	saveGame := fmt.Sprintf("/data/saves/%s/%s", config.Current.Server.User, gameID)

	if err := initRepo(server, verbose, passwordFile, repo); err != nil {
		return err
	}

	_, err := RunCmd(server, verbose, "restic", 
		"--password-file", passwordFile, 
		"--repo", repo,
		"backup", saveGame)

	if err != nil {
		return nil
	}

	return nil
}
