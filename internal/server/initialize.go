package server

import "os"

func CreateDirs() error {
	for _, dir := range []string{AppDir} {
		if err := os.MkdirAll(dir, 0775); err != nil {
			return err
		}
	}
	return nil
}
