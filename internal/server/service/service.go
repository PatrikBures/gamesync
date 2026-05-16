package service

import (
	"context"
	"database/sql/driver"
	"gamesync/internal/query"
)

const tokenLen = 33

// This type is requred as the gorm driver does not support inserting bytes
type HashBytes []byte
func (h HashBytes) Value() (driver.Value, error) {
	return []byte(h), nil
}

type contextKey int
const (
	userContextKey contextKey = iota
	userPermsContextKey
)

type Service struct {
	q *query.Query
	o ServiceOpts
}
type ServiceOpts struct {
	DefaultRoleID int32
}
func NewService(query *query.Query, opts ServiceOpts) *Service {
	return &Service{
		q: query,
		o: opts,
	}
}
func (s *Service) GetHealth(ctx context.Context) error {
	return nil
}


