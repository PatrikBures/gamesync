package config

import (
	"errors"
	"fmt"
	"math/bits"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Loads content from a cached file if it exists into target, cache file does not exist it simply returns.
//
// I wanted this to be in the Config struct, but methods do not support type parameters.
// It could be modified to accept 'any' instead of 'T', but that would
// not be compile-time safe, Which i prefer.
func loadCachedItem[T int | int32 | int64 | string](cacheDir string, name CacheNames, target *T) error {
	p := filepath.Join(cacheDir, string(name))
	
	var valueString string
	if content, err := os.ReadFile(p); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("did not find:", p)
			return nil
		}
		return err
	} else {
		valueString = string(content)
		valueString = strings.TrimSpace(valueString)
	}

	errMsg := fmt.Errorf("loading cache '%s'", p)
	switch t := any(target).(type) {
	case *int:
		i, err := strconv.ParseInt(valueString, 10, bits.UintSize)
		if err != nil { return errors.Join(errMsg, err) }
		*t = int(i)
	case *int32:
		i, err := strconv.ParseInt(valueString, 10, 32)
		if err != nil { return errors.Join(errMsg, err) }
		*t = int32(i)
	case *int64:
		i, err := strconv.ParseInt(valueString, 10, 64)
		if err != nil { return errors.Join(errMsg, err) }
		*t = int64(i)
	case *string:
		*t = valueString
	default:
		panic(fmt.Sprintf("undefined type when loading cache: %T", target))
	}
	return nil
}

// Sets the value of a cache file. 
func SetCacheItem[T int | int32 | int64 | string](cacheDir string, name CacheNames, value T) (err error) {
	p := filepath.Join(cacheDir, string(name))
	
	f, err := os.OpenFile(p, os.O_CREATE | os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	var valueString string

	switch t := any(value).(type) {
	case int, int32, int64:
		i := t.(int64)
		valueString = strconv.FormatInt(i, 10)
	case string:
		valueString = t
	}
	if _, err = f.WriteString(valueString); err != nil {
		return err
	}

	fmt.Printf("set cache '%s' to '%s'\n", string(name), valueString)

	return nil
}

func (c *Config) cacheDir() error {
	if c.Global.CacheDir == "" {
		cd, err := os.UserCacheDir()
		if err != nil {
			return err
		}
		c.Global.CacheDir = filepath.Join(cd, projectName)
	}

	s, err := os.Stat(c.Global.CacheDir)
	if err != nil {
		return err
	}
	if !s.IsDir() {
		return fmt.Errorf("cache dir is not a dir: %s", c.Global.CacheDir)
	}

	c.Global.ChunkDir = filepath.Join(c.Global.CacheDir, "chunks")

	if err := createDirs([]string{
		c.Global.ChunkDir,
	}); err != nil {
		return err
	}
	return nil
}
