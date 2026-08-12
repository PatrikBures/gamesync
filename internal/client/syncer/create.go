package syncer

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/schollz/progressbar/v3"
	api "go.pabu.dev/gamesync/internal/ogen"
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
			return r, nil
		default:
			return nil, fmt.Errorf("unrecognized type %T with result: %v", r, r)
		}
	}
	return nil, fmt.Errorf("did not succeed creating snapshot, this error should not occcur")
}

func (s *syncer) uploadMissing(ctx context.Context, chunkHashes []string) error {
	var totalSize int64
	for _, ch := range chunkHashes {
		stat, err := os.Stat(s.conf.ChunkHashPath(ch))
		if err != nil {
			return err
		}
		totalSize += stat.Size()
	}

	bar := progressbar.NewOptions64(totalSize,
		progressbar.OptionSetDescription(fmt.Sprintf("[0/%d] Starting uploading chunks...", len(chunkHashes))),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionUseIECUnits(true),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() { fmt.Fprintln(os.Stderr) } ),
	)

	for i, ch := range chunkHashes {
		path := s.conf.ChunkHashPath(ch)
		f, err := os.Open(path)
		if err != nil {
			return err
		}

		br := progressbar.NewReader(f, bar)
		
		result, err := s.client.PutChunk(ctx, api.PutChunkReq{Data: &br}, api.PutChunkParams{ChunkHash: ch})
		if err != nil {
			return err
		}
		switch r := result.(type) {
		case *api.PutChunkOK, *api.PutChunkCreated:
			bar.Describe(fmt.Sprintf("[%d/%d] Uploading %s...", i+1, len(chunkHashes), ch[:8]))
		default:
			return fmt.Errorf("uploading chunk: unrecognized type %T with result: %v", r, r)
		}
	}

	if err := bar.Finish(); err != nil {
		return err
	}
	
	return nil
}
