#!/usr/bin/env bash

client=""
single=""
for _ in {1..2}; do
    case "$1" in 
    "-c")
        client="yes"
        shift
    ;;
    "-s")
        single="yes"
        shift
    esac
done


expected_arg_count=2
if [[ -n "$single" ]]; then
    expected_arg_count=1
fi

if [[ $# -lt $expected_arg_count ]]; then
    echo "less than $expected_arg_count args"
    exit 2
elif [[ $# -gt $expected_arg_count ]]; then
    echo "more thant $expected_arg_count args"
    exit 2
fi

if ! [[ -d "./.git" ]]; then
    echo "not in root of repo"
    exit 3
fi

verb=$1
resource=$2

[[ -n "$single" ]] && cmdString="Cmd" || cmdString="cmd"
[[ -n "$single" ]] && newCmdString="New" || newCmdString="new${resource^}Cmd"

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
	${cmdString} *cobra.Command
	opts ${resource}Opts
}

type ${resource}Opts struct {}

func ${newCmdString}(conf *config.Config) *${resource}Cmd {
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

	root.${cmdString} = cmd
	return &root
}

func populate${resource^}Opts(opts *${resource}Opts, args []string) {}

func run${resource^}Cmd(${client:+"c *api.Client, "}conf *config.Config, opts *${resource}Opts) error {
	return nil
}
EOF

dir="./internal/client/cmd/${verb}"
mkdir -p "$dir"

if [[ -n "$single" ]]; then
    resource_file="${dir}/${verb}.go"
    echo -n "$resource_cmd_content" > "$resource_file"
    echo "created single $resource_file"
    exit 0
fi

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

verb_file="${dir}/${verb}.go"
if ! [[ -f "$verb_file" ]]; then
    echo -n "$verb_cmd_content" > "$verb_file"
    echo "created $verb_file"
fi


resource_file="${dir}/${resource}.go"
if ! [[ -f "$resource_file" ]]; then
    echo -n "$resource_cmd_content" > "$resource_file"
    echo "created $resource_file"
fi
