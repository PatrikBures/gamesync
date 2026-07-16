package client

import (
	"context"
	"fmt"
	clientAuth "gamesync/internal/client/auth"
	config "gamesync/internal/client/config"
	api "gamesync/internal/ogen"
)

// if UserID in config is < 1, request one from server
func Client(c *config.Config) (*api.Client, error) {
	client, err := api.NewClient(
		c.Server.Url,
		clientAuth.NewAuth(c.Server.Token))
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}
	if c.Server.UserID < 1 {
		user, err := client.GetMe(context.Background())
		if err != nil {
			return nil, fmt.Errorf("getting userid: %w", err)
		}
		c.Server.UserID = user.UserID
		config.SetCacheItem(c.Global.CacheDir, config.CacheNameUserID, user.UserID)
	}
	return client, nil
}
