package rrsync

import (
	"fmt"
	"os"
	"regexp"
	"syscall"
)

func Run(restrictedDir string, rsyncArgs []string) error {
	rDirStat, err := os.Stat(restrictedDir)
	if err != nil {
		return fmt.Errorf("restricted dir: %w", err)
	}
	if !rDirStat.IsDir() {
		return fmt.Errorf("restricted dir is not a dir: %s", restrictedDir)
	}

	syncDir, err := getSyncDir(rsyncArgs)
	if err != nil {
		return err
	}
	if err := validateSyncDir(syncDir); err != nil {
		return err
	}

	if err := syscall.Chdir(restrictedDir); err != nil {
		return fmt.Errorf("changing dir to %s: %w", restrictedDir, err)
	}
	fullArgs := append([]string{"rsync"}, rsyncArgs...)
	if err := syscall.Exec("/usr/bin/rsync", fullArgs, os.Environ()); err != nil {
		return err
	}
	return nil
}

func getSyncDir(args []string) (string, error) {
	var syncDir string
	nextIsDir := false
	for _, arg := range args {
		if syncDir != "" {
			return "", fmt.Errorf("found multiple dirs, only one is allowed")
		}
		if nextIsDir {
			syncDir = arg
		} else if arg == "." {
			nextIsDir = true
		}
	}
	if syncDir == "" {
		return "", fmt.Errorf("could not find sync dir")
	}
	return syncDir, nil
}

func validateSyncDir(dir string) error {
	matched, err := regexp.MatchString("^[A-Za-z0-9]+/?$", dir)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("dir is not valid: %s", dir)
	}
	return nil
}
