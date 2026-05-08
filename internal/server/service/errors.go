package service

import "errors"

var (
	ErrMissingBody      = errors.New("request body is required")
	ErrDatabase         = errors.New("database error")
	ErrDuplicateKey     = errors.New("unique constraint failed")
	ErrNotFound         = errors.New("resource not found")
)
