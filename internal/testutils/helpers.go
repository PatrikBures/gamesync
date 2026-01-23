package testutils

import (
	"os"
	"testing"
)

func SetupTest(t *testing.T) (string, string) {
	t.Helper()

	configDir, err := os.MkdirTemp("", "test_config")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(configDir); err != nil {
			t.Errorf("Failed to cleanup config dir: %v", err)
		}
	})


	stateDir, err := os.MkdirTemp("", "test_state")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(stateDir); err != nil {
			t.Errorf("Failed to cleanup state dir: %v", err)
		}
	})

	t.Setenv("XDG_CONFIG_HOME", configDir)
	t.Setenv("XDG_STATE_HOME", stateDir)

	return configDir, stateDir

}
