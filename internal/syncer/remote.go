package syncer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gamesync/internal/config"
	"gamesync/internal/state"
	"os/exec"
	"path"
)

func RemoveSaveGame(gameID string, server config.ServerConfig, verbose bool) (string, error) {
	output, err := RunCmd(server, verbose, "rm", "-r", 
		fmt.Sprintf("/data/saves/%s/%s", server.User, gameID))

	if err != nil {
		return output, err
	}

	return output, nil
}

func RunCmd(server config.ServerConfig, verbose bool, cmds ...string) (string, error) {
	sshArgs := []string{
		"-p", server.Port,
		"-i", server.IdentityFile,
		fmt.Sprintf("%s@%s", server.User, server.Host),
	}

	sshArgs = append(sshArgs, cmds...)

	cmd := exec.Command("ssh", sshArgs...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()

	outStr := stdoutBuf.String()
	errStr := stderrBuf.String()

	if verbose {
		fmt.Printf("Ran cmd:\n%s\n", cmd.String())
		fmt.Printf("output:\n%s\n", outStr)
		fmt.Printf("error:\n%s\n", errStr)
	}

	if err != nil {
		return errStr, fmt.Errorf("failed to run cmd: %w", err)
	}

	return outStr, nil
}

func GetRemoteState(gameID string, server config.ServerConfig, verbose bool) (map[string]state.FileMeta, error) {
	output, err := RunCmd(server, verbose, "gamesync-state", "/"+path.Join("data", "saves", server.User, gameID))
	if err != nil {
		return nil, err
	}

	s := make(map[string]state.FileMeta)

	if err := json.Unmarshal([]byte(output), &s); err != nil {
		return nil, err
	}

	return s, nil
}
