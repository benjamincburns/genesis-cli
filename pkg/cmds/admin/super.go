package admin

import (
	"github.com/whiteblock/genesis-cli/pkg/auth"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var superCmd = &cobra.Command{
	Use:     "super <enabled|disabled> <user-id>",
	Short:   "returns your identity",
	Long:    `returns your identity`,
	Aliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 2, 2)
		switch args[0] {
		case "enabled":
			_, err := auth.Post(conf.UpdateUserSuperAdminURL(args[1]), nil)
			if err != nil {
				util.ErrorFatal(err)
			}
		case "disabled":
			_, err := auth.Delete(conf.UpdateUserSuperAdminURL(args[1]), nil)
			if err != nil {
				util.ErrorFatal(err)
			}
		}

		util.Print("success")
	},
}

func init() {
	Command.AddCommand(superCmd)
}
