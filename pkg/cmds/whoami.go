package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/auth"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var whoamiCmd = &cobra.Command{
	Use:     "whoami",
	Short:   "returns your identity",
	Long:    `returns your identity`,
	Aliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		out := map[string]interface{}{}
		err := auth.Get(conf.APIEndpoint()+conf.GetSelfURI, &out)
		if err != nil {
			util.ErrorFatal(err)
		}
		util.PrintS(0, out)
	},
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}
