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

type contextKey int
const (
	ckUser contextKey = iota
	ckRolePerms
)

type Service struct {
	conn dbx.DBconn
	o ServiceOpts
}
type ServiceOpts struct {
	DefaultRoleID int32
}
func NewService(conn dbx.DBconn, opts ServiceOpts) *Service {
	return &Service{
		conn: conn,
		o: opts,
	}
}
func (s *Service) GetHealth(ctx context.Context) error {
	return nil
}


