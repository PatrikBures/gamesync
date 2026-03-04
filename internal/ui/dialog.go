package ui

import (
	"fmt"
	"time"

	"github.com/ncruces/zenity"
)


type ConflictType int
const (
	ConflictCancel ConflictType = iota
	ConflictPull
	ConflictPush
	ConflictIgnore
	ConflictError
)

type conflict struct {
	conflictType ConflictType
	msg string
}

func DialogConflict(gameID string, dateTimeRemote int64, dateTimeLocal int64) (ConflictType, error) {

	conflicts := []conflict{
		{ conflictType: ConflictCancel, msg: "Cancel launch" },
		{ conflictType: ConflictPull,   msg: fmt.Sprintf("Force pull, latest remote %s", time.Unix(dateTimeRemote, 0).Format("2006-01-02 15:04:05")) },
		{ conflictType: ConflictPush,   msg: fmt.Sprintf("Force push, latest local %s",  time.Unix(dateTimeLocal,  0).Format("2006-01-02 15:04:05")) },
		{ conflictType: ConflictIgnore, msg: "Launch game, ignoring conflict (Not recommended)" },
	}

	opts := []string{}
	for _, c := range conflicts {
		opts = append(opts, c.msg)
	}

	pickedOpt, err := zenity.List(
		"A conflict occured while syncing GAME_ID, what action do you want to take?",
		opts,
		zenity.DisallowEmpty(),
	)

	if err != nil {
		return ConflictCancel, err
	}

	pickedType := ConflictError
	for i, o := range opts {
		if o == pickedOpt {
			pickedType = conflicts[i].conflictType
			break
		}
	}
	if pickedType == ConflictError {
		return ConflictCancel, fmt.Errorf("Invalid conflict option: %s", pickedOpt)
	}

	Debug("picked: %s, %d\n", pickedOpt, pickedType)

	return pickedType, nil
}
