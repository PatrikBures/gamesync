package util

import "os"


func CreateDirs(dirs []string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0775); err != nil {
			return err
		}
	}
	return nil
}

