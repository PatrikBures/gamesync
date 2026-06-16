package service

import (
	"context"
	"database/sql/driver"
	"gamesync/internal/dbx"
)

const tokenLen = 33

// This type is requred as the gorm driver does not support inserting bytes
type HashBytes []byte

func (h HashBytes) Value() (driver.Value, error) {
	return []byte(h), nil
}

type ContextKey int
const (
	CkUser ContextKey = iota
	CkRolePerms
)

type Service struct {
	db *dbx.DB
	o  ServiceOpts
}
type ServiceOpts struct {
	DefaultRoleID int32
}

func NewService(db *dbx.DB, opts ServiceOpts) *Service {
	return &Service{
		db: db,
		o:  opts,
	}
}
func (s *Service) GetHealth(ctx context.Context) error {
	return nil
}
