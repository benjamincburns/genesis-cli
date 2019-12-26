package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/auth"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate the CLI with new credentials",
	Long:  `Authenticate the CLI with new credentials`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		_, err := auth.Login()
		if err != nil {
			util.ErrorFatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
