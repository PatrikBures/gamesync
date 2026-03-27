package dbm

import "strings"

func listSQL(qty int) string {
	return "(" + strings.Repeat("?,", qty-1) + "?" + ")"
}
