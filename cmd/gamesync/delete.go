package main

import (
	"errors"
	"fmt"
	"gamesync/internal/client/config"
	"gamesync/internal/client/profiler"
	"strings"

	"github.com/spf13/cobra"
)


type deleteCmd struct {
	cmd *cobra.Command
}

func newDeleteCmd(conf *config.Config) *deleteCmd {
	root := deleteCmd{}

	cmd := &cobra.Command{
		Use: "delete",
		Short: "Delete resources",
	}

	cmd.AddCommand(
		newDeleteProfileCmd(conf).cmd,
	)

	root.cmd = cmd
	return &root
}

type deleteProfileCmd struct {
	cmd *cobra.Command
	opts deleteProfileOpts
}

type deleteProfileOpts struct {
	force bool
}

func newDeleteProfileCmd(conf *config.Config) *deleteProfileCmd {
	root := deleteProfileCmd{}

	cmd := &cobra.Command{
		Use: "profile PROFILE...",
		Short: "Delete profiles. Errors if any profile is not found",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeleteProfileCmd(conf, &root.opts, args)
		},
	}

	cmd.Flags().BoolVarP(&root.opts.force, "force", "f", false, "Does not error if a profile is not found")

	root.cmd = cmd
	return &root
}

func runDeleteProfileCmd(conf *config.Config, opts *deleteProfileOpts, args []string) (err error) {
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
			fmt.Println("Deleted profile", slug)
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
