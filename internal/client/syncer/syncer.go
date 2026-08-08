package syncer

import (
	"os"

	"go.pabu.dev/gamesync/internal/client/config"
	"go.pabu.dev/gamesync/internal/client/profiler"
	"go.pabu.dev/gamesync/internal/client/stater"
	api "go.pabu.dev/gamesync/internal/ogen"

	"github.com/klauspost/compress/zstd"
)

type syncer struct {
	conf    *config.Config
	client  *api.Client
	profile profiler.Profile
	stater  *stater.Stater
}

type puller struct {
	client           *api.Client
	decoder          *zstd.Decoder
	file             *os.File
	fileBytesWritten int64
	chunkDir         string
}

func New(conf *config.Config, c *api.Client, profile profiler.Profile) *syncer {
	return &syncer{
		client:  c,
		conf:    conf,
		profile: profile,
		stater: stater.New(conf.ProfileStateDir()),
	}
}
