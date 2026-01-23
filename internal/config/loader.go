package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)


func Load(current *Current) error {

	data, err := os.ReadFile(current.ConfigPath)
	if err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &current.Config); err != nil {
		return fmt.Errorf("parsing config file: %w", err)
	}

	return nil
}

func GetGame(current Current, gameID string) (GameConfig, int, error) {
	for i, g := range current.Config.Games {
		if g.ID == gameID {
			return g,i , nil
		}
	}
	return GameConfig{}, -1,  fmt.Errorf("game id %s could not found in config", gameID)
}
