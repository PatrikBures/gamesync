package util

import (
	"errors"
	"fmt"
	api "go.pabu.dev/gamesync/internal/ogen"

	"github.com/spf13/cobra"
)

func ErrHandler(err error) error {
	if gErr, ok := errors.AsType[*api.GlobalErrorStatusCode](err); ok {
		return fmt.Errorf("server returned error (%d): %s", gErr.Response.Code, gErr.Response.Message)
	}
	return err
}

func MarkFlagsRequired(cmd *cobra.Command, names []string) {
	for _, n := range names {
		if err := cmd.MarkFlagRequired(n); err != nil {
			panic(fmt.Errorf("Error marking flag: %w", err))
		}
	}
}
