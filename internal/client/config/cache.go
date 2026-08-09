package config

import (
	"os"
	"path/filepath"

	"go.pabu.dev/gamesync/internal/client/fileval"
	"go.pabu.dev/gamesync/internal/client/util"
)

// Loads content from a cached file if it exists into target, otherwise it returns nil.
//
// I wanted this to be in the Config struct, but methods do not support type parameters.
// It could be modified to accept 'any' instead of 'T', but that would
// not be compile-time safe, Which i prefer.
func loadCachedItem[T int | int32 | int64 | string](cacheDir string, name CacheNames, target *T) error {
	p := filepath.Join(cacheDir, string(name))

	err := fileval.Read(p, target)
	if os.IsNotExist(err) {
		return nil
	}

	return err
}

// Sets the value of a cache file.
func SetCacheItem[T int | int32 | int64 | string](cacheDir string, name CacheNames, value T) (err error) {
	p := filepath.Join(cacheDir, string(name))

	return fileval.Write(p, value)
}

func (c *Config) cacheDir() error {
	if c.Global.CacheDir == "" {
		cd, err := os.UserCacheDir()
		if err != nil {
			return err
		}
		c.Global.CacheDir = filepath.Join(cd, projectName)
	}

	c.Global.ChunkDir = filepath.Join(c.Global.CacheDir, "chunks")

	if err := util.CreateDirs([]string{
		c.Global.ChunkDir,
	}); err != nil {
		return err
	}
	return nil
}
