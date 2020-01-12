package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:     "stop <id>",
	Short:   "Stop a test",
	Long:    `Stop a test`,
	Aliases: []string{"kill", "rm"},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		isDef, err := cmd.Flags().GetBool("def")
		if err != nil {
			util.ErrorFatal(err)
		}
		err = service.StopTest(args[0], isDef)
		if err != nil {
			util.ErrorFatal(err)
		}
		util.Print("stop request sent")
	},
}

func init() {
	stopCmd.Flags().BoolP("def", "d", false, "treat the id as a definition")
	rootCmd.AddCommand(stopCmd)
}
