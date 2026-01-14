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

func GetState(path string) ([]FileMeta, error) {
	var files []FileMeta
	err := filepath.Walk(path, 
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			var file FileMeta

			file.ModTime = info.ModTime().Unix()
			file.Size = info.Size()

			files = append(files, file)

			return nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("walking dir: %w", err)
	}

	return files, nil
}
