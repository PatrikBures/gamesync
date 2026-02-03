package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gamesync/internal/config"
	"gamesync/internal/ui"
)

func GetOld(current config.Current, gameID string) (map[string]FileMeta, error) {
	stateDir, err := config.GetStateDir()
	if err != nil { return nil, err }

	gameStateFile := filepath.Join(stateDir, gameID+".json")

	if _, err := os.Stat(gameStateFile); errors.Is(err, os.ErrNotExist) {
		ui.Verbose("Old local state does not exist, loading empty state: %s\n", gameStateFile)
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
