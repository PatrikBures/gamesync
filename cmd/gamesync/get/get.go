package get

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	util "gamesync/cmd/gamesync/_util"
	"gamesync/internal/client"
	"gamesync/internal/client/config"
	"gamesync/internal/client/profiler"
	api "gamesync/internal/ogen"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type getCmd struct {
	Cmd *cobra.Command
}

func New(conf *config.Config) *getCmd {
	root := getCmd{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get resource",
	}
	cmd.AddCommand(
		newGetUserCmd(conf).cmd,
		newGetProfileCmd(conf).cmd,
	)

	root.Cmd = cmd
	return &root
}

type getUserCmd struct {
	cmd  *cobra.Command
	opts getUserOpts
}
type getUserOpts struct {
	all    bool
	userID int64
}

func newGetUserCmd(config *config.Config) *getUserCmd {
	root := getUserCmd{}

	cmd := &cobra.Command{
		Use:   "user",
		Short: "Get all users or a specific one",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := client.New(config)
			if err != nil {
				return err
			}
			if err := populateGetUserOpts(&root.opts, args); err != nil {
				return err
			}
			if err := runGetUserCmd(client, root.opts); err != nil {
				return err
			}
			return nil
		},
	}

	root.cmd = cmd
	return &root
}

func populateGetUserOpts(opts *getUserOpts, args []string) error {
	if len(args) == 0 {
		opts.all = true
	} else {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		opts.userID = id
	}
	return nil
}

func runGetUserCmd(client *api.Client, opts getUserOpts) (err error) {
	if opts.all {
		err = allUsers(client)
	} else {
		err = oneUser(client, opts.userID)
	}
	return err
}

func allUsers(client *api.Client) error {
	users, err := client.GetUsers(context.Background())
	if err := util.ErrHandler(err); err != nil {
		return err
	}
	for _, user := range users {
		fmt.Println("name:", user.UserName, "ID:", user.UserID, "roleID:", user.RoleID)
	}
	return nil
}

func oneUser(client *api.Client, userID int64) error {
	user, err := client.GetUser(
		context.Background(),
		api.GetUserParams{UserID: userID},
	)
	if err := util.ErrHandler(err); err != nil {
		return err
	}
	fmt.Println("name:", user.UserName, "ID:", user.UserID, "roleID:", user.RoleID)
	return nil
}




type getProfileCmd struct {
	cmd *cobra.Command
	opts getProfileOpts
}
type getProfileOpts struct {
}

func newGetProfileCmd(conf *config.Config) *getProfileCmd {
	root := getProfileCmd{}

	cmd := &cobra.Command{
		Use: "profile [PROFILE]...",
		Short: "Get profile(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetProfileCmd(conf, args)
		},
	}

	root.cmd = cmd
	return &root
}

func runGetProfileCmd(conf *config.Config, args []string) (err error) {
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

