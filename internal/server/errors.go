package server

import "errors"

var (
	ErrMissingBody   = errors.New("request body is required")
	ErrDatabase      = errors.New("database error")
	ErrDuplicateKey  = errors.New("unique constraint failed")
	ErrNotFound      = errors.New("resource not found")
	ErrToken         = errors.New("token error")
	ErrAuth          = errors.New("authentication error")
	ErrNotAuthorized = errors.New("not authorized")
	ErrPermNotFound  = errors.New("permission not found")
	ErrContext       = errors.New("context error")
)
