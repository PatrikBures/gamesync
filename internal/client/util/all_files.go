package util

import (
	"io/fs"
	"path/filepath"
)


func FilesInDirRecursive(dir string) (map[string]struct{}, error) {
	var files map[string]struct{}
	if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		files[path] = struct{}{}

		return nil
	}); err != nil {
		return nil, err
	}
	return files, nil
}
