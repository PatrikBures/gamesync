package syncer

import (
	"fmt"
	"gamesync/internal/config"
	"gamesync/internal/ui"
	"gamesync/internal/vars"
)

func SameApiVersion(server config.ServerConfig) error {
	if vars.ApiVersion == "0" {
		ui.Debug("ignored checking api version\n")
		return nil
	}
	remoteApiVersion, err := RunCmd(server, false, "api-version")
	if err != nil {
		return fmt.Errorf("failed getting remote version: %w", err)
	}
	if remoteApiVersion != vars.ApiVersion {
		return fmt.Errorf("api version numbers do not match! client version: '%s', remote version: '%s'", vars.ApiVersion, remoteApiVersion)
	}
	ui.Debug("Client and remote api versions match, client: %s, remote: %s\n", vars.ApiVersion, remoteApiVersion)
	return nil
}
