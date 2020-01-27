package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
	"github.com/whiteblock/definition/pkg/search"
)

var maxCmd = &cobra.Command{
	Use:   "max <file>",
	Short: "Get the whiteblock provided environment vars for a test spec",
	Long:  `Get the max count of each service`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		def, err := service.DefinitionFromFile(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}
		max := search.FindServiceMaxCounts(def.Spec)
		util.PrintS(0, max)
	},
}

func init() {
	rootCmd.AddCommand(maxCmd)
}
