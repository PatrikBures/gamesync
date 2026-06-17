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
	ErrRepoNotFound  = errors.New("repo not found")
	ErrContext       = errors.New("context error")
	ErrDecompress    = errors.New("decompressing error")
	ErrHashing       = errors.New("hashing error")
	ErrFile          = errors.New("file operation failed")
	ErrFileClose     = errors.New("closing file failed")
	ErrFileTooBig    = errors.New("file is too big")
	ErrHashMismatch = errors.New("hashes do not match")
)
