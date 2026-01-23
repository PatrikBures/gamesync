package config

import (
	"gamesync/internal/testutils"
	"os"
	"path/filepath"
	"testing"
)

func TestGetStateDir(t *testing.T) {
	_, tmpDir := testutils.SetupTest(t)

	if err := os.Setenv("XDG_STATE_HOME", tmpDir); err != nil {
		t.Fatal(err)
	}
	
	wantStateDir := filepath.Join(tmpDir, "gamesync")

	gotStateDir, err := GetStateDir()
	if err != nil {
		t.Errorf("GetStateDir returned error: %v", err)
	}

	if wantStateDir != gotStateDir {
		t.Errorf("state dir do not match, want: %s, got: %s", wantStateDir, gotStateDir)
	}

	info, err := os.Stat(wantStateDir)
	if err != nil {
		t.Errorf("error getting state of stateDir: %v", err)
	}

	if ! info.IsDir() {
		t.Errorf("state dir is not a dir")
	}
}

func removeDir(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		panic(err)
	}
}
