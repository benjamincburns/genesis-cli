package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var commandsCmd = &cobra.Command{
	Use:    "commands <file>",
	Hidden: true,
	Short:  "get commands",
	Long:   `get commands!`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)

		tests, def, err := service.ProcessDefinitionFromFile(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}
		meta, err := cmd.Flags().GetBool("meta")
		if err != nil {
			util.ErrorFatal(err)
		}

		for i := range tests {
			util.Printf("%s:", def.Spec.Tests[i].Name)
			for j := range tests[i].Commands {
				util.PrintKV(1, j, "")
				for k := range tests[i].Commands[j] {
					if meta {
						util.PrintKV(2, tests[i].Commands[j][k].Order.Type, tests[i].Commands[j][k].Meta)
					} else {
						util.PrintKV(2, tests[i].Commands[j][k].Order.Type, tests[i].Commands[j][k].Order.Payload)
					}

				}
			}
		}
	},
}

func init() {
	commandsCmd.Flags().BoolP("meta", "m", false, "show command meta")
	rootCmd.AddCommand(commandsCmd)
}
