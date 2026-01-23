package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetStateDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "state_dir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	stateEnvVar := "XDG_STATE_HOME"

	os.Setenv(stateEnvVar, tmpDir)
	
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

