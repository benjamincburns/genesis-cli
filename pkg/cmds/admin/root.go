package admin

import (
	"github.com/whiteblock/genesis-cli/pkg/config"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var (
	conf = config.NewConfig()
)
var Command = &cobra.Command{
	Use:     "admin",
	Short:   "here be dragons",
	Long:    `here be dragons`,
	Aliases: []string{},
	Hidden:  true,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		util.ErrorFatal("command not found")
	},
}
