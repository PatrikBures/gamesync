package delete

import (
	"errors"
	"fmt"
	"go.pabu.dev/gamesync/internal/client/config"
	"go.pabu.dev/gamesync/internal/client/profiler"
	"strings"

	"github.com/spf13/cobra"
)

type profileCmd struct {
	cmd *cobra.Command
	opts profileOpts
}

type profileOpts struct {
	force bool
}

func newProfileCmd(conf *config.Config) *profileCmd {
	root := profileCmd{}

	cmd := &cobra.Command{
		Use: "profile PROFILE...",
		Short: "Delete profiles. Errors if any profile is not found",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProfileCmd(conf, &root.opts, args)
		},
	}

	cmd.Flags().BoolVarP(&root.opts.force, "force", "f", false, "Does not error if a profile is not found")

	root.cmd = cmd
	return &root
}

func runProfileCmd(conf *config.Config, opts *profileOpts, args []string) (err error) {
	p, err := profiler.New(conf.Global.ProfilesFile)
	if err != nil {
		return fmt.Errorf("initializing profiler: %w", err)
	}
	defer func() {
		err = errors.Join(err, p.Close())
	}()

	missingProfiles := []string{}

	for _, slug := range args {
		ok := p.Delete(slug)
		if ok {
			fmt.Println("d profile", slug)
		} else {
			if !opts.force {
				missingProfiles = append(missingProfiles, slug)
			} else {
				fmt.Println("Did not find profile", slug)
			}
		}
	}

	if !opts.force && len(missingProfiles) > 0 {
		fmt.Println("Changes were not saved")
		return fmt.Errorf("profile(s) not found: %s", strings.Join(missingProfiles, ", "))
	}

	return p.Save()
}
