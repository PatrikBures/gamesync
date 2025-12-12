package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

const poolName = ".gamesync-pool.yml"

func InitPool(id string, dirPath string) error {
	file, err := os.Stat(dirPath)
	if err != nil {
		return fmt.Errorf("Failed stating dir %w", err)
	}
	if ! file.IsDir() {
		return fmt.Errorf("Path is not dir: %s", dirPath)
	}

	configPath := filepath.Join(dirPath, poolName)

	_, err = os.Stat(configPath)
	if os.IsExist(err) {
		return fmt.Errorf("Pool already exists at %s", configPath)
	}

	pool, _ := GetPool(id)
	if pool != nil {
		return fmt.Errorf("A pool with the id \"%s\" already exists.", id)
	}

	config := PoolConfig{
		ID: id,
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("Failed marshaling config to yaml, %w", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("Failed writing pool config file, %w", err)
	}

	return nil
}
