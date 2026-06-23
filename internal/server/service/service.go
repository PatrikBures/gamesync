package service

import (
	"context"
	"gamesync/internal/dbx"
	"gamesync/internal/server/storage"
)

const tokenLen = 33

type ContextKey int
const (
	CkUser ContextKey = iota
	CkRolePerms
	CkRepoID
	CkBranchID
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
