package stater

import (
	"fmt"
	"testing"
)

func TestEqual(t *testing.T) {
	testData := []struct {
		name string
		in1 map[string]FileState
		in2 map[string]FileState
		out bool
	} {
		{
			name: "one file equal",
			in1: map[string]FileState{ "asdf": { Size: 3, ModTime: 1234 } },
			in2: map[string]FileState{ "asdf": { Size: 3, ModTime: 1234 } },
			out: true,
		},
		{
			name: "one file size diff",
			in1: map[string]FileState{ "asdf": { Size: 2, ModTime: 1234 } },
			in2: map[string]FileState{ "asdf": { Size: 3, ModTime: 1234 } },
			out: false,
		},
		{
			name: "one file mod diff",
			in1: map[string]FileState{ "asdf": { Size: 3, ModTime: 1234 } },
			in2: map[string]FileState{ "asdf": { Size: 3, ModTime: 1 } },
			out: false,
		},
		{
			name: "multi files equal",
			in1: map[string]FileState{
				"asdf":      { Size: 3, ModTime: 1234 },
				"asdf/asdf": { Size: 1000, ModTime: 1234 },
				"asdf/1":    { Size: 0, ModTime: 1 },
			},
			in2: map[string]FileState{
				"asdf":      { Size: 3, ModTime: 1234 },
				"asdf/asdf": { Size: 1000, ModTime: 1234 },
				"asdf/1":    { Size: 0, ModTime: 1 },
			},
			out: true,
		},
		{
			name: "multi files mod diff",
			in1: map[string]FileState{
				"asdf":      { Size: 3, ModTime: 1234 },
				"asdf/asdf": { Size: 1000, ModTime: 1234 },
				"asdf/1":    { Size: 0, ModTime: 1 },
			},
			in2: map[string]FileState{
				"asdf":      { Size: 3, ModTime: 1234 },
				"asdf/asdf": { Size: 1000, ModTime: 1234 },
				"asdf/1":    { Size: 0, ModTime: 2 },
			},
			out: false,
		},
		{
			name: "file count diff",
			in1: map[string]FileState{
				"asdf":      { Size: 3, ModTime: 1234 },
				"asdf/asdf": { Size: 1000, ModTime: 1234 },
				"asdf/1":    { Size: 0, ModTime: 1 },
			},
			in2: map[string]FileState{
				"asdf":      { Size: 3, ModTime: 1234 },
				"asdf/asdf": { Size: 1000, ModTime: 1234 },
			},
			out: false,
		},
	}

	for i, td := range testData {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if result := Equal(td.in1, td.in2); result != td.out {
				t.Errorf("Equal(\n%+v, \n%+v\n) = %v, expected %v", td.in1, td.in2, result, td.out)
			}
		})
	}
}

