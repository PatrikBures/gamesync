package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

func AddSave(current *Current, gameID string, savePath string, update bool) error {
	_, gameIdx, getGameErr := GetGame(*current, gameID)

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
		current.Config.Games[gameIdx].SavePath = savePath
	} else {
		gameConfig := GameConfig{
			ID: gameID,
			SavePath: savePath,
		}

		current.Config.Games = append(current.Config.Games, gameConfig)
	}
	if err := WriteGlobalConfig(*current); err != nil {
		return err
	}

	return nil
}

func WriteGlobalConfig(current Current) error {
	data, err := yaml.Marshal(current.Config)
	if err != nil {
		return fmt.Errorf("marshaling config: %v", err)
	}

	err = os.WriteFile(current.ConfigPath, data, 0644)
	if err != nil {
		return fmt.Errorf("writing config: %v", err)
	}

	return nil
}
