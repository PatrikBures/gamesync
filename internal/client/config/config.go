package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	ChunkDir string
	Token    string
	Server   string // example: http://localhost:8080, https://example.org
	UserID   int64
}

const projectName = "gamesync"

func LoadConfig(config *Config) error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	config.ChunkDir = filepath.Join(cacheDir, projectName)

	if err := createDirs([]string{
		config.ChunkDir,
	}); err != nil {
		return err
	}
	return nil
}

func createDirs(dirs []string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0775); err != nil {
			return err
		}
	}
	return nil
}
