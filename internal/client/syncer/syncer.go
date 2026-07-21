package syncer

import (
	"gamesync/internal/client/config"
	"gamesync/internal/client/profiler"
	api "gamesync/internal/ogen"
	"os"

	"github.com/klauspost/compress/zstd"
)


type syncer struct {
	conf    *config.Config
	client  *api.Client
	profile profiler.Profile
}

type puller struct {
	client   *api.Client
	decoder  *zstd.Decoder
	file     *os.File
	fileBytesWritten int64
	chunkDir string
}

func New(conf *config.Config, c *api.Client, profile profiler.Profile) *syncer { 
	return &syncer{
		client: c,
		conf: conf,
		profile: profile,
	}
}

