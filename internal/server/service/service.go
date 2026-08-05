package service

import (
	"context"
	"go.pabu.dev/gamesync/internal/dbx"
	"go.pabu.dev/gamesync/internal/server/storage"
)

const tokenLen = 33

type ContextKey int
const (
	CkUser ContextKey = iota
	CkRolePerms
	CkRepoID
	CkBranch
	CkSnapshot
)

type Service struct {
	db *dbx.DB
	storage storage.Chunk
	o ServiceOpts
}
type ServiceOpts struct {
	DefaultRoleID int32
}

func NewService(db *dbx.DB, storage storage.Chunk, opts ServiceOpts) *Service {
	return &Service{
		db: db,
		storage: storage,
		o:  opts,
	}
}
func (s *Service) GetHealth(ctx context.Context) error {
	return nil
}
