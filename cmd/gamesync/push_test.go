package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPush(t *testing.T) {
	setupTest(t)

	cmd :=newPushCmd()

	cmd.cmd.SetArgs([]string{"game1"})

	require.NoError(t, cmd.cmd.Execute())
}
