package state

import (
	"encoding/json"
	"errors"
	"fmt"
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

const (
	CompareStateConflict = iota
	CompareStatePull
	CompareStatePush
	CompareStateUnchanged
	CompareStateError
)

func Compare(localState map[string]FileMeta, oldLocalState map[string]FileMeta, remoteState map[string]FileMeta, loose bool, verbose bool) (int, error) {
	localChange := true
	remoteChange := true

	if verbose { fmt.Println("Comparing local to old local...") }
	if stateEqual(localState, oldLocalState, loose, verbose) {
		localChange = false
	}
	if verbose { fmt.Println("Comparing remote to old local...") }
	if stateEqual(remoteState, oldLocalState, loose, verbose) {
		remoteChange = false
	}

	if localChange  && remoteChange  { return CompareStateConflict, nil }
	if localChange  && !remoteChange { return CompareStatePush, nil }
	if !localChange && remoteChange  { return CompareStatePull, nil }
	if !localChange && !remoteChange { return CompareStateUnchanged, nil }

	return CompareStateError, fmt.Errorf("There is something very wrong, this error should not be possible")
}

func stateEqual(state1 map[string]FileMeta, state2 map[string]FileMeta, loose bool, verbose bool) bool {
	for path, meta1 := range state1 {
		meta2, ok := state2[path]
		if !ok {
			if verbose { fmt.Println("does not exist:", path) }
			return false
		}

		if meta1.Size != meta2.Size {
			if verbose { fmt.Println("size:", path)}
			return false
		}

		if loose {
			diff := meta1.ModTime - meta2.ModTime
			if diff < -2 || diff > 2 {
				if verbose { fmt.Println("loose modtime:", path) }
				return false
			}
		} else {
			if meta1.ModTime != meta2.ModTime {
				if verbose { fmt.Println("modtime:", path) }
				return false
			}
		}
		if verbose { fmt.Println("same:", path) }
	}

	if verbose { fmt.Println("state 1 and 2 were same") }
	return true
}
