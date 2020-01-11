package cmds

import (
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
		util.CheckArguments(cmd, args, 0, 0)
		org, err := organization.Get(orgNameOrId, client)
		if err != nil {
			util.ErrorFatal(err)
		}

		util.Print(org.Name)
	},
}

func init() {
	rootCmd.AddCommand(orgCmd)
}
