package server

import "errors"

var (
	ErrNotAuthorized      = errors.New("not authorized")
	ErrPermNotFound       = errors.New("permission not found")

	ErrRepoNotFound       = errors.New("repo not found")
	ErrRepoNameConflict   = errors.New("repo name already exists")

	ErrBranchNotFound     = errors.New("branch not found")
	ErrBranchNameConflict = errors.New("branch name already exists")

	ErrRolePermsNotFound  = errors.New("role perms not found")

	ErrRoleNameConflict   = errors.New("role name already exists")
	ErrRoleNotFound       = errors.New("role not found")

	ErrUserNotFound       = errors.New("user not found")
	ErrUserNameConflict   = errors.New("user name already exists")

	ErrTokenLength        = errors.New("token is too long or short")

	ErrBase64             = errors.New("could not decode base64")

	ErrHashMismatch       = errors.New("hashes do not match")
	ErrInvalidHash        = errors.New("invalid hash")

	ErrChunkTooBig        = errors.New("chunk is too big")
)
