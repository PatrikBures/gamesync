package push

import (
	"fmt"

	"go.pabu.dev/gamesync/internal/client"
	"go.pabu.dev/gamesync/internal/client/config"
	"go.pabu.dev/gamesync/internal/client/profiler"
	"go.pabu.dev/gamesync/internal/client/syncer"
	api "go.pabu.dev/gamesync/internal/ogen"

	"github.com/spf13/cobra"
)

type pushCmd struct {
	Cmd *cobra.Command
	opts pushOpts
}

type pushOpts struct {
	profile profiler.Profile
	force bool
	mode syncer.SyncMode
}

func New(conf *config.Config) *pushCmd {
	root := pushCmd{}

	cmd := &cobra.Command{
		Use: "push PROFILE",
		Short: "Pushes a profile to server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := populatePushOpts(conf, &root.opts, args); err != nil {
				return err
			}

			c, err := client.New(conf)
			if err != nil { return err }

			return runPushCmd(c, conf, &root.opts)
		},
	}

	cmd.Flags().BoolVarP(&root.opts.force, "force", "f", false, "Force pushes to server")

	root.Cmd = cmd
	return &root
}

func populatePushOpts(conf *config.Config, opts *pushOpts, args []string) error {
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
		opts.mode = syncer.ModePushForce
	} else {
		opts.mode = syncer.ModePush
	}

	return nil
}

func runPushCmd(c *api.Client, conf *config.Config, opts *pushOpts) error {
	s := syncer.New(conf, c, opts.profile)

	return s.Sync(opts.mode)
}

