package config

import (
	"fmt"
	"gamesync/internal/ui"
	"os"
	"path/filepath"
	"slices"
)

func RemoveGames(current *Current, gameIdsToRemove []string) error {
	gamesRemovedQty := 0
	var gamesRemoved []GameConfig

	for _, gameIdToRemove := range gameIdsToRemove {
		found := false
		
		for i, game := range current.Config.Games {
			if gameIdToRemove == game.ID {
				current.Config.Games = slices.Delete(current.Config.Games, i, i+1)
				gamesRemoved = append(gamesRemoved, game)
				gamesRemovedQty++
				found = true
				break
			}
		}

		if ! found {
			return fmt.Errorf("game not found in config: %s", gameIdToRemove)
		} 
	}

	if gamesRemovedQty == 0 {
		return fmt.Errorf("removed no games, how did you manage to get this error?")
	}

	for _, gameID := range gameIdsToRemove {
		if err := RemoveStateFile(gameID); err != nil {
			return fmt.Errorf("removing state file: %w", err)
		}
	}

	if err := WriteGlobalConfig(*current); err != nil {
		return err
	}

	for _, game := range gamesRemoved {
		ui.Info("%s\n", game.ID)
	}

	return nil
}

func RemoveStateFile(gameID string) error {
	stateDir, err := GetStateDir()
	if err != nil {
		return fmt.Errorf("getting state dir: %v", err)
	}

	stateFile := filepath.Join(stateDir, gameID+".json")
	if err := os.RemoveAll(stateFile); err != nil {
		return fmt.Errorf("removing state file: %v", err)
	}

	return nil
}

