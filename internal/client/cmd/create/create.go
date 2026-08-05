package create

import (
	"go.pabu.dev/gamesync/internal/client/config"

	"github.com/spf13/cobra"
)

type createCmd struct {
	Cmd *cobra.Command
}

func New(conf *config.Config) *createCmd {
	root := createCmd{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create resource",
	}
	cmd.AddCommand(
		newRepoCmd(conf).cmd,
		newSnapshotCmd(conf).cmd,
		newProfileCmd(conf).cmd,
		newBranchCmd(conf).cmd,
	)

	root.Cmd = cmd
	return &root
}
