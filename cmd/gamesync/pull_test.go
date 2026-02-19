package main

import (
	"gamesync/internal/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPull(t *testing.T) {
	setupTest(t)

	pushGameRemote(t, "game_1")

	require.NoError(t, config.RemoveGames(&current, []string{"game_1"}))

	cmd :=newPullCmd()

	t.Cleanup(func() {
		cleanupRemoteGame(t, "game_1")
	})

	cmd.cmd.SilenceUsage = true
	cmd.cmd.SilenceErrors = true

	cmd.cmd.SetArgs([]string{"game_1"})

	require.NoError(t, cmd.cmd.Execute())
}
