package pull

import (
	"fmt"

	"go.pabu.dev/gamesync/internal/client"
	"go.pabu.dev/gamesync/internal/client/config"
	"go.pabu.dev/gamesync/internal/client/profiler"
	"go.pabu.dev/gamesync/internal/client/syncer"
	api "go.pabu.dev/gamesync/internal/ogen"

	"github.com/spf13/cobra"
)

type pullCmd struct {
	Cmd *cobra.Command
	opts pullOpts
}

type pullOpts struct {
	profile profiler.Profile
	force bool
	mode syncer.SyncMode
}

func New(conf *config.Config) *pullCmd {
	root := pullCmd{}

	cmd := &cobra.Command{
		Use: "pull PROFILE",
		Short: "Pulls a profile from server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := populatePullOpts(conf, &root.opts, args); err != nil {
				return err
			}

			c, err := client.New(conf)
			if err != nil { return err }

			return runPullCmd(c, conf, &root.opts)
		},
	}

	cmd.Flags().BoolVarP(&root.opts.force, "force", "f", false, "Force pulls from server")

	root.Cmd = cmd
	return &root
}

func populatePullOpts(conf *config.Config, opts *pullOpts, args []string) error {
	profileName := args[0]
	profile, ok, err := profiler.Get(profileName, conf.Global.ProfilesFile)
	if err != nil {
		return fmt.Errorf("getting profile '%s': %w", profileName, err)
	}
	if !ok {
		return fmt.Errorf("profile '%s' not found", profileName)
	}
	opts.profile = profile

	if opts.force {
		opts.mode = syncer.ModePullForce
	} else {
		opts.mode = syncer.ModePull
	}

	return nil
}

func runPullCmd(c *api.Client, conf *config.Config, opts *pullOpts) error {
	s := syncer.New(conf, c, opts.profile)

	return s.Sync(opts.mode)
}

