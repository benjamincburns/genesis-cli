package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/auth"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate the CLI with new credentials",
	Long: `Authenticate the CLI with new credentials. This command will attempt to open a tab in your browser,
where you can complete the login process. If your browser window does not open, it will also print a URL 
that you can paste in your browser to complete the login process.`,
	Aliases: []string{},
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
