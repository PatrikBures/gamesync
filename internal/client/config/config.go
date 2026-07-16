package config

import (
	"errors"
	"fmt"
	"math/bits"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"go.pabu.dev/ini"
)

type CacheNames string
const (
	CacheNameUserID CacheNames = "user_id"
)

type Config struct {
	Global Global
	Server Server
}

type Global struct {
	ChunkDir string
	CacheDir string
}

type Server struct {
	UserID    int64
	TokenFile string
	Token     string
	Url       string // example: http://localhost:8080, https://example.org
}

type Sync struct {
	Name   string
	Repo   string
	Branch string
	Path   string
}

const projectName = "gamesync"

func LoadConfig(config *Config, configFile string) error {
	if configFile == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return err
		}
		configFile = filepath.Join(configDir, projectName, "config.ini")
	}
	configContent, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("reading config file '%s': %w", configFile, err)
	}
	if err := ini.Unmarshal(configContent, config); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	if err := insertHome(config); err != nil {
		return err
	}
	if err := loadToken(config); err != nil {
		return err
	}
	if err := cacheDir(config); err != nil {
		return err
	}

	if err := loadCachedItem(config.Global.CacheDir, CacheNameUserID, &config.Server.UserID); err != nil {
		return err
	}
	return nil
}

func insertHome(config *Config) error {
	home := ""
	for _, p := range []*string{
		&config.Global.ChunkDir,
	} {
		s := *p
		if !strings.HasPrefix(s, "~/") {
			continue
		}
		if home == "" {
			var err error
			home, err = os.UserHomeDir()
			if err != nil {
				return err
			}
		}
		newPath := filepath.Join(home, s[2:])
		*p = newPath
	}
	return nil
}

func loadToken(config *Config) error {
	if config.Server.TokenFile != "" {
		token, err := os.ReadFile(config.Server.TokenFile)
		if err != nil {
			return fmt.Errorf("loading token file: %w", err)
		}
		// trims new line if it ends with it
		if token[len(token)-1] == '\n' {
			token = token[:len(token)-1]
		}
		tokenString := string(token)
		if len(tokenString) != 44 {
			return fmt.Errorf("expected token length to be %d, got %d", 44, len(tokenString))
		}
		config.Server.Token = string(token)
	}
	if config.Server.Token == "" {
		return fmt.Errorf("token can not be empty")
	}
	return nil
}

func loadCachedItem[T int | int32 | int64 | string](cacheDir string, name CacheNames, target *T) (error) {
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

func cacheDir(config *Config) error {
	if config.Global.CacheDir == "" {
		cd, err := os.UserCacheDir()
		if err != nil {
			return err
		}
		config.Global.CacheDir = filepath.Join(cd, projectName)
	}

	s, err := os.Stat(config.Global.CacheDir)
	if err != nil {
		return err
	}
	if !s.IsDir() {
		return fmt.Errorf("cache dir is not a dir: %s", config.Global.CacheDir)
	}

	config.Global.ChunkDir = filepath.Join(config.Global.CacheDir, "chunks")

	if err := createDirs([]string{
		config.Global.ChunkDir,
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
