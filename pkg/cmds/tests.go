package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var testsCmd = &cobra.Command{
	Use:     "tests [org]",
	Short:   "List active tests",
	Long:    `List active tests`,
	Aliases: []string{"validate"},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 1)
		org := ""
		if len(args) == 1 {
			org = args[0]
		}

		tests, err := service.GetTests(org)
		if err != nil {
			util.ErrorFatal(err)
		}
		if len(tests) == 0 {
			util.Print("you do not have any active tests")
		}
		for i := range tests {
			util.Print(tests[i])
		}
	},
}

func init() {
	rootCmd.AddCommand(testsCmd)
}
