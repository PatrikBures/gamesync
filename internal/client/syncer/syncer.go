package syncer

import (
	"go.pabu.dev/gamesync/internal/client/config"
	"go.pabu.dev/gamesync/internal/client/profiler"
	"go.pabu.dev/gamesync/internal/client/stater"
	api "go.pabu.dev/gamesync/internal/ogen"
)

type syncer struct {
	conf    *config.Config
	client  *api.Client
	profile profiler.Profile
	stater  *stater.Stater
}

func New(conf *config.Config, c *api.Client, profile profiler.Profile) *syncer {
	return &syncer{
		client:  c,
		conf:    conf,
		profile: profile,
		stater: stater.New(conf.ProfileStateDir()),
	}
}
