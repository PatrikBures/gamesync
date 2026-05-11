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
)

type Service struct {
	q *query.Query
}
func NewService(query *query.Query) *Service {
	return &Service{
		q: query,
	}
}
func (s *Service) GetHealth(ctx context.Context) error {
	return nil
}


