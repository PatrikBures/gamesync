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

	ui.SetLevel(ui.LevelDebug)

	current.Config = config.Config{}
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

	for i := range 2 {
		saveName := "game_"+strconv.Itoa(i)
		savePath := filepath.Join(tmpDir, saveName)
		odd := i % 2 == 1

		createTestGame(t, saveName, savePath, odd)
	}
}


func createTestGame(t *testing.T, ID, savePath string, appendGame bool) {
	game := config.GameConfig{
		ID: ID,
		SavePath: savePath,
	}

	if appendGame {
		current.Config.Games = append(current.Config.Games, game)
	}

	if err := os.Mkdir(savePath, 0770); err != nil {
		t.Fatal(err)
	}

	testFile := filepath.Join(savePath, "save.txt")

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
	_, err := syncer.HandleSync(current, gameID, syncer.ModePush, false, false, false)

	if err != nil {
		t.Logf("Error pushing game: %v\n", err)
	}
}
