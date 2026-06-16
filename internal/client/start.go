package client

import (
	"fmt"
	clientAuth "gamesync/internal/client/auth"
	clientConfig "gamesync/internal/client/config"
	api "gamesync/internal/ogen"
)

func Client(config clientConfig.Config) (*api.Client, error) {
	client, err := api.NewClient(
		config.Server,
		clientAuth.NewAuth(config.Token))
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}
	return client, nil
}
