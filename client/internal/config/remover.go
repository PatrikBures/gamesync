package config

import (
	"fmt"
	"slices"
)

func RemoveGames(gameIdsToRemove []string, configPath string) error {
	gamesRemovedQty := 0

	for _, gameIdToRemove := range gameIdsToRemove {
		found := false
		
		for i, game := range Current.Games {
			if gameIdToRemove == game.ID {
				fmt.Println("removed:", game.ID)
				Current.Games = slices.Delete(Current.Games, i, i+1)
				gamesRemovedQty++
				found = true
				break
			}
		}
		fmt.Printf("result:\n")
		for _, game := range Current.Games {
			println(game.ID)
		}

		if ! found {
			return fmt.Errorf("game not found in config: %s", gameIdToRemove)
		}
	}

	if gamesRemovedQty == 0 {
		return fmt.Errorf("removed no games")
	}

	if err := WriteGlobalConfig(configPath); err != nil {
		return err
	}

	return nil
}
