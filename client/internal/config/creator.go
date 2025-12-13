package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

const poolName = ".gamesync-pool.yml"

func InitPool(id string, dirPath string, updateConfigPath string) error {
	file, err := os.Stat(dirPath)
	if err != nil {
		return fmt.Errorf("stating dir %w", err)
	}
	if ! file.IsDir() {
		return fmt.Errorf("path is not dir: %s", dirPath)
	}

	poolPath := filepath.Join(dirPath, poolName)

	_, err = os.Stat(poolPath)
	if err == nil {
		return fmt.Errorf("pool already exists at %s", poolPath)
	}

	poolConfig, _ := GetPool(id)
	if poolConfig != nil {
		return fmt.Errorf("pool with id \"%s\" already exists.", id)
	}

	pool := Pool{
		ID: id,
	}

	data, err := yaml.Marshal(pool)
	if err != nil {
		return fmt.Errorf("marshaling config to yaml, %w", err)
	}

	err = os.WriteFile(poolPath, data, 0644)
	if err != nil {
		return fmt.Errorf("writing pool config file, %w", err)
	}

	if updateConfigPath != "" {
		poolConfig := PoolConfig{
			ID: id,
			Path: poolPath,
		}
		Current.Pools = append(Current.Pools, poolConfig)
		err = WriteGlobalConfig(updateConfigPath)
		if err != nil {
			return fmt.Errorf("updating config, %w", err)
		}
	}

	return nil
}

func WriteGlobalConfig(configPath string) error {
	data, err := yaml.Marshal(Current)
	if err != nil {
		return fmt.Errorf("marshaling config.")
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("writing config.")
	}

	return nil
}
