package config

import (
	"fmt"
	"path/filepath"

	"go.pabu.dev/packer"
	"go.pabu.dev/packer/data/local"
	"go.pabu.dev/packer/metadata/bolt"
)

func (conf *Config) InitPackager() (*packer.Packager, error) {
	ms := metadata.New(filepath.Join(conf.Global.ChunkDir, "metadata.db"))
	ds := data.New(conf.Global.ChunkDir)

	p := packer.New(ms, ds)

	if err := p.Init(); err != nil {
		return nil, fmt.Errorf("initializing packager: %w", err)
	}

	return p, nil
}
