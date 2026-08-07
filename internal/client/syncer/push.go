package syncer

import (
	"context"
	"errors"
	"fmt"
	api "go.pabu.dev/gamesync/internal/ogen"
	"go.pabu.dev/gamesync/internal/snapshoter"
	"os"
	"path/filepath"
)

// Creates a snapshot. If it receives a list of missing chunks, it will
// upload them and try to create the snapshot again
//
// It expects the chunks to already exist localy
func (s *syncer) CreateSnapshot(files []api.File) (*api.Snapshot, error) {
	request := api.Files{
		Files: files,
	}

	params := api.PostSnapshotParams{
		UserID:     s.conf.Server.UserID,
		RepoName:   s.profile.RepoName,
		BranchName: s.profile.BranchName,
	}

	for range 2 {
		result, err := s.client.PostSnapshot(context.Background(), &request, params)
		if err != nil {
			// this is a duplicate of a function in util.go in client cmd
			if gErr, ok := errors.AsType[*api.GlobalErrorStatusCode](err); ok {
				return nil, fmt.Errorf("server returned error (%d): %s", gErr.Response.Code, gErr.Response.Message)
			}
			return nil, fmt.Errorf("posting snapshot: %w", err)
		}
		switch r := result.(type) {
		case *api.PostSnapshotFailedDependency:
			if err := s.uploadMissing(context.Background(), r.ChunkHashes); err != nil {
				return nil, err
			}
		case *api.Snapshot:
			fmt.Printf("Created snapshot for '%s' repo on branch '%s'\n", s.profile.RepoName, s.profile.BranchName)

			if err := s.stater.SetProfileSnapshot(s.profile.Slug, r.SnapshotID); err != nil {
				return nil, fmt.Errorf("setting profile snapshot: %w", err)
			}
			return r, nil
		default:
			return nil, fmt.Errorf("unrecognized type %T with result: %v", r, r)
		}
	}
	return nil, fmt.Errorf("did not succeed creating snapshot, this error should not occcur")
}

func (s *syncer) uploadMissing(ctx context.Context, chunkHashes []string) error {
	for i, ch := range chunkHashes {
		path := filepath.Join(s.conf.Global.ChunkDir, snapshoter.DirsForChunk(ch), ch)
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		result, err := s.client.PutChunk(ctx, api.PutChunkReq{Data: f}, api.PutChunkParams{ChunkHash: ch})
		if err != nil {
			return err
		}
		switch r := result.(type) {
		case *api.PutChunkOK, *api.PutChunkCreated:
			fmt.Printf("uploaded chunk %d/%d\n", i+1, len(chunkHashes))
		default:
			return fmt.Errorf("uploading chunk: unrecognized type %T with result: %v", r, r)
		}
	}
	return nil
}
