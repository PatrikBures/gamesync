package syncer

import (
	"bytes"
	"fmt"
	"gamesync/internal/config"
	"os/exec"
)

func RunCmd(server config.ServerConfig, verbose bool, cmds ...string) (string, error){
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
