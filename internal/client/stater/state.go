package stater

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type Stater struct {
	stateDir string
}

type State struct {
	SnapshotID int64
	FileStates map[string]FileState
}

type FileState struct {
	Size int64
	ModTime int64
}

// make sure the provided stateDir is ProfileStateDir
func New(stateDir string) *Stater {
	return &Stater{
		stateDir: stateDir,
	}
}

func (s *Stater) resolvePath(profile string) string {
	return filepath.Join(s.stateDir, filepath.Base(profile) + ".json")
}

// reads current profile state. returns  os.ErrNotExist if missing.
func (s *Stater) Get(profile string) (State, error) {
	file, err := os.Open(s.resolvePath(profile))
	if err != nil {
		return State{}, err
	}

	defer func() {
		_ = file.Close()
	}()

	var state State
	if err = json.NewDecoder(file).Decode(&state); err != nil {
		return State{}, err
	}

	return state, nil
}

// Atomically writes state
func (s *Stater) Set(profile string, state State) error {
	tempFile, err := os.CreateTemp(s.stateDir, ".tmp-*")
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tempFile.Close()
			_ = os.Remove(tempFile.Name())
		}
	}()

	encoder := json.NewEncoder(tempFile)
	encoder.SetIndent("", "    ")

	if err = encoder.Encode(state); err != nil {
		return err
	}

	if err = tempFile.Close(); err != nil {
		return err
	}

	return os.Rename(tempFile.Name(), s.resolvePath(profile))
}

func (s *Stater) Delete(profile string) error {
	err := os.Remove(s.resolvePath(profile))
	if errors.Is(err, os.ErrNotExist) {
		err = nil
	}
	return err
}

func (s *Stater) Update(profile string, snapshotID int64, dir string) error {
	state := State{
		SnapshotID: snapshotID,
		FileStates: make(map[string]FileState),
	}

	if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		fileState, err := StateFile(path)
		if err != nil {
			return err
		}

		state.FileStates[relPath] = fileState

		return nil
	}); err != nil {
		return err
	}

	return s.Set(profile, state)
}

func StateFile(path string) (FileState, error) {
	stat, err :=  os.Stat(path)
	if err != nil {
		return FileState{}, err
	}
	return FileState{
		ModTime: stat.ModTime().Unix(),
		Size: stat.Size(),
	}, nil
}


func Equal(a, b map[string]FileState) bool {
	if len(a) != len(b) {
		return false
	}

	for k, fsa := range a {
		fsb, ok := b[k]
		if !ok || fsa != fsb {
			return false
		}
	}

	return true
}

func CopyToSimpleMap(dir string, state map[string]FileState) map[string]struct{} {
	simple := make(map[string]struct{}, len(state))
	for f := range state {
		simple[filepath.Join(dir, f)] = struct{}{}
	}
	return simple
}
