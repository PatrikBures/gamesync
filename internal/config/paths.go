package config

import (
	"os"
	"path/filepath"
)

func GetStateDir() (string, error) {
	var stateDir string

	if env := os.Getenv("XDG_STATE_HOME"); env != "" {
		stateDir = filepath.Join(env, "gamesync")
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		stateDir = filepath.Join(home, ".local", "state", "gamesync")
	}

	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return "", err
	}

	return stateDir, nil
}
