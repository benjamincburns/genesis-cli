package cmds

import (
	"fmt"

	"github.com/whiteblock/genesis-cli/pkg/auth"
	organization "github.com/whiteblock/genesis-cli/pkg/org"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var roleCmd = &cobra.Command{
	Use:     "role <org>",
	Short:   "returns your role in the org",
	Long:    `returns your role in the org`,
	Aliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		orgInfo, err := organization.GetOrgInfo(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}

		out := map[string]interface{}{}
		err = auth.Get(conf.APIEndpoint()+fmt.Sprintf(conf.GetOrgRoleURI, orgInfo.ID), &out)
		if err != nil {
			util.ErrorFatal(err)
		}
		util.PrintS(0, out)
	},
}

var roleStatusCmd = &cobra.Command{
	Use:     "role-status <org>",
	Short:   "returns your role in the org",
	Long:    `returns your role in the org`,
	Aliases: []string{},
	Hidden:  true,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		orgInfo, err := organization.GetOrgInfo(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}
		err = auth.Get(conf.APIEndpoint()+fmt.Sprintf(conf.CheckMemberURI, orgInfo.ID), nil)
		if err != nil {
			util.Errorf("Org Membership... Failed")
		} else {
			util.Print("Org Membership... OK")
		}

		err = auth.Get(conf.APIEndpoint()+fmt.Sprintf(conf.CheckAdminURI, orgInfo.ID), nil)
		if err != nil {
			util.Errorf("Org Admin... Failed")
		} else {
			util.Print("Org Admin... OK")
		}
	},
}

func init() {
	rootCmd.AddCommand(roleCmd)
	rootCmd.AddCommand(roleStatusCmd)
}
