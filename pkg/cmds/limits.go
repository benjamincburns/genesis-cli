package cmds

import (
	"fmt"

	"github.com/whiteblock/genesis-cli/pkg/auth"
	"github.com/whiteblock/genesis-cli/pkg/config"
	organization "github.com/whiteblock/genesis-cli/pkg/org"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var limitsCmd = &cobra.Command{
	Use:     "limits <org>",
	Short:   "get your genesis limits",
	Long:    `get your genesis limits`,
	Aliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		orgInfo, err := organization.GetOrgInfo(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}
		out := map[string]interface{}{}
		err = auth.Get(conf.APIEndpoint()+fmt.Sprintf(conf.LimitsURI, orgInfo.ID, config.Product), &out)
		if err != nil {
			util.ErrorFatal(err)
		}
		util.PrintS(0, out)
	},
}

func init() {
	rootCmd.AddCommand(limitsCmd)
}
