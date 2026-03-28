package dbm

import (
	"fmt"
	"strings"
)

func listSQL(qty int) string {
	return "(" + strings.Repeat("?,", qty-1) + "?" + ")"
}

func FlattenArgs[T any](cols ...[]T) []any {
	if len(cols) == 0 {
		return nil
	}
	nRows := len(cols[0])
	for i, col := range cols {
		if len(col) != nRows {
			panic(fmt.Errorf("column %d length missmatch: got %d, expected %d", i, len(col), nRows))
		}
	}
	args := make([]any, 0, nRows*len(cols))
	for i := range nRows {
		for _, col := range cols {
			args = append(args, col[i])
		}
	}
	return args
}
