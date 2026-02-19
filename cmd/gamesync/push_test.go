package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPush(t *testing.T) {
	setupTest(t)

	cmd :=newPushCmd()

	t.Cleanup(func() {
		cleanupRemoteGame(t, "game_1")
	})

	cmd.cmd.SilenceUsage = true
	cmd.cmd.SilenceErrors = true

	cmd.cmd.SetArgs([]string{"game_1"})

	require.NoError(t, cmd.cmd.Execute())
}
