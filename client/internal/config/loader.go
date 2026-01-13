package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)


func Load(configPath string) error {

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &Current); err != nil {
		return fmt.Errorf("parsing config file: %w", err)
	}

	return nil
}

func GetGame(gameID string) (*GameConfig, int, error) {
	for i, g := range Current.Games {
		if g.ID == gameID {
			return &g,i , nil
		}
	}
	return nil, -1,  fmt.Errorf("game id %s could not found in config", gameID)
}
