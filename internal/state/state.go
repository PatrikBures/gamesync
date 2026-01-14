package state

import (
	"encoding/json"
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
		os.Remove(path)
		return err
	}

	return nil
}
