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
		latest, err := cmd.Flags().GetBool("latest")
		if err != nil {
			util.ErrorFatal(err)
		}
		if latest {
			test, err := service.GetMostRecentTest(org)
			if err != nil {
				util.ErrorFatal(err)
			}
			util.Print(test.ID)
			return
		}

		tests, err := service.GetTests(org)
		if err != nil {
			util.ErrorFatal(err)
		}
		if len(tests) == 0 {
			util.Print("you do not have any active tests")
		}
		details, err := cmd.Flags().GetBool("details")
		if err != nil {
			util.ErrorFatal(err)
		}
		simple, err := cmd.Flags().GetBool("simple")
		if err != nil {
			util.ErrorFatal(err)
		}

		for i := range tests {
			if simple {
				util.Print(tests[i].ID)
				continue
			}
			util.PrintKV(0, "Test", tests[i].ID)
			if details {
				util.PrintKV(1, "Definition", tests[i].DefinitionID)
				util.PrintKV(1, "Started", tests[i].CreatedAt)
			}
		}
	},
}

func init() {
	testsCmd.Flags().BoolP("details", "d", false, "show more test details")
	testsCmd.Flags().BoolP("simple", "1", false, "list the tests without any labels or formatting")
	testsCmd.Flags().BoolP("latest", "l", false, "show the ID of only the most recent active test")
	rootCmd.AddCommand(testsCmd)
}
