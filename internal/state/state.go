package state

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

type FileMeta struct {
	ModTime int64 `json:"mtime"`
	Size	int64 `json:"size"`
}

func GetState(path string) (map[string]FileMeta, error) {
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
