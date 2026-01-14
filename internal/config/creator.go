package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

func AddSave(gameID string, savePath string, configPath string, update bool) error {
	_, gameIdx, getGameErr := GetGame(gameID)

	if !update && getGameErr == nil {
		return fmt.Errorf("game with id \"%s\" already exists", gameID)
	}

	file, err := os.Stat(savePath)
	if err != nil {
		return fmt.Errorf("stating dir: %v", err)
	}
	if !file.IsDir() {
		return fmt.Errorf("path is not dir: %s", savePath)
	}

	if update && getGameErr == nil {
		Current.Games[gameIdx].SavePath = savePath
	} else {
		gameConfig := GameConfig{
			ID: gameID,
			SavePath: savePath,
		}

		Current.Games = append(Current.Games, gameConfig)
	}
	if err := WriteGlobalConfig(configPath); err != nil {
		return err
	}

	return nil
}

func WriteGlobalConfig(configPath string) error {
	data, err := yaml.Marshal(Current)
	if err != nil {
		return fmt.Errorf("marshaling config: %v", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("writing config: %v", err)
	}

	return nil
}
