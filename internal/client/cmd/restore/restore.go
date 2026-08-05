package restore

import (
	"fmt"
	"go.pabu.dev/gamesync/internal/client"
	"go.pabu.dev/gamesync/internal/client/config"
	"go.pabu.dev/gamesync/internal/client/profiler"
	"go.pabu.dev/gamesync/internal/client/syncer"
	api "go.pabu.dev/gamesync/internal/ogen"
	"strconv"

	"github.com/spf13/cobra"
)


type restoreCmd struct {
	Cmd *cobra.Command
	opts restoreOpts
}
type restoreOpts struct {
	profile profiler.Profile
	snapshotID int64
}

func New(conf *config.Config) *restoreCmd {
	root := restoreCmd{}

	cmd := &cobra.Command{
		Use:   "restore PROFILE SNAPSHOT_ID",
		Short: "Restore a sync to a specific snapshot",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := populateRestoreOpts(conf, &root.opts, args); err != nil { return err }

			c, err := client.New(conf)
			if err != nil { return err }

			if err := runRestoreCmd(c, root.opts, conf); err != nil { return err }

			return nil
		},
	}

	root.Cmd = cmd
	return &root
}

func populateRestoreOpts(conf *config.Config, opts *restoreOpts, args []string) error {
	profile, ok, err := profiler.Get(args[0], conf.Global.ProfilesFile)
	if err != nil {
		return fmt.Errorf("initializing profiler: %s", err)
	}
	if !ok {
		return fmt.Errorf("profile '%s' does not exist", args[0])
	}
	opts.profile = profile

	opts.snapshotID, err = strconv.ParseInt(args[1], 10, 64)
	if err != nil { return err }

	return nil
}

func runRestoreCmd(c *api.Client, opts restoreOpts, conf *config.Config) (err error) {
	syncer := syncer.New(conf, c, opts.profile)

	if err := syncer.Pull(opts.snapshotID); err != nil {
		return fmt.Errorf("pulling: %w", err)
	}

	return nil
}
