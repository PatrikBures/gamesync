package syncer

import (
	"os"
	"path/filepath"

	"go.pabu.dev/gamesync/internal/client/fileval"
)

type stater struct {
	stateDir string
}

// make sure the provided stateDir is ProfileStateDir
func NewStater(stateDir string) *stater {
	return &stater{
		stateDir: stateDir,
	}
}

// get current profile snapshotID
//
// if state does not exist, err is os.ErrNotExist
func (s *stater) GetProfileSnapshot(profile string) (snapshotID int64, err error) {
	p := filepath.Join(s.stateDir, profile)
	err = fileval.Read(p, &snapshotID)
	return
}

func (s *stater) SetProfileSnapshot(profile string, snapshotID int64) error {
	p := filepath.Join(s.stateDir, profile)
	return fileval.Write(p, snapshotID)
}

func (s *stater) DeleteProfileSnapshot(profile string) error {
	p := filepath.Join(s.stateDir, profile)
	err := os.Remove(p)
	if os.IsNotExist(err) {
		err = nil
	}
	return err
}
