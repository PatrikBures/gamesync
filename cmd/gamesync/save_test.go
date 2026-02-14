package main

import (
	"bytes"
	"fmt"
	"gamesync/internal/ui"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSaveLs(t *testing.T) {
	setupTest(t)

	var buf bytes.Buffer

	oldOut := ui.OutWriter
	ui.OutWriter = &buf

	t.Cleanup(func() {
		ui.OutWriter = oldOut
	})

	cmd := newSaveLsCmd()

	cmd.opts.quiet = true

	cmd.cmd.SilenceUsage = true
	cmd.cmd.SilenceErrors = true

	err := cmd.cmd.Execute()
	require.NoError(t, err)

	require.Equal(t, getExpectedLs(), buf.String())
}

func getExpectedLs() string {
	var out string
	for _, game := range current.Config.Games {
		out = fmt.Sprintf("%s%s\n", out, game.ID)
	}
	return out
}
