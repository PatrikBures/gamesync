package main

import (
	"fmt"
	"os"

	createCmd "go.pabu.dev/gamesync/internal/client/cmd/create"
	deleteCmd "go.pabu.dev/gamesync/internal/client/cmd/delete"
	docsCmd "go.pabu.dev/gamesync/internal/client/cmd/docs"
	getCmd "go.pabu.dev/gamesync/internal/client/cmd/get"
	restoreCmd "go.pabu.dev/gamesync/internal/client/cmd/restore"
	syncCmd "go.pabu.dev/gamesync/internal/client/cmd/sync"
	pushCmd "go.pabu.dev/gamesync/internal/client/cmd/push"
	pullCmd "go.pabu.dev/gamesync/internal/client/cmd/pull"
	wrapCmd "go.pabu.dev/gamesync/internal/client/cmd/wrap"
	"go.pabu.dev/gamesync/internal/client/config"

	"github.com/spf13/cobra"
)

type rootCmd struct {
	cmd    *cobra.Command
	opts   rootOpts
	config config.Config
}

type rootOpts struct {
	configPath string
}

func newRootCmd() *rootCmd {
	root := rootCmd{}

	cmd := &cobra.Command{
		Use:   "gamesync",
		Short: "Syncs save games to a server",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := config.Load(&root.config, root.opts.configPath); err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&root.opts.configPath, "config", "", "Config file path")
	cmd.PersistentFlags().StringVar(&root.config.Global.ChunkDir, "chunk-dir", "", "Dir where chunks are stored")
	cmd.PersistentFlags().StringVar(&root.config.Server.Token, "token", "", "Token used to authenticate with server")
	cmd.PersistentFlags().StringVar(&root.config.Server.Url, "url", "", "Server url")
	cmd.PersistentFlags().Int64Var(&root.config.Server.UserID, "userid", 0, "User id")

	cmd.AddCommand(
		docsCmd.New().Cmd,
		getCmd.New(&root.config).Cmd,
		createCmd.New(&root.config).Cmd,
		restoreCmd.New(&root.config).Cmd,
		deleteCmd.New(&root.config).Cmd,
		syncCmd.New(&root.config).Cmd,
		pushCmd.New(&root.config).Cmd,
		pullCmd.New(&root.config).Cmd,
		wrapCmd.New(&root.config).Cmd,
	)

	cmd.DisableAutoGenTag = true
	cmd.SilenceUsage = true

	root.cmd = cmd

	return &root
}

func Execute() {
	if err := newRootCmd().cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
