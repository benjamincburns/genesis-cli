package cmds

import (
	"bytes"
	"io"
	"os"

	"github.com/whiteblock/genesis-cli/pkg/auth"
	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs <container> [lines]",
	Short: "view the logs of a container",
	Long:  `view the logs of a container`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 2)
		testID := util.GetStringFlagValue(cmd, "test-id")
		if len(testID) == 0 {
			test, err := service.GetMostRecentTest(util.GetStringFlagValue(cmd, "org"))
			if err != nil {
				util.ErrorFatal(err)
			}
			testID = test.ID
		}
		lines := "20"
		if len(args) > 1 {
			lines = args[1]
		}
		data, err := auth.GetRaw(conf.ContainerLogsURL(testID, args[0], lines))
		if err != nil {
			util.ErrorFatal(err)
		}
		io.Copy(os.Stdout, bytes.NewReader(data))
	},
}

func init() {
	logsCmd.Flags().StringP("test-id", "t", "", "get from a specify test")
	logsCmd.Flags().StringP("org", "o", "", "get from the most recent test of the specified org")
	rootCmd.AddCommand(logsCmd)
}
