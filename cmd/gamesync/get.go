package main

import (
	"context"
	"fmt"
	"gamesync/internal/client"
	"gamesync/internal/client/config"
	api "gamesync/internal/ogen"

	"github.com/spf13/cobra"
)

type getCmd struct {
	cmd *cobra.Command
}

func newGetCmd(config *config.Config) *getCmd {
	root := getCmd{}

	cmd := &cobra.Command{
		Use: "get",
		Short: "Get resource",
	}
	cmd.AddCommand(
		newGetUserCmd(config).cmd,
	)

	root.cmd = cmd
	return &root
}


type getUserCmd struct {
	cmd *cobra.Command
	opts getUserOpts
}
type getUserOpts struct {
	all bool
}

func newGetUserCmd(config *config.Config) *getUserCmd {
	root := getUserCmd{}

	cmd := &cobra.Command{
		Use: "user",
		Short: "Get all users or a specific one",
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := client.Client(*config)
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
		return fmt.Errorf("not implemented")
	}
	return nil
}

func runGetUserCmd(client *api.Client, opts getUserOpts) error {
	usersRes, err := client.GetUsers(context.Background())
	if err != nil {
		return err
	}
	switch res := usersRes.(type) {
	case *api.GetUsersOKApplicationJSON:
		for _, user := range *res {
			fmt.Println("name:", user.UserName, "ID:", user.UserID, "roleID:", user.RoleID)
		}
	case *api.GetUsersUnauthorized:
		return fmt.Errorf("unauthorized")
	case *api.GetUsersInternalServerError:
		return fmt.Errorf("server error: %v", res)
	default:
		return fmt.Errorf("unrecognized result: %v", res)
	}
	return nil
}

