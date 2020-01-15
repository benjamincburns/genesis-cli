package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var forkCmd = &cobra.Command{
	Use:     "fork <def id> [org id]",
	Short:   "fork",
	Long:    `fork`,
	Aliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 2)
		orgID := ""
		if len(args) > 1 {
			orgID = args[1]
		}
		res, err := service.Fork(args[0], orgID)
		if err != nil {
			util.ErrorFatal(err)
		}
		util.Print("successfully forked the project!")
		util.PrintS(0, res)
	},
}

func init() {
	rootCmd.AddCommand(forkCmd)
}
