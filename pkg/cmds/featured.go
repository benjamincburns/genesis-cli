package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/auth"
	"github.com/whiteblock/genesis-cli/pkg/org"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var featuredCmd = &cobra.Command{
	Use:     "featured",
	Short:   "featured",
	Long:    `featured`,
	Aliases: []string{},
	Hidden:  true,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		var val []org.Organization
		err := auth.Get(conf.FeaturedOrgsURL()+"?start=0&max=200", &val)
		if err != nil {
			util.ErrorFatal(err)
		}
		util.PrintS(0, val)
	},
}

func init() {
	rootCmd.AddCommand(featuredCmd)
}
