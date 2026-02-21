package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"gamesync/internal/ui"
	"io/fs"
	"os"
	"path/filepath"
)

type FileMeta struct {
	ModTime int64 `json:"mtime"`
	Size	int64 `json:"size"`
}

func Get(root string) (map[string]FileMeta, error) {
	results := make(map[string]FileMeta)

	if _, err := os.Stat(root); errors.Is(err, fs.ErrNotExist) {
		return results, nil
	}

	err := filepath.WalkDir(root, 
		func(path string, d fs.DirEntry, err error) error {
			if err != nil { return err }
			if d.IsDir() { return nil }

			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}

			info, _ := d.Info()
			results[relPath] = FileMeta{
				ModTime: info.ModTime().Unix(),
				Size: info.Size(),
			}

			return nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("walking dir: %w", err)
	}

	ui.Debug("got state of dir: %s\n", root)

	return results, nil
}

func Write(state map[string]FileMeta, path string) error {
	stateJson, err := json.MarshalIndent(state, "", "    ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, stateJson, 0664)
	err != nil {
		errRemove := os.RemoveAll(path)
		if errRemove != nil {
			return fmt.Errorf("%w, %w", err, errRemove)
		}
		return err
	}

	return nil
}

type SyncState int
const (
	SyncStateConflict SyncState = iota
	SyncStatePull
	SyncStatePush
	SyncStateUnchanged
	SyncStateError
)

func Compare(localState map[string]FileMeta, oldLocalState map[string]FileMeta, remoteState map[string]FileMeta, loose bool) (SyncState, error) {
	if len(localState) == 0 && len(remoteState) >  0 { return SyncStatePull, nil }
	if len(localState) >  0 && len(remoteState) == 0 { return SyncStatePush, nil }
	if len(localState) == 0 && len(remoteState) == 0 { return SyncStateError, fmt.Errorf("local and remote empty")}

	localChange := true
	remoteChange := true

	ui.Verbose("Comparing local to old local...\n")
	if stateEqual(localState, oldLocalState, loose) {
		localChange = false
	}
	ui.Verbose("Comparing remote to old local...\n")
	if stateEqual(remoteState, oldLocalState, loose) {
		remoteChange = false
	}

	if localChange  && remoteChange  { return SyncStateConflict, nil }
	if localChange  && !remoteChange { return SyncStatePush, nil }
	if !localChange && remoteChange  { return SyncStatePull, nil }
	if !localChange && !remoteChange { return SyncStateUnchanged, nil }

	return SyncStateError, fmt.Errorf("there is something very wrong, this error should not be possible")
}

func stateEqual(state1 map[string]FileMeta, state2 map[string]FileMeta, loose bool) bool {
	if len(state1) != len(state2) {
		ui.Verbose("len of states not same, state1: %d, state2: %d\n", len(state1), len(state2))
		return false
	}
	for path, meta1 := range state1 {
		meta2, ok := state2[path]
		if !ok {
			ui.Verbose("does not exist: %s\n", path)
			return false
		}

		if meta1.Size != meta2.Size {
			ui.Verbose("size: %s\n", path)
			return false
		}

		if loose {
			diff := meta1.ModTime - meta2.ModTime
			if diff < -2 || diff > 2 {
				ui.Verbose("loose modtime: %s\n", path)
				return false
			}
		} else {
			if meta1.ModTime != meta2.ModTime {
				ui.Verbose("modtime: %s\n", path)
				return false
			}
		}
		ui.Debug("same: %s\n", path)
	}

	ui.Verbose("state 1 and 2 were same\n")
	return true
}
