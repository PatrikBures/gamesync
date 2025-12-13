package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)


func Load(configPath string) error {

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Failed reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &Current); err != nil {
		return fmt.Errorf("Failed parsing config file: %w", err)
	}

	return nil
}

func GetGame(gameID string) (*GameConfig, error) {
	for _, g := range Current.Games {
		if g.ID == gameID {
			return &g, nil
		}
	}
	return nil, fmt.Errorf("Game id %s not found in config.", gameID)
}

func GetPool(poolID string) (*PoolConfig, error) {
	for _, p := range Current.Pools {
		if p.ID == poolID {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("Pool id %s could not be found in config.", poolID)
}
