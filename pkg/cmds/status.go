package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:     "status <id>",
	Short:   "status",
	Long:    `status`,
	Hidden:  true,
	Aliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		res, err := service.GetStatus(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}
		util.PrintS(0, res)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
