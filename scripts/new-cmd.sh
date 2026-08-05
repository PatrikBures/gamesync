#!/usr/bin/env bash

client=""
if [[ "$1" == "-c" ]]; then
    client="yes"
    shift
fi


if [[ $# -lt 2 ]]; then
    echo "less than 2 args"
    exit 2
elif [[ $# -gt 2 ]]; then
    echo "more thant 2 args"
    exit 2
fi

if ! [[ -d "./.git" ]]; then
    echo "not in root of repo"
    exit 3
fi

verb=$1
resource=$2


resource_cmd_content=""
IFS= read -r -d '' resource_cmd_content << EOF
package $verb

import (${client:+"
	\"go.pabu.dev/gamesync/internal/client\"
	api \"go.pabu.dev/gamesync/internal/ogen\""}
	"go.pabu.dev/gamesync/internal/client/config"

	"github.com/spf13/cobra"
)

type ${resource}Cmd struct {
	cmd *cobra.Command
	opts ${resource}Opts
}

type ${resource}Opts struct {}

func new${resource^}Cmd(conf *config.Config) *${resource}Cmd {
	root := ${resource}Cmd{}

	cmd := &cobra.Command{
		Use: "${resource}",
		Short: "SUMMARY_PLACEHOLDER",
		RunE: func(cmd *cobra.Command, args []string) error {${client:+"
			c, err := client.New(conf)
			if err != nil {
				return err
			}"}
			populate${resource^}Opts(&root.opts, args)

			return run${resource^}Cmd(${client:+"c, "}conf, &root.opts)
		},
	}

	root.cmd = cmd
	return &root
}

func populate${resource^}Opts(opts *${resource}Opts, args []string) error {
	return nil
}

func run${resource^}Cmd(${client:+"c *api.Client, "}conf *config.Config, opts *${resource}Opts) error {
	return nil
}
EOF

verb_cmd_content=""
IFS= read -r -d '' verb_cmd_content << EOF
package $verb

import (
	"go.pabu.dev/gamesync/internal/client/config"

	"github.com/spf13/cobra"
)

type ${verb}Cmd struct {
	Cmd *cobra.Command
}

func New(conf *config.Config) *${verb}Cmd {
	root := ${verb}Cmd{}

	cmd := &cobra.Command{
		Use: "${verb}",
		Short: "${verb^} resources",
	}

	cmd.AddCommand(
		new${resource^}Cmd(conf).cmd,
	)

	root.Cmd = cmd
	return &root
}
EOF

dir="./cmd/gamesync/${verb}"
verb_file="${dir}/${verb}.go"
if ! [[ -d "$dir" ]]; then
    mkdir -p "$dir"
    echo -n "$verb_cmd_content" > "$verb_file"
    echo "created $verb_file"
fi

resource_file="${dir}/${resource}.go"

if ! [[ -f "$resource_file" ]]; then
    echo -n "$resource_cmd_content" > "$resource_file"
    echo "created $resource_file"
fi
