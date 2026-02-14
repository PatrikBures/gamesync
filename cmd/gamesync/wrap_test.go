package main

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestWrap(t *testing.T) {
	setupTest(t)
	cmd := newWrapCmd()

	cmd.cmd.SilenceUsage = true
	cmd.cmd.SilenceErrors = true

	cmd.opts.exitOnError = true
	cmd.cmd.SetArgs([]string{"game_1", "--", "sleep 1"})
	require.NoError(t, cmd.cmd.Execute())
}
