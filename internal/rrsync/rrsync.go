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

	syncDir := getSyncDir(rsyncArgs)
	if syncDir == "" {
		return fmt.Errorf("could not get rsync dir")
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

func getSyncDir(args []string) string {
	var syncDir string
	nextIsDir := false
	for _, arg := range args {
		if nextIsDir {
			syncDir = arg
		}
		if arg == "." {
			nextIsDir = true
		}
	}
	return syncDir
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
