package admin

import (
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var owoCmd = &cobra.Command{
	Use:     "OwO",
	Short:   "OwO",
	Long:    `OwO`,
	Aliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		util.Print(`
  ___            ___  
 / _ \__      __/ _ \ 
| | | \ \ /\ / / | | |
| |_| |\ V  V /| |_| |
 \___/  \_/\_/  \___/ 
                      
`)
	},
}

func init() {
	Command.AddCommand(owoCmd)
}
