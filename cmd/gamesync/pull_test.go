package main

import (
	"gamesync/internal/config"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPull(t *testing.T) {
	setupTest(t)

	pushGameRemote(t, "game_1")

	game, _, err := config.GetGame(current, "game_1")
	require.NoError(t, err)
	require.NoError(t, os.RemoveAll(game.SavePath))

	cmd :=newPullCmd()

	t.Cleanup(func() {
		cleanupRemoteGame(t, "game_1")
	})

	cmd.cmd.SilenceUsage = true
	cmd.cmd.SilenceErrors = true

	cmd.cmd.SetArgs([]string{"game_1"})

	require.NoError(t, cmd.cmd.Execute())
}
