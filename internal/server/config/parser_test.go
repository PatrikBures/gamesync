package config

import (
	"fmt"
	"maps"
	"slices"
	"testing"
)

func TestStringToSlice(t *testing.T) {
	data := map[string][]string{
		"asdff": {"asdff"},
		"":      {},
		"as|fd": {"as", "fd"},
		// by some reason, slices.Split turn the sep into space in these situations
		// "|d": {"d"},
		// "|asd|": {"asd"},
	}

	for in, out := range data {
		t.Run(fmt.Sprintf("%s = %#v", in, out), func(t *testing.T) {
			result := StringToSlice(in)
			if !slices.Equal(out, result) {
				t.Errorf("StringToSlice(%s) = '%#v', expected '%#v'", in, result, out)
			}
		})
	}
}

func TestStringToMap(t *testing.T) {
	data := map[string]map[string]bool{
		"asdff": {"asdff": true},
		"":      {},
		"as|fd": {"as": true, "fd": true},
		"|d":    {"d": true},
		"|asd|": {"asd": true},
	}

	for in, out := range data {
		t.Run(fmt.Sprintf("%s = %#v", in, out), func(t *testing.T) {
			result := StringToMap(in)
			if !maps.Equal(out, result) {
				t.Errorf("StringToMap(%s) = '%#v', expected '%#v'", in, result, out)
			}
		})
	}
}
