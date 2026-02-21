package main

import (
	"bytes"
	"gamesync/internal/ui"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoteLs(t *testing.T) {
	setupTest(t)

	pushGameRemote(t, "game_1")

	var buf bytes.Buffer

	oldOut := ui.OutWriter
	ui.OutWriter = &buf

	oldLevel := ui.GetLevel()
	ui.SetLevel(ui.LevelNormal)

	t.Cleanup(func() {
		ui.OutWriter = oldOut
		ui.SetLevel(oldLevel)
	})

	cmd := newRemoteLsCmd()

	cmd.cmd.SilenceUsage = true
	cmd.cmd.SilenceErrors = true

	require.NoError(t, cmd.cmd.Execute())
	require.Equal(t, "game_1\n", buf.String())
}
