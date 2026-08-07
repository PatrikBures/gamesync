package syncer

import (
	"context"
	"fmt"
	"os"

	api "go.pabu.dev/gamesync/internal/ogen"
	"go.pabu.dev/gamesync/internal/snapshoter"
)

// Syncs a profile by automatically determining if it should push/pull or 
// if it up to date or there is a conflict. 
func (s *syncer) Sync() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, err := s.client.GetBranchHead(ctx, api.GetBranchHeadParams{
		UserID: s.conf.Server.UserID,
		RepoName: s.profile.RepoName,
		BranchName: s.profile.BranchName,
	})
	if err != nil {
		return err
	}

	chunkGen := snapshoter.NewChunkGen(s.conf.Global.ChunkDir)

	var snapshotID int64
	var parentSnapshotID int64

	switch snapshot := res.(type) {
	case *api.GetBranchHeadNotFound:
		fmt.Println("no snapshot in branch")
		files, err := chunkGen.ChunkFilesInDirApiFile(ctx, s.profile.Dir)
		if err != nil {
			return err
		}
		newSnapshot, err := s.CreateSnapshot(files)
		if err != nil {
			return fmt.Errorf("creating snapshot: %w", err)
		}
		if err := s.stater.SetProfileSnapshot(s.profile.Slug, newSnapshot.SnapshotID); err != nil {
			return fmt.Errorf("setting profile snapshot state: %w", err)
		}
		return nil
	case *api.Snapshot:
		snapshotID = snapshot.SnapshotID
		if !snapshot.ParentSnapshotID.Null {
			parentSnapshotID = snapshot.ParentSnapshotID.Value
		}
	default:
		return fmt.Errorf("unexpected result type: %T", snapshot)
	}

	unknownState := false
	stateSnapshotID, err := s.stater.GetProfileSnapshot(s.profile.Slug)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("getting profile snapshot state: %w", err)
		} 
		unknownState = true
	}

	latestSnapshot, err := s.client.GetSnapshot(ctx, api.GetSnapshotParams{
		UserID: s.conf.Server.UserID,
		RepoName: s.profile.RepoName,
		BranchName: s.profile.BranchName,
		SnapshotID: snapshotID,
	})
	if err != nil {
		return err
	}

	files, err := chunkGen.ChunkFilesInDirApiFile(ctx, s.profile.Dir)
	if err != nil {
		return err
	}

	diff := false
	localFileCount := 0
	for _, fr := range files {
		if !diff {
			for _, f := range latestSnapshot.Files {
				if f.Path != fr.Path {
					continue
				}
				if f.Hash != fr.Hash {
					diff = true
					break
				}
			}
		}
		localFileCount++
	}

	if localFileCount != len(latestSnapshot.Files) {
		diff = true
	}

	if unknownState && diff {
		return fmt.Errorf("no previous local state and difference between remote and local. Force pull or push.")
	}

	if diff && snapshotID == stateSnapshotID {
		fmt.Println("Pushing...")
		_, err := s.CreateSnapshot(files)
		if err != nil {
			return fmt.Errorf("creating snapshot: %w", err)
		}
	} else if diff && parentSnapshotID == stateSnapshotID {
		fmt.Println("Pulling...")
		if err := s.Restore(snapshotID); err != nil {
			return fmt.Errorf("restoring snapshot: %w", err)
		}
	} else if diff && parentSnapshotID != stateSnapshotID && snapshotID != stateSnapshotID {
		// this will only occur if it is 2 snapshots behind. this is an issue
		// when it is only behind by one snapshot and there are changes on both sides. 
		// in which case it will just pull. 
		return fmt.Errorf("Changes on remote and local, either force push or pull.")
	} else if !diff && snapshotID == stateSnapshotID {
		fmt.Println("Already up to date")
	} else if !diff {
		// also handles case when unknownState is true
		fmt.Println("Fast forwarding...")
		if err := s.stater.SetProfileSnapshot(s.profile.Slug, snapshotID); err != nil {
			return fmt.Errorf("fast forwarding snapshot: %w", err)
		}
	} else {
		fmt.Println("diff:", diff)
		fmt.Println("snapshotID", snapshotID)
		fmt.Println("parentSnapshotID", parentSnapshotID)
		fmt.Println("stateSnapshotID", stateSnapshotID)

		return fmt.Errorf("unhandled combination")
	}

	return nil
}
