package wrap

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"go.pabu.dev/gamesync/internal/client"
	"go.pabu.dev/gamesync/internal/client/config"
	"go.pabu.dev/gamesync/internal/client/profiler"
	"go.pabu.dev/gamesync/internal/client/syncer"
	api "go.pabu.dev/gamesync/internal/ogen"

	"github.com/spf13/cobra"
)

type wrapCmd struct {
	Cmd *cobra.Command
	opts wrapOpts
}

type wrapOpts struct {
	profile profiler.Profile
	noPull bool
	noPush bool
	cmd string
	cmdArgs []string
}

func New(conf *config.Config) *wrapCmd {
	root := wrapCmd{}

	cmd := &cobra.Command{
		Use: "wrap PROFILE -- CMD...",
		Short: "Wraps any command, pulling before process starts and pushing after the process exits",
		ArgAliases: []string{"--"},
		RunE: func(cmd *cobra.Command, args []string) error {
			dashIdx := cmd.ArgsLenAtDash()
			if dashIdx == -1 {
				return fmt.Errorf("found no '--'")
			}
			userArgs := args[:dashIdx]
			cmdArgs := args[dashIdx:]

			if len(userArgs) < 1 {
				return fmt.Errorf("missing PROFILE")
			}
			if len(cmdArgs) < 1 {
				return fmt.Errorf("missing command after '--'")
			}

			if err := populateWrapOpts(conf, &root.opts, userArgs, cmdArgs); err != nil {
				return err
			}

			c, err := client.New(conf)
			if err != nil { return err }

			return runWrapCmd(c, conf, &root.opts)
		},
	}

	cmd.Flags().BoolVar(&root.opts.noPull, "no-pull", false, "Does not pull before running cmd")
	cmd.Flags().BoolVar(&root.opts.noPush, "no-push", false, "Does not push after cmd exits")

	root.Cmd = cmd
	return &root
}

func populateWrapOpts(conf *config.Config, opts *wrapOpts, args []string, cmdArgs []string) error {
	profileName := args[0]
	profile, ok, err := profiler.Get(profileName, conf.Global.ProfilesFile)
	if err != nil {
		return fmt.Errorf("getting profile '%s': %w", profileName, err)
	}
	if !ok {
		return fmt.Errorf("profile '%s' not found", profileName)
	}
	opts.profile = profile


	opts.cmd = cmdArgs[0]
	if len(cmdArgs) > 1 {
		opts.cmdArgs = cmdArgs[1:]
	}

	return nil

}

func runWrapCmd(c *api.Client, conf *config.Config, opts *wrapOpts) error {
	s := syncer.New(conf, c, opts.profile)

	if !opts.noPull {
		if err := s.Sync(syncer.ModePull); err != nil {
			if errors.Is(err, syncer.ErrNoExistingSnapshot) {
				fmt.Println("nothing to pull")
			} else {
				return fmt.Errorf("pulling: %w", err)
			}
		}
	}

	fmt.Printf("about to execute: %s %v\n", opts.cmd, opts.cmdArgs)

	cmd := exec.Command(opts.cmd, opts.cmdArgs...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	exitCode := 0
	var exitErr error
	if err := cmd.Run(); err != nil {
		if exitError, ok := errors.AsType[*exec.ExitError](err); ok {
			exitCode = exitError.ExitCode()
			exitErr = exitError
		} else {
			exitErr = err
			exitCode = 1
		}
	}

	if exitErr != nil {
		return fmt.Errorf("cmd received error with code '%d': %w", exitCode, exitErr)
	}

	if !opts.noPush {
		if err := s.Sync(syncer.ModePush); err != nil {
			return fmt.Errorf("pushing: %w", err)
		}
	}

	fmt.Println("cmd exited with code", exitCode)

	return nil
}
