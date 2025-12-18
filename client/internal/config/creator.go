package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

func AddSave(gameID string, savePath string, updateConfigPath string) error {
	if _, err := GetGame(gameID); err == nil {
		return fmt.Errorf("game with id \"%s\" already exists", gameID)
	}

	file, err := os.Stat(savePath)
	if err != nil {
		return fmt.Errorf("stating dir %s", savePath)
	}
	if !file.IsDir() {
		return fmt.Errorf("path is not dir %s", savePath)
	}

	if updateConfigPath != "" {
		gameConfig := GameConfig{
			ID: gameID,
			SavePath: savePath,
		}
		Current.Games = append(Current.Games, gameConfig)
		if err := WriteGlobalConfig(updateConfigPath); err != nil {
			return err
		}
	}

	return nil
}

func WriteGlobalConfig(configPath string) error {
	data, err := yaml.Marshal(Current)
	if err != nil {
		return fmt.Errorf("marshaling config")
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("writing config")
	}

	return nil
}
