package syncer

import (
	"go.pabu.dev/gamesync/internal/client/config"
	"go.pabu.dev/gamesync/internal/client/profiler"
	api "go.pabu.dev/gamesync/internal/ogen"
	"os"

	"github.com/klauspost/compress/zstd"
)

type syncer struct {
	conf    *config.Config
	client  *api.Client
	profile profiler.Profile
	stater  *stater
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
		stater: NewStater(conf.ProfileStateDir()),
	}
}
