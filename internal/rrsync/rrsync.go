package rrsync

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"syscall"
)

var allowedOptsLong = []string {
	"--server",
	"--sender",
	"--delete",
}
const allowedOptsShortBeforeDot = "logDtprze"

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

	if err := validOpts(rsyncArgs[:dotIdx]); err != nil {
		return err
	}
	syncDir := rsyncArgs[dotIdx+1]
	if err := validSyncDir(syncDir); err != nil {
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

func validSyncDir(dir string) error {
	matched, err := regexp.MatchString("^[A-Za-z0-9]+/?$", dir)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("dir is not valid: %s", dir)
	}
	return nil
}

func validOpts(opts []string) error {
	hasD := false
	hasDD := false
	for _, opt := range opts {
		if strings.HasPrefix(opt, "--") {
			if !slices.Contains(allowedOptsLong, opt) {
				return fmt.Errorf("used opt which is not allowed: %s", opt)
			}
			hasDD = true
		} else if strings.HasPrefix(opt, "-") {
			for _, o := range opt[1:] {
				if o == '.' {
					break
				}
				if !strings.ContainsRune(allowedOptsShortBeforeDot, o) {
					return fmt.Errorf("used a short opt which is not allowed: %c", o)
				}
			}
			hasD = true
		}
	}
	if hasD == false && hasDD == false {
		return fmt.Errorf("does not have - or --")
	}
	return nil
}
