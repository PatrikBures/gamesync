package config

import (
	"fmt"
	"os"
	"path"
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

	poolConfig, _ := GetPoolConfig(id)
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
			Path: dirPath,
		}
		Current.Pools = append(Current.Pools, poolConfig)
		err = WriteGlobalConfig(updateConfigPath)
		if err != nil {
			return fmt.Errorf("updating config, %w", err)
		}
	}

	return nil
}

func AddSave(gameID string, poolID string, savePath string, updateConfigPath string, move bool) error {
	if _, err := GetGame(gameID); err == nil {
		return fmt.Errorf("game with id \"%s\" already exists.", gameID)
	}

	poolConfig, err := GetPoolConfig(poolID) 
	if err != nil {
		return fmt.Errorf("pool with id \"%s\" does not exist in global config.", poolID)
	}
	pool := Pool{ID: poolID}
	if err := verifyPool(pool, poolConfig.Path); err != nil {
		return err
	}

	file, err := os.Stat(savePath)
	if err != nil {
		return fmt.Errorf("stating dir %s", savePath)
	}
	if !file.IsDir() {
		return fmt.Errorf("path is not dir %s", savePath)
	}

	if move {
		gameIdPathInPool := path.Join(poolConfig.Path, gameID)
		if err := moveAndSymlink(savePath, gameIdPathInPool); err != nil {
			return err
		}
	}

	if updateConfigPath != "" {
		gameConfig := GameConfig{
			ID: gameID,
			PoolID: poolID,
			SavePath: savePath,
		}
		Current.Games = append(Current.Games, gameConfig)
		if err := WriteGlobalConfig(updateConfigPath); err != nil {
			return err
		}
	}

	return nil
}

func moveAndSymlink(source string, dest string) error {
	fmt.Printf("from: %s\nto:   %s\n", source, dest)
	if err := os.Rename(source, dest); err != nil {
		return err
	}
	if err := os.Symlink(dest, source); err != nil {
		return err
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

func verifyPool(poolExpect Pool, dirPath string) error {
	filePath := filepath.Join(dirPath, poolName)
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var pool Pool
	if err := yaml.Unmarshal(data, &pool); err != nil {
		return err
	}

	if pool.ID != poolExpect.ID {
		return fmt.Errorf("Pool ids do not match, expected: %s, got: %s", poolExpect.ID, pool.ID)
	}

	return nil
}
