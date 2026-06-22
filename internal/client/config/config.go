package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.pabu.dev/PatrikBures/ini"
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
	if err := cache(config); err != nil {
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
		return fmt.Errorf("Token can not be empty")
	}
	return nil
}

func cache(config *Config) error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	config.Global.ChunkDir = filepath.Join(cacheDir, projectName)

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
