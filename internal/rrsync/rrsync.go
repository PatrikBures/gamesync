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

	dotIdx := -1
	for i, o := range rsyncArgs {
		if o == "." {
			dotIdx = i
		}
	}
	if dotIdx == -1 {
		return fmt.Errorf("could not find required '.'")
	}
	// only allow one item after dot, which is the syncDir
	if len(rsyncArgs) -2 != dotIdx {
		return fmt.Errorf("more that one dirs after '.'")
	}

	syncDir := rsyncArgs[dotIdx+1]
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
