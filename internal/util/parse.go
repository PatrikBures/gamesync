package util

import (
	"bytes"
)

func ParseArgs(s string) []string {
	var args []string

	var quotedDoube bool
	var quotedSingle bool

	var buffer bytes.Buffer
	for _, c := range s {
		switch c {
		case '"':
			if quotedSingle {
				buffer.WriteString(string(c))
			} else if !quotedDoube {
				quotedDoube = true
			} else {
				quotedDoube = false
				args = append(args, buffer.String())
				buffer.Reset()
			}
		case '\'':
			if quotedDoube {
				buffer.WriteString(string(c))
			} else if !quotedSingle {
				quotedSingle = true
			} else {
				quotedSingle = false
				args = append(args, buffer.String())
				buffer.Reset()
			}
		case ' ':
			if quotedDoube || quotedSingle {
				buffer.WriteString(string(c))
			} else if buffer.Len() > 0 {
				args = append(args, buffer.String())
				buffer.Reset()
			}
		default:
			buffer.WriteString(string(c))
		}
	}
	if buffer.Len() > 0 {
		args = append(args, buffer.String())
	}

	return args
}

