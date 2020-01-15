package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:     "info <test id>",
	Short:   "fetch test info",
	Long:    `fetch test info`,
	Aliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		res, err := service.TestInfo(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}
		showSpec, err := cmd.Flags().GetBool("show-spec")
		if err != nil {
			util.ErrorFatal(err)
		}
		if !showSpec {
			res.SpecFile = ""
		}
		util.PrintS(0, res)
	},
}

func init() {
	infoCmd.Flags().Bool("show-spec", false, "print out the spec with the other info")
	rootCmd.AddCommand(infoCmd)
}
