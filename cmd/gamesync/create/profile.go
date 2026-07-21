package create

import (
	"errors"
	"fmt"
	"gamesync/internal/client/config"
	"gamesync/internal/client/profiler"

	"github.com/spf13/cobra"
)

type profileCmd struct {
	cmd *cobra.Command
	opts profileOpts
}
type profileOpts struct {
	slug    string
	force   bool
	profile profiler.Profile
}

func newProfileCmd(conf *config.Config) *profileCmd {
	root := profileCmd{}

	cmd := &cobra.Command{
		Use: "profile PROFILE REPO BRANCH DIR",
		Short: "Create new profile to sync",
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			populateProfileOpts(&root.opts, args)

			if err := runProfileCmd(root.opts, conf); err != nil { return err }

			return nil
		},
	}
	cmd.Flags().BoolVarP(&root.opts.force, "force", "f", false, "Overwrites any exising profile")
	root.cmd = cmd
	return &root
}

func populateProfileOpts(opts *profileOpts, args []string) {
	opts.slug = args[0]
	opts.profile = profiler.Profile{
		RepoName: args[1],
		BranchName: args[2],
		Dir: args[3],
	}
}

func runProfileCmd(opts profileOpts, conf *config.Config) (err error) {
	pr, err := profiler.New(conf.Global.ProfilesFile)
	if err != nil {
		return fmt.Errorf("initializing profiler: %w", err)
	}
	defer func() {
		err = errors.Join(err, pr.Close())
	}()

	if opts.force {
		pr.AddOverwrite(opts.slug, opts.profile)
	} else {
		if err := pr.Add(opts.slug, opts.profile); err != nil { return err }
	}

	if opts.force && pr.Exists(opts.slug) {
		fmt.Printf("Force created profile '%s'\n", opts.slug)
	} else {
		fmt.Printf("Created profile '%s'\n", opts.slug)
	}

	if err := pr.Save(); err != nil { return err }

	return nil
}
