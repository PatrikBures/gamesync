package config

import (
	"fmt"
	"os"
	"path/filepath"
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
	ProfilesFile string
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

// Loads config into c from configFile
func Load(c *Config, configFile string) error {
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
	if err := ini.Unmarshal(configContent, c); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	if err := c.insertHome(); err != nil {
		return err
	}
	if err := c.loadToken(); err != nil {
		return err
	}
	if err := c.cacheDir(); err != nil {
		return err
	}

	if err := loadCachedItem(c.Global.CacheDir, CacheNameUserID, &c.Server.UserID); err != nil {
		return err
	}

	if c.Global.ProfilesFile == "" {
		configDir := filepath.Dir(configFile)
		c.Global.ProfilesFile = filepath.Join(configDir, "profiles.json")
	}

	return nil
}

// Inserts current users home for paths that start with ~/
func (c *Config) insertHome() error {
	home := ""
	for _, p := range []*string{
		&c.Global.ChunkDir,
		&c.Global.CacheDir,
		&c.Global.ProfilesFile,
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

// loads token from from TokenFile
func (c *Config) loadToken() error {
	if c.Server.TokenFile != "" {
		token, err := os.ReadFile(c.Server.TokenFile)
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
		c.Server.Token = string(token)
	}
	if c.Server.Token == "" {
		return fmt.Errorf("token can not be empty")
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
