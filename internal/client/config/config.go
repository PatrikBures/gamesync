package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.pabu.dev/gamesync/internal/client/util"
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
	ChunkDir     string
	CacheDir     string
	StateDir     string
	ProfilesFile string
}

type Server struct {
	UserID    int64
	TokenFile string
	Token     string
	Url       string // example: http://localhost:8080, https://example.org
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

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	if err := c.insertHome(home); err != nil {
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

	if c.Global.StateDir == "" {
		stateDir := os.Getenv("XDG_STATE_HOME")
		if stateDir == "" {
			stateDir = filepath.Join(home, ".local", "state", projectName)
		}
		c.Global.StateDir = stateDir
	}

	if err := util.CreateDirs([]string{
		c.Global.StateDir,
		c.ProfileStateDir(),
	}); err != nil {
		return err
	}

	return nil
}

// replaces ~/ with home for paths in config
func (c *Config) insertHome(home string) error {
	for _, p := range []*string{
		&c.Global.ChunkDir,
		&c.Global.CacheDir,
		&c.Global.StateDir,
		&c.Global.ProfilesFile,
	} {
		s := *p
		if !strings.HasPrefix(s, "~/") {
			continue
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


func (c *Config) ProfileStateDir() string {
	return filepath.Join(c.Global.StateDir, "profiles")
}
