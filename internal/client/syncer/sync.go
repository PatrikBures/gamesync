package syncer

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.pabu.dev/gamesync/internal/client/stater"
	api "go.pabu.dev/gamesync/internal/ogen"
	"go.pabu.dev/gamesync/internal/snapshoter"
)

type SyncMode int
const (
	ModeAuto SyncMode = iota
	ModePush
	ModePushForce
	ModePull
	ModePullForce
)

var (
	ErrNoExistingSnapshot = errors.New("no snapshot on branch, nothing to pull")
)

// Syncs a profile by automatically determining if it should push/pull or
// if it up to date or there is a conflict.
func (s *syncer) Sync(mode SyncMode) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	noHead, snapshotID, parentSnapshotID, err := s.head(ctx)
	if err != nil {
		return err
	}

	if noHead && (mode == ModePull || mode == ModePullForce) {
		return ErrNoExistingSnapshot
	}

	unknownState := false
	previousState, err := s.stater.Get(s.profile.Slug)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("getting profile snapshot state: %w", err)
		} 
		unknownState = true
	}

	currentFileStates, files, err := chunkFiles(ctx, s.conf.Global.ChunkDir, s.profile.Dir)
	if err != nil {
		return err
	}


	if noHead {
		fmt.Println("Pushing with no head...")
		newSnapshot, err := s.CreateSnapshot(files)
		if err != nil {
			return fmt.Errorf("creating snapshot: %w", err)
		}
		return s.stater.Set(s.profile.Slug, stater.State{
			SnapshotID: newSnapshot.SnapshotID,
			FileStates: currentFileStates,
		})
	}

	localDiff := !stater.Equal(currentFileStates, previousState.FileStates)
	remoteDiff := snapshotID != previousState.SnapshotID

	if (localDiff && !remoteDiff && mode != ModePullForce) || mode == ModePushForce {
		if mode == ModePull {
			return fmt.Errorf("local files are newer than remote, either push or force pull")
		}
		fmt.Println("Pushing...")
		newSnapshot, err := s.CreateSnapshot(files)
		if err != nil {
			return fmt.Errorf("creating snapshot: %w", err)
		}
		return s.stater.Set(s.profile.Slug, stater.State{
			SnapshotID: newSnapshot.SnapshotID,
			FileStates: currentFileStates,
		})
	} else if (!localDiff && remoteDiff) || mode == ModePullForce {
		if mode == ModePush {
			return fmt.Errorf("remote is newer than local, either pull or force push")
		}
		fmt.Println("Pulling...")
		if err := s.Restore(snapshotID, stater.CopyToSimpleMap(s.profile.Dir, currentFileStates)); err != nil {
			return fmt.Errorf("restoring snapshot: %w", err)
		}
		return s.stater.Update(s.profile.Slug, snapshotID, s.profile.Dir)
	} else if localDiff && remoteDiff {
		if unknownState {
			return fmt.Errorf("no previous local state. Force pull or push")
		}
		return fmt.Errorf("changes on remote and local, either force push or pull")
	} else if !localDiff && !remoteDiff{
		fmt.Println("Already up to date")
	} else {
		fmt.Println("localDiff:", localDiff)
		fmt.Println("remoteDiff:", remoteDiff)
		fmt.Println("snapshotID", snapshotID)
		fmt.Println("parentSnapshotID", parentSnapshotID)
		fmt.Println("previous SnapshotID", previousState.SnapshotID)

		return fmt.Errorf("unhandled combination")
	}

	return nil
}


func (s *syncer) head(ctx context.Context) (bool, int64, int64, error) {
	branches, err := s.client.GetBranches(ctx, api.GetBranchesParams{
		UserID: s.conf.Server.UserID,
		RepoName: s.profile.RepoName,
		BranchName: api.OptString{
			Value: s.profile.BranchName,
			Set: true,
		},
	})
	if err != nil {
		return false, 0, 0, err
	}

	if len(branches) == 0 {
		return true, 0, 0, nil
	}

	branch := branches[0]

	var snapshotID int64
	if branch.ParentSnapshotID.Null {
		snapshotID = -1
	} else {
		snapshotID = branch.ParentSnapshotID.Value
	}

	return false, branch.HeadSnapshotID, snapshotID, nil
}


func chunkFiles(ctx context.Context, chunkDir string, profileDir string) (map[string]stater.FileState, []api.File, error) {
	chunkGen := snapshoter.NewChunkGen(chunkDir)
	stream, err := chunkGen.ChunkFilesInDir(ctx, profileDir)
	if err != nil {
		return nil, nil, err
	}

	currentFileStates := make(map[string]stater.FileState)
	var files []api.File
	for fr := range stream.Ch {
		if fr.Err != nil {
			return nil, nil, fr.Err
		}
		files = append(files, api.File{
			ChunkHashes: fr.Hashes,
			Hash: fr.Hash,
			Path: fr.Path,
		})
		currentFileStates[fr.Path] = fr.State
		chunkGen.ProcessedFile()
	}

	fmt.Println("New created chunks:", chunkGen.Info.Created(), ", Already existing chunks:", chunkGen.Info.Skipped())

	return currentFileStates, files, nil
}
