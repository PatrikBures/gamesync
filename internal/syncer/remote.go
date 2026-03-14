package syncer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gamesync/internal/config"
	"gamesync/internal/state"
	"gamesync/internal/ui"
	"os/exec"
	"path"
	"strings"
)

func RemoveSaveGame(current config.Current, gameID string) (string, error) {
	output, err := RunCmd(current.Config.Server, "rm", "-r", 
		fmt.Sprintf("%s/%s/%s", config.RemoteSavesDir, current.Config.Server.User, gameID))

	if err != nil {
		return output, err
	}

	return output, nil
}

func RunCmd(server config.ServerConfig, cmds ...string) (string, error) {
	var sshArgs []string

	if server.SshHost == "" {
		sshArgs = []string{
			"-p", server.Port,
			"-i", server.IdentityFile,
			fmt.Sprintf("%s@%s", server.User, server.Host),
		}
	} else {
		sshArgs = []string{server.SshHost}
	}

	sshArgs = append(sshArgs, cmds...)

	cmd := exec.Command("ssh", sshArgs...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()

	outStr := strings.TrimSuffix(stdoutBuf.String(), "\n")
	errStr := strings.TrimSuffix(stderrBuf.String(), "\n")

	ui.Debug("Ran cmd:\n%s\n", cmd.String())
	ui.Debug("output:\n%s\n", outStr)
	ui.Debug("error:\n%s\n", errStr)

	if err != nil {
		return errStr, fmt.Errorf("failed to run cmd: %w", err)
	}

	return outStr, nil
}

func GetRemoteState(current config.Current, gameID string) (map[string]state.FileMeta, error) {
	output, err := RunCmd(current.Config.Server, 
		"gamesync-state",
		"/"+path.Join("data", "saves", current.Config.Server.User, gameID))
	if err != nil {
		return nil, err
	}

	s := make(map[string]state.FileMeta)

	if err := json.Unmarshal([]byte(output), &s); err != nil {
		return nil, err
	}

	return s, nil
}
