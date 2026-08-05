package get

import (
	"context"
	"fmt"
	"go.pabu.dev/gamesync/internal/client"
	util "go.pabu.dev/gamesync/internal/client/cmd/_util"
	"go.pabu.dev/gamesync/internal/client/config"
	api "go.pabu.dev/gamesync/internal/ogen"
	"strconv"

	"github.com/spf13/cobra"
)

type userCmd struct {
	cmd  *cobra.Command
	opts userOpts
}
type userOpts struct {
	all    bool
	userID int64
}

func newUserCmd(config *config.Config) *userCmd {
	root := userCmd{}

	cmd := &cobra.Command{
		Use:   "user",
		Short: "Get all users or a specific one",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := client.New(config)
			if err != nil {
				return err
			}
			if err := populateUserOpts(&root.opts, args); err != nil {
				return err
			}
			if err := runUserCmd(client, root.opts); err != nil {
				return err
			}
			return nil
		},
	}

	root.cmd = cmd
	return &root
}

func populateUserOpts(opts *userOpts, args []string) error {
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

func runUserCmd(client *api.Client, opts userOpts) (err error) {
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
