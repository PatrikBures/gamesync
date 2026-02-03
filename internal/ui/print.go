package ui

import (
	"fmt"
	"os"
)

type Level int

const (
	LevelDebug Level = iota
	LevelVerbose
	LevelNormal
	LevelError 
)

var currentLevel = LevelNormal

func SetLevel(l Level) {
	currentLevel = l
}

func Printf(l Level, format string, args ...any) {
	if currentLevel <= l {
		if l == LevelError {
			fmt.Fprintf(os.Stderr, format, args...)
			return
		}

		fmt.Printf(format, args...)
	}
}
