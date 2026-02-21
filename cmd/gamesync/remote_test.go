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

	t.Cleanup(func() {
		ui.OutWriter = oldOut
	})

	cmd := newRemoteLsCmd()

	cmd.cmd.SilenceUsage = true
	cmd.cmd.SilenceErrors = true

	require.Equal(t, "game_1\n", buf.String())
}
