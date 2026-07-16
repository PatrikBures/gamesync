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

type Config struct {
	Global Global
	Server Server
}

type Global struct {
	ChunkDir string
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
	cache, err := cacheDir(config)
	if err != nil {
		return err
	}

	if err := loadCachedItem(cache, "user_id", &config.Server.UserID); err != nil {
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

func loadCachedItem[T int | int32 | int64 | string](cacheDir string, name string, target *T) (error) {
	p := filepath.Join(cacheDir, name)
	
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

func cacheDir(config *Config) (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	projectCache := filepath.Join(cacheDir, projectName)
	config.Global.ChunkDir = filepath.Join(projectCache, "chunks")

	if err := createDirs([]string{
		config.Global.ChunkDir,
	}); err != nil {
		return "", err
	}
	return projectCache, nil
}

func createDirs(dirs []string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0775); err != nil {
			return err
		}
	}
	return nil
}
