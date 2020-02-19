package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var psCmd = &cobra.Command{
	Use:     "ps",
	Short:   "what is running",
	Long:    `what is running`,
	Aliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		testID := util.GetStringFlagValue(cmd, "test-id")
		if len(testID) == 0 {
			test, err := service.GetMostRecentTest(util.GetStringFlagValue(cmd, "org"))
			if err != nil {
				util.ErrorFatal(err)
			}
			testID = test.ID
		}
		res, err := service.ListContainers(testID)
		if err != nil {
			util.ErrorFatal(err)
		}

		util.PrintS(-1, res)
	},
}

func init() {
	psCmd.Flags().StringP("test-id", "t", "", "get from a specify test")
	psCmd.Flags().StringP("org", "o", "", "get from the most recent test of the specified org")
	rootCmd.AddCommand(psCmd)
}
