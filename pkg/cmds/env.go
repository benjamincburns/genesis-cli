package cmds

import (
	"fmt"

	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
	"sort"
)

var envCmd = &cobra.Command{
	Use:   "env <file>",
	Short: "Get the whiteblock provided environment vars for a test spec",
	Long:  `Check the resource usage of a spec`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		envs, def, err := service.DefinitionEnv(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}
		for i := range envs {
			util.Printf("%s:", def.Spec.Tests[i].Name)
			out := []string{}
			for k, v := range envs[i] {
				out = append(out, fmt.Sprintf("\t%s: %s", k, v))
			}
			sort.Strings(out)
			for j := range out {
				util.Print(out[j])
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
}
