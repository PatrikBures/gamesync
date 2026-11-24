package syncer

import (
	"go.pabu.dev/gamesync/internal/client/config"
	"go.pabu.dev/gamesync/internal/client/profiler"
	"go.pabu.dev/gamesync/internal/client/stater"
	api "go.pabu.dev/gamesync/internal/ogen"
	"go.pabu.dev/packer"
)

type syncer struct {
	conf    *config.Config
	client  *api.Client
	profile profiler.Profile
	stater  *stater.Stater
	p       *packer.Packager
}

func New(conf *config.Config, c *api.Client, profile profiler.Profile, p *packer.Packager) *syncer {
	return &syncer{
		client:  c,
		conf:    conf,
		profile: profile,
		stater:  stater.New(conf.ProfileStateDir()),
		p:       p,
	}
}
