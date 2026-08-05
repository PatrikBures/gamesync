package delete

import (
	"go.pabu.dev/gamesync/internal/client/config"

	"github.com/spf13/cobra"
)

type deleteCmd struct {
	Cmd *cobra.Command
}

func New(conf *config.Config) *deleteCmd {
	root := deleteCmd{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete resources",
	}

	cmd.AddCommand(
		newProfileCmd(conf).cmd,
		newBranchCmd(conf).cmd,
		newRepoCmd(conf).cmd,
	)

	root.Cmd = cmd
	return &root
}
