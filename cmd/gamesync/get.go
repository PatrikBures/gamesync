package main

import (
	"context"
	"fmt"
	"gamesync/internal/client"
	"gamesync/internal/client/config"
	api "gamesync/internal/ogen"
	"strconv"

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
	userID int64
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
	result, err := client.GetUsers(context.Background())
	if err != nil {
		return err
	}
	switch r := result.(type) {
	case *api.GetUsersOKApplicationJSON:
		for _, user := range *r {
			fmt.Println("name:", user.UserName, "ID:", user.UserID, "roleID:", user.RoleID)
		}
	case *api.GetUsersUnauthorized:
		return fmt.Errorf("unauthorized")
	case *api.GetUsersInternalServerError:
		return fmt.Errorf("server error: %v", r)
	default:
		return fmt.Errorf("unrecognized type %T with result: %v", r, r)
	}
	return nil
}

func oneUser(client *api.Client, userID int64) error {
	result, err := client.GetUser(
		context.Background(), 
		api.GetUserParams{UserID: userID},
	)
	if err != nil {
		return err
	}
	switch r := result.(type) {
	case *api.User:
		user := *r
		fmt.Println("name:", user.UserName, "ID:", user.UserID, "roleID:", user.RoleID)
	case *api.GetUserUnauthorized:
		return fmt.Errorf("unauthorized")
	case *api.GetUserNotFound:
		return fmt.Errorf("user not found")
	case *api.GetUserInternalServerError:
		return fmt.Errorf("server error: %v", r)
	default:
		return fmt.Errorf("unrecognized type %T with result: %v", r, r)
	}
	return nil
}
