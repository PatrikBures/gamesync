package get

import (
	"bytes"
	"errors"
	"fmt"
	"go.pabu.dev/gamesync/internal/client/config"
	"go.pabu.dev/gamesync/internal/client/profiler"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type profileCmd struct {
	cmd *cobra.Command
}

func newProfileCmd(conf *config.Config) *profileCmd {
	root := profileCmd{}

	cmd := &cobra.Command{
		Use: "profile [PROFILE]...",
		Short: "Get profile(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProfileCmd(conf, args)
		},
	}

	root.cmd = cmd
	return &root
}

func runProfileCmd(conf *config.Config, args []string) (err error) {
	p, err := profiler.New(conf.Global.ProfilesFile)
	if err != nil {
		return fmt.Errorf("initializing profiler: %w", err)
	}
	defer func() {
		err = errors.Join(err, p.Close())
	}()

	if len(p.Profiles) == 0 {
		fmt.Println("No profiles found")
		return nil
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 4, 2, ' ', 0)

	if _, err := fmt.Fprintf(w, "Name\tRepo\tBranch\tDir\n"); err != nil { return err }
	if len(args) > 0 {
		for _, slug := range args {
			profile, ok := p.Get(slug)
			if !ok {
				return fmt.Errorf("Profile not found: %s", slug)
			}
			if err := printProfile(w, slug, profile); err != nil { return err }
		}

	} else {
		for slug, profile := range p.Profiles {
			if err := printProfile(w, slug, profile); err != nil { return err }
		}
	}
	if err := w.Flush(); err != nil { return err }

	_, err = buf.WriteTo(os.Stdout)
	return err
}

func printProfile(w *tabwriter.Writer, slug string, profile profiler.Profile) error {
	_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", slug, profile.RepoName, profile.BranchName, profile.Dir)
	return err
}

