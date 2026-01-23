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

func RemoveSaveGame(current config.Current, gameID string) (string, error) {
	output, err := RunCmd(current, "rm", "-r", 
		fmt.Sprintf("/data/saves/%s/%s", current.Config.Server.User, gameID))

	if err != nil {
		return output, err
	}

	return output, nil
}

func RunCmd(current config.Current, cmds ...string) (string, error) {
	sshArgs := []string{
		"-p", current.Config.Server.Port,
		"-i", current.Config.Server.IdentityFile,
		fmt.Sprintf("%s@%s", current.Config.Server.User, current.Config.Server.Host),
	}

	sshArgs = append(sshArgs, cmds...)

	cmd := exec.Command("ssh", sshArgs...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()

	outStr := stdoutBuf.String()
	errStr := stderrBuf.String()

	if current.Verbose {
		fmt.Printf("Ran cmd:\n%s\n", cmd.String())
		fmt.Printf("output:\n%s\n", outStr)
		fmt.Printf("error:\n%s\n", errStr)
	}

	if err != nil {
		return errStr, fmt.Errorf("failed to run cmd: %w", err)
	}

	return outStr, nil
}

func GetRemoteState(current config.Current, gameID string) (map[string]state.FileMeta, error) {
	output, err := RunCmd(current, 
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
