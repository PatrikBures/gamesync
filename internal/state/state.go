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

func Get(path string) (map[string]FileMeta, error) {
	files := make(map[string]FileMeta)

	err := filepath.Walk(path, 
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			var meta FileMeta

			meta.ModTime = info.ModTime().Unix()
			meta.Size = info.Size()

			files[path] = meta

			return nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("walking dir: %w", err)
	}

	return files, nil
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
