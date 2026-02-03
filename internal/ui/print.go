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

func GetLevel() Level {
	return currentLevel
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

