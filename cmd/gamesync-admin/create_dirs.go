package main

import (
	"fmt"
	"gamesync/internal/vars"
	"os"

	"github.com/spf13/cobra"
)

type createDirsCmd struct {
	cmd *cobra.Command
}

func newCreateDirsCmd() *createDirsCmd {
	root := createDirsCmd{}
	cmd := &cobra.Command{
		Use: "init-dirs",
		Short: "Creates the necessary dirs for the server, moves old paths",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, path := range vars.RemoteSaveDirOld {
				f, err := os.Stat(path)
				if err != nil || !f.IsDir(){
					continue
				}
				os.Rename(path, vars.RemoteSaveDir)
				fmt.Printf("moved save dir from %s to %s\n", path, vars.RemoteSaveDir)
				break
			}

			dirs := []string{
				vars.RemoteSaveDir,
				vars.RemoteBackupDir,
			}

			for _, dir := range dirs {
				if err := os.MkdirAll(dir, 0775); err != nil {
					return err
				}
				fmt.Println("Ensured dir exists:", dir)
			}

			return nil
		},
	}
	root.cmd = cmd
	return &root
}
