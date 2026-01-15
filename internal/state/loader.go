package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"gamesync/internal/config"
	"os"
	"path/filepath"
)

func GetOld(gameID string, verbose bool) (map[string]FileMeta, error) {
	stateDir, err := config.GetStateDir()
	if err != nil { return nil, err }

	gameStateFile := filepath.Join(stateDir, gameID+".json")

	if _, err := os.Stat(gameStateFile); errors.Is(err, os.ErrNotExist) {
		if verbose { fmt.Println("old local state does not exist, loading empty state:", gameStateFile) }
		return make(map[string]FileMeta), nil
	}

	stateBytes, err := os.ReadFile(gameStateFile)
	if err != nil { return nil, err }

	state := make(map[string]FileMeta)

	if err := json.Unmarshal(stateBytes, &state); err != nil {
		return nil, err
	}

	return state, nil
}
