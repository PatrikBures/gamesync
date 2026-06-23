package snapshoter

import (
	"context"
	"fmt"
	api "gamesync/internal/ogen"
	"os"
	"path/filepath"
)


type uploader struct {
	client   *api.Client
	params   api.PostUserRepoBranchSnapshotParams
	chunkDir string
	attempts int
}

type UploaderOpts struct {
}

func NewUploader(client *api.Client, params api.PostUserRepoBranchSnapshotParams, chunkDir string, attempts int) uploader {
	if attempts < 1 {
		attempts = 1
	}
	return uploader{
		client: client,
		params: params,
		chunkDir: chunkDir,
		attempts: attempts,
	}
}


// attempts to create snapshot 'attempts' amount of times with the provided 'files'. 
// if the issue is that there are missing chunks on server, it will upload them and 
// try again if there are more than one 'attempts'
func (u *uploader) CreateSnapshot(files []FileResults) error {

	request := api.SnapshotNew{}
	// make capacity big enough for all files
	request.Files = make([]api.File, 0, len(files))
	for _, f := range files {
		request.Files = append(request.Files, api.File{
			Path: f.Path,
			Hash: f.Hash,
			ChunkHashes: f.Hashes,
		})
	}
	for range 2 {
		result, err := u.client.PostUserRepoBranchSnapshot(context.Background(), &request, u.params)
		if err != nil {
			return err
		}
		switch r := result.(type) {
		case *api.PostUserRepoBranchSnapshotFailedDependency:
			if err := u.uploadMissing(context.Background(), r.ChunkHashes); err != nil {
				return err
			}
		case *api.PostUserRepoBranchSnapshotCreated:
			fmt.Printf("Created snapshot for '%s' repo on branch '%s'\n", u.params.RepoName, u.params.BranchName)
			return nil
		default:
			return fmt.Errorf("unrecognized type %T with result: %v", r, r)
		}
	}
	return nil
}

func (s *uploader) uploadMissing(ctx context.Context, chunkHashes []string) error {
	for i, ch := range chunkHashes {
		path := filepath.Join(s.chunkDir, DirsForChunk(2, 2, ch), ch)
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
