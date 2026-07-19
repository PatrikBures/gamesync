package syncer

import (
	"gamesync/internal/client/config"
	api "gamesync/internal/ogen"
	"os"

	"github.com/klauspost/compress/zstd"
)


type syncer struct {
	conf    *config.Config
	client  *api.Client
	repoDir string
	repo    string
	branch  string
}

type puller struct {
	client   *api.Client
	decoder  *zstd.Decoder
	file     *os.File
	fileBytesWritten int64
	chunkDir string
}

func New(conf *config.Config, c *api.Client, repo string, branch string, repoDir string) *syncer { 
	return &syncer{
		client: c,
		conf: conf,
		repoDir: repoDir,
		repo: repo,
		branch: branch,
	}
}

