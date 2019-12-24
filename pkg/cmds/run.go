package cmds

import (
	"os"

	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <file> [org]",
	Short: "Run a test",
	Long:  `Run a test`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 2)
		org := ""
		if len(args) > 1 {
			org = args[1]
		}
		res, err := service.TestExecute(args[0], org)
		if err != nil {
			util.Error(err)
			if len(res) > 0 {
				util.Errorf("Response : %s", res)
			}
			os.Exit(1)
		}
		util.Print(res)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
