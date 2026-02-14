package ui

import (
	"fmt"
	"io"
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

var OutWriter io.Writer = os.Stdout
var ErrWriter io.Writer = os.Stderr

func SetLevel(l Level) {
	currentLevel = l
}

func GetLevel() Level {
	return currentLevel
}

func Printf(l Level, format string, args ...any) {
	if currentLevel <= l {
		var err error
		if l == LevelError {
			_, err = fmt.Fprintf(ErrWriter, format, args...)
		} else {
			_, err = fmt.Fprintf(OutWriter, format, args...)
		}

		if err != nil {
			panic("Error printing")
		}
	}
}

func Info(format string, args ...any) {
	Printf(LevelNormal, format, args...)
}
func Error(format string, args ...any) {
	Printf(LevelError, format, args...)
}
func Verbose(format string, args ...any) {
	Printf(LevelVerbose, format, args...)
}
func Debug(format string, args ...any) {
	Printf(LevelDebug, format, args...)
}

