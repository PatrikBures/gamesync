package main

import (
	"bytes"
	"fmt"
	"gamesync/internal/config"
	"gamesync/internal/ui"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSaveLs(t *testing.T) {
	setupTest(t)

	var buf bytes.Buffer

	oldOut := ui.OutWriter
	ui.OutWriter = &buf

	oldLevel := ui.GetLevel()
	ui.SetLevel(ui.LevelNormal)

	t.Cleanup(func() {
		ui.OutWriter = oldOut
		ui.SetLevel(oldLevel)
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




func TestSaveAdd(t *testing.T) {
	setupTest(t)

	var cmd *saveAddCmd
	var err error
	var game config.GameConfig

	game, _, err = config.GetGame(current, "game_1")
	require.NoError(t, err)

	cmd = testSaveAddCmd(game.ID, game.SavePath, false)
	require.Error(t, cmd.cmd.Execute())

	cmd = testSaveAddCmd(game.ID, game.SavePath, true)
	require.NoError(t, cmd.cmd.Execute())



	_, _, err = config.GetGame(current, "game_2")
	require.Error(t, err)

	// using game_1's savepath for testing, there is not a variable currently which stores game_2's save path
	cmd = testSaveAddCmd("game_2", game.SavePath, true)
	require.NoError(t, cmd.cmd.Execute())

	_, _, err = config.GetGame(current, "game_2")
	require.NoError(t, err)
}

func testSaveAddCmd(gameID string, savePath string, update bool) *saveAddCmd {
	cmd := newSaveAddCmd()

	cmd.cmd.SilenceUsage = true
	cmd.cmd.SilenceErrors = true

	cmd.opts.update = update

	cmd.cmd.SetArgs([]string{gameID, savePath})

	return cmd
}



func TestSaveRm(t *testing.T) {
	setupTest(t)

	var cmd *saveRmCmd

	cmd = testSaveRmCmd("game_1")
	require.NoError(t, cmd.cmd.Execute())

	cmd = testSaveRmCmd("game_1")
	require.Error(t, cmd.cmd.Execute())
}

func testSaveRmCmd(gameID string) *saveRmCmd {
	cmd := newSaveRmCmd()

	cmd.cmd.SilenceUsage = true
	cmd.cmd.SilenceErrors = true

	cmd.cmd.SetArgs([]string{gameID})

	return cmd
}
