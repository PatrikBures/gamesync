package ui

import (
	"github.com/gen2brain/beeep"
)

func Notify(title string, message string) {
	beeep.AppName = "gamesync"

	err := beeep.Notify(title, message, "")

	if err != nil {
		Error("Error notifying: %s\n", err)
	}
}

