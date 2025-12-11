package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)


func Load(customPath string) error {
	var configPath string

	if customPath != "" {
		configPath = customPath
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("Could not find home directory: %w", err)
		}
		configPath = filepath.Join(home, ".config", "gamesync", "config.yml")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Failed reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &Current); err != nil {
		return fmt.Errorf("Failed parsing config file: %w", err)
	}

	return nil
}
