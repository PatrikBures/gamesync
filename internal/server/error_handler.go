package server

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"
)


func ErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	var code = http.StatusInternalServerError
	var ogenErr ogenerrors.Error
	switch {
	case errors.Is(err, ErrRepoNotFound):
		code = http.StatusNotFound
	case errors.Is(err, ErrBranchNotFound):
		code = http.StatusNotFound
	case errors.As(err, &ogenErr):
		code = ogenErr.Code()
	}
	w.WriteHeader(code)
	_, _ = io.WriteString(w, http.StatusText(code))
}
