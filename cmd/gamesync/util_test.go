package main

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"gamesync/internal/config"
	"gamesync/internal/syncer"
	"gamesync/internal/ui"
)



func setupTest(t *testing.T) {
	t.Helper()

	ui.SetLevel(ui.LevelError)

	current.Config.Server.User = "test-user"
	current.Config.Server.Host = "server"
	current.Config.Server.Port = "22"
	current.Config.Server.IdentityFile = "/tmp/ssh/key"

	setupTestGames(t)
	setupTestStateDir(t)
	setupTestConfigFile(t)
	populateConfigFile(t)
}

func setupTestGames(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gamesync_test_saves_")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to cleanup save dir: %v", err)
		}
	})

	for i := range 3 {
		saveName := "game_"+strconv.Itoa(i)
		saveDir := filepath.Join(tmpDir, saveName)

		game := config.GameConfig{
			ID: saveName,
			SavePath: saveDir,
		}

		current.Config.Games = append(current.Config.Games, game)

		createTestGame(t, saveDir)
	}
}


func createTestGame(t *testing.T, path string) {
	if err := os.Mkdir(path, 0770); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(path, "save.txt")

	f, err := os.OpenFile(testFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = f.Close() }()

	if _, err := f.WriteString("test"); err != nil {
		t.Fatal(err)
	}
}

func setupTestStateDir(t *testing.T) {
	stateDir, err := os.MkdirTemp("", "gamesync_test_state_")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(stateDir); err != nil {
			t.Errorf("Failed to cleanup state dir: %v", err)
		}
	})

	t.Setenv("GAMESYNC_STATE", stateDir)
}

func setupTestConfigFile(t *testing.T) {
	f, err := os.CreateTemp("", "gamesync_test_config_")
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = f.Close() }()

	t.Cleanup(func() {
		if err := os.RemoveAll(f.Name()); err != nil {
			t.Errorf("Failed to clean up config file: %v", err)
		}
	})

	current.ConfigPath = f.Name()
	t.Setenv("GAMESYNC_CONFIG", f.Name())
}

func populateConfigFile(t *testing.T) {
	if err := config.WriteGlobalConfig(current); err != nil {
		t.Fatal(err)
	}
}

func cleanupRemoteGame(t *testing.T, gameID string) {
	output, err := syncer.RemoveSaveGame(current, gameID)

	if err != nil {
		t.Logf("Error removing game from remote: %v\n%s\n", err, output)
	}
}

func pushGameRemote(t *testing.T, gameID string) {
	err := syncer.HandleSync(current, gameID, syncer.ModePush, false)

	if err != nil {
		t.Logf("Error pushing game: %v\n", err)
	}
}
