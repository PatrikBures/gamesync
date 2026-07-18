package main

import (
	"context"
	"errors"
	"fmt"
	"gamesync/internal/client"
	"gamesync/internal/client/config"
	api "gamesync/internal/ogen"
	"gamesync/internal/snapshoter"
	"os"
	"path/filepath"
	"strconv"

	"github.com/klauspost/compress/zstd"
	"github.com/spf13/cobra"
)


type restoreCmd struct {
	cmd *cobra.Command
	opts restoreOpts
}
type restoreOpts struct {
	branch string
	repo string
	dir string
	snapshotID int64
}

func newRestoreCmd(conf *config.Config) *restoreCmd {
	root := restoreCmd{}

	cmd := &cobra.Command{
		Use:   "restore DIR SNAPSHOT_ID",
		Short: "Restore a sync to a specific snapshot",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := populateRestoreOpts(&root.opts, args); err != nil { return err }

			c, err := client.Client(conf)
			if err != nil { return err }

			if err := runRestoreCmd(c, root.opts, *conf); err != nil { return err }

			return nil
		},
	}

	cmd.Flags().StringVarP(&root.opts.repo,   "repo", "r", "", "")
	cmd.Flags().StringVarP(&root.opts.branch, "branch", "b", "", "")

	if err := markFlagsRequired(cmd, []string{"repo", "branch"}); err != nil {
		panic(err)
	}

	root.cmd = cmd
	return &root
}

func populateRestoreOpts(opts *restoreOpts, args []string) error {
	opts.dir = args[0]

	d, err := os.Stat(opts.dir)
	if err != nil { return err }
	if !d.IsDir() {
		return fmt.Errorf("provided dir is not a dir: %s", opts.dir)
	}

	opts.snapshotID, err = strconv.ParseInt(args[1], 10, 64)
	if err != nil { return err }

	return nil
}

func runRestoreCmd(c *api.Client, opts restoreOpts, conf config.Config) (err error) {
	snapshot, err := c.GetSnapshot(context.Background(), api.GetSnapshotParams{
		UserID: conf.Server.UserID,
		RepoName: opts.repo,
		BranchName: opts.branch,
		SnapshotID: opts.snapshotID,
	})
	if err != nil {
		return fmt.Errorf("getting snapshot: %w", err)
	}

	// if it errors anywhere in the loop. then the sync will not be up to date. if it errors there
	// should be a defer which restores its state the the previous snapshot.
	for _, f := range snapshot.Files {
		// currently only downloads the chunks. needs to also create the files from the chunks
		// when checking if the files exists, it will just delete it and write a new one for now. 
		// until there is a state which stores the current chunk hash, size and modtime, which 
		// can be used to quickly check if it has changed. 
		//
		// when it writes the chunk, it needs to write the chunk and also append to the file. 
		// if the chunk already exists, it just needs to append to the file.

		filePath := filepath.Join(opts.dir, f.Path)

		fmt.Printf("\nwriting to: %s\n", filePath)

		dir := filepath.Base(filePath)
		if err = os.MkdirAll(dir, 0775); err != nil {
			return fmt.Errorf("failed creating dir: %w", err)
		}


		decoder, err := zstd.NewReader(nil)
		if err != nil {
			return fmt.Errorf("initializing zstd decoder: %w", err)
		}
		for _, ch := range f.ChunkHashes {
			chunk, err := c.GetChunk(context.Background(), api.GetChunkParams{ChunkHash: ch})
			if err != nil {
				return fmt.Errorf("failed getting chunk with hash '%s': %w", ch, err)
			}

			chunkFile, _, err := snapshoter.CreateChunkFile(conf.Global.ChunkDir, ch)
			if err != nil {
				if errors.Is(err, os.ErrExist) {
					fmt.Println("skipped:", ch)
					continue
				}
				return err
			}
			defer func() {
				err = errors.Join(err, chunkFile.Close())
			}()

			// need to verify that the hash matches the content
			if err := decoder.Reset(chunk.Data); err != nil {
				return fmt.Errorf("reseting decoder while downloading chunk '%s': %w", ch, err)
			}
			if i, err := decoder.WriteTo(chunkFile); err != nil {
				return fmt.Errorf("decoding chunk '%s', wrote %d bytes: %w", ch, i, err)
			}
			fmt.Println("Downloaded:", ch)
		}
	}

	return nil
}
