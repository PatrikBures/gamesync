package sync

import (
	"fmt"

	"go.pabu.dev/gamesync/internal/client"
	util "go.pabu.dev/gamesync/internal/client/cmd/_util"
	"go.pabu.dev/gamesync/internal/client/config"
	"go.pabu.dev/gamesync/internal/client/profiler"
	"go.pabu.dev/gamesync/internal/client/syncer"
	api "go.pabu.dev/gamesync/internal/ogen"

	"github.com/spf13/cobra"
)

type syncCmd struct {
	Cmd *cobra.Command
	opts syncOpts
}

type syncOpts struct {
	profile profiler.Profile
}

func New(conf *config.Config) *syncCmd {
	root := syncCmd{}

	cmd := &cobra.Command{
		Use: "sync PROFILE",
		Short: "Syncs a profile with server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := populateSyncOpts(conf, &root.opts, args); err != nil {
				return err
			}

			c, err := client.New(conf)
			if err != nil { return err }

			return runSyncCmd(c, conf, &root.opts)
		},
	}

	root.Cmd = cmd
	return &root
}

func populateSyncOpts(conf *config.Config, opts *syncOpts, args []string) error {
	profileName := args[0]
	profile, ok, err := profiler.Get(profileName, conf.Global.ProfilesFile)
	if err != nil {
		return fmt.Errorf("getting profile '%s': %w", profileName, err)
	}
	if !ok {
		return fmt.Errorf("profile '%s' not found", profileName)
	}
	opts.profile = profile
	return nil
}

func runSyncCmd(c *api.Client, conf *config.Config, opts *syncOpts) error {
	s := syncer.New(conf, c, opts.profile)

	if err := s.Sync(syncer.ModeAuto); err != nil {
		return util.ErrHandler(err)
	}

	return nil
}
