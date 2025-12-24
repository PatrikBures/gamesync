package config

import (
	"fmt"
	"slices"
)

func RemoveGames(gameIdsToRemove []string, configPath string) error {
	gamesRemovedQty := 0
	var gamesRemoved []*GameConfig

	for _, gameIdToRemove := range gameIdsToRemove {
		found := false
		
		for i, game := range Current.Games {
			if gameIdToRemove == game.ID {
				Current.Games = slices.Delete(Current.Games, i, i+1)
				gamesRemoved = append(gamesRemoved, &game)
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

	if err := WriteGlobalConfig(configPath); err != nil {
		return err
	}

	for _, game := range gamesRemoved {
		fmt.Println(&game.ID)
	}

	return nil
}
