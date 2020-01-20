package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/auth"
	organization "github.com/whiteblock/genesis-cli/pkg/org"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var orgCmd = &cobra.Command{
	Use:     "org",
	Short:   "What org am I acting as?",
	Long:    `What org am I acting as?`,
	Aliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 1)
		if len(args) == 0 {
			client, err := auth.GetClient()
			if err != nil {
				util.ErrorFatal(err)
			}
			org, err := organization.Get("", client)
			if err != nil {
				util.ErrorFatal(err)
			}

			util.Print(org.Name)
			return
		}

		orgInfo, err := organization.GetOrgInfo(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}
		util.PrintS(0, orgInfo)
	},
}

func init() {
	rootCmd.AddCommand(orgCmd)
}
