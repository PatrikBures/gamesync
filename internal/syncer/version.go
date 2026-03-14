package syncer

import (
	"fmt"
	"gamesync/internal/config"
	"gamesync/internal/ui"
)

func SameApiVersion(server config.ServerConfig) error {
	remoteApiVersion, err := RunCmd(server, false, "api-version")
	if err != nil {
		return fmt.Errorf("failed getting remote version: %w", err)
	}
	if remoteApiVersion != config.ApiVersion {
		return fmt.Errorf("api version numbers do not match! client version: '%s', remote version: '%s'", config.ApiVersion, remoteApiVersion)
	}
	ui.Debug("Client and remote api versions match, client: %s, remote: %s\n", config.ApiVersion, remoteApiVersion)
	return nil
}
