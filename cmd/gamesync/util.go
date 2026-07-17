package main

import (
	"errors"
	"fmt"
	api "gamesync/internal/ogen"

	"github.com/spf13/cobra"
)

func errHandler(err error) error {
	if gErr, ok := errors.AsType[*api.GlobalErrorStatusCode](err); ok {
		return fmt.Errorf("server returned error (%d): %s", gErr.Response.Code, gErr.Response.Message)
	}
	return err
}

func markFlagsRequired(cmd *cobra.Command, names []string) error {
	for _, n := range names {
		if err := cmd.MarkFlagRequired(n); err != nil {
			return err
		}
	}
	return nil
}
