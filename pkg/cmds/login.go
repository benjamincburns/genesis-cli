package cmds

import (
	"encoding/json"

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
	Aliases: []string{"auth"},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		_, err := auth.Login()
		if err != nil {
			util.ErrorFatal(err)
		}
	},
}

var printAccessTokenCmd = &cobra.Command{
	Use:     "print-access-token",
	Short:   "show access token",
	Long:    `Show access token`,
	Aliases: []string{"auth"},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		token := auth.GetToken()
		if token == nil {
			util.Error("you are not logged in")
			_, err := auth.Login()
			if err != nil {
				util.ErrorFatal(err)
			}
		}
		token = auth.GetToken()
		if token == nil {
			util.ErrorFatal("failed to get the auth token")
		}
		data, err := json.Marshal(*token)
		if token == nil {
			util.ErrorFatal(err)
		}
		util.Print(string(data))
	},
}

func init() {
	loginCmd.AddCommand(printAccessTokenCmd)
	rootCmd.AddCommand(loginCmd)
}
