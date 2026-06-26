package main

import (
	"errors"
	"fmt"
	api "gamesync/internal/ogen"
)

func errHandler(err error) error {
	if gErr, ok := errors.AsType[*api.GlobalErrorStatusCode](err); ok {
		return fmt.Errorf("server returned error (%d): %s", gErr.Response.Code, gErr.Response.Message)
	}
	return err
}
