package service

import (
	"context"
	"errors"
	api "gamesync/internal/ogen"
	"gamesync/internal/server"
	"log/slog"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"
)

func (s *Service) NewError(ctx context.Context, err error) *api.GlobalErrorStatusCode {
	code := http.StatusInternalServerError
	msg  := "An internal server error occurred"
	
	// ogen auth wraps the error as a different type so we need to unwrap it
	if oErr, ok := errors.AsType[*ogenerrors.SecurityError](err); ok {
		err = oErr.Unwrap()
	}

	// err.Error() can probably be used for the message
	switch err {
	case server.ErrNotAuthorized:
		code = http.StatusUnauthorized
		msg = "Not authorized"

	case server.ErrUserNotFound:
		code = http.StatusNotFound
		msg = "User not found"
	case server.ErrRepoNotFound:
		code = http.StatusNotFound
		msg = "Repo not found"
	case server.ErrBranchNotFound:
		code = http.StatusNotFound
		msg = "Branch not found"
	case server.ErrRolePermsNotFound:
		code = http.StatusNotFound
		msg = "Role perms not found"
	case server.ErrRoleNotFound:
		code = http.StatusNotFound
		msg = "Role not found"

	case server.ErrUserNameConflict:
		code = http.StatusConflict
		msg = "User name already exists"
	case server.ErrRepoNameConflict:
		code = http.StatusConflict
		msg = "Repo name already exists"
	case server.ErrBranchNameConflict:
		code = http.StatusConflict
		msg = "Branch name already exists in repo"
	case server.ErrRoleNameConflict:
		code = http.StatusConflict
		msg = "Role name already exist"
	
	case server.ErrInvalidHash:
		code = http.StatusUnprocessableEntity
	}

	if code == http.StatusInternalServerError {
		var ie *server.InternalError
		if errors.As(err, &ie) {
			logArgs := append(ie.Args(), "error", ie.Unwrap())
			slog.Error(ie.Msg(), logArgs...)
		} else {
			slog.Error("internal issue", "error", err)
		}
	}

	return &api.GlobalErrorStatusCode{
		StatusCode: code,
		Response: api.Error{
			Code: int32(code),
			Message: msg,
		},
	}
}
