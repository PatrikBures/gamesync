package client

import (
	"context"
	"fmt"
	clientAuth "gamesync/internal/client/auth"
	config "gamesync/internal/client/config"
	api "gamesync/internal/ogen"
)

// if UserID in config is < 1, request one from server
func Client(config *config.Config) (*api.Client, error) {
	client, err := api.NewClient(
		config.Server.Url,
		clientAuth.NewAuth(config.Server.Token))
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}
	if config.Server.UserID < 1 {
		user, err := client.GetMe(context.Background())
		if err != nil {
			return nil, fmt.Errorf("getting userid: %w", err)
		}
		fmt.Println("got uid:", user.UserID)
		config.Server.UserID = user.UserID
	}
	return client, nil
}
