package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"
)

var logsCmd = &cobra.Command{
	Use:   "logs [org]",
	Short: "Most recent logs",
	Long:  `Gives the last most recent logs of the tests for the organization`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 1)
		org := ""
		if len(args) > 0 {
			org = args[0]
		}
		items, err := service.GetLogs(org)
		if err != nil {
			fmt.Println(err)
		} else {
			for _, item := range items {
				fmt.Println(item.Message.Text)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
}
