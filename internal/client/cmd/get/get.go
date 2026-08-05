package get

import (
	"go.pabu.dev/gamesync/internal/client/config"

	"github.com/spf13/cobra"
)

type getCmd struct {
	Cmd *cobra.Command
}

func New(conf *config.Config) *getCmd {
	root := getCmd{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get resource",
	}
	cmd.AddCommand(
		newUserCmd(conf).cmd,
		newProfileCmd(conf).cmd,
		newRepoCmd(conf).cmd,
		newBranchCmd(conf).cmd,
	)

	root.Cmd = cmd
	return &root
}
