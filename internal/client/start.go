package client

import (
	"fmt"
	clientAuth "gamesync/internal/client/auth"
	config "gamesync/internal/client/config"
	api "gamesync/internal/ogen"
)

func Client(config config.Config) (*api.Client, error) {
	client, err := api.NewClient(
		config.Server.Url,
		clientAuth.NewAuth(config.Server.Token))
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}
	return client, nil
}
