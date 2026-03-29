package util

import (
	"bytes"
)

func ParseArgs(s string) []string {
	var args []string

	var quotedDouble bool
	var quotedSingle bool

	var buffer bytes.Buffer
	for _, c := range s {
		switch c {
		case '"':
			if quotedSingle {
				buffer.WriteString(string(c))
			} else if !quotedDouble {
				quotedDouble = true
			} else {
				quotedDouble = false
				args = append(args, buffer.String())
				buffer.Reset()
			}
		case '\'':
			if quotedDouble {
				buffer.WriteString(string(c))
			} else if !quotedSingle {
				quotedSingle = true
			} else {
				quotedSingle = false
				args = append(args, buffer.String())
				buffer.Reset()
			}
		case ' ', '\t', '\n':
			if quotedDouble || quotedSingle {
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

