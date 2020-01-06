package cmds

import (
	"fmt"

	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var burnCmd = &cobra.Command{
	Use:   "burn <file>",
	Short: "Check the total resource usage of a spec",
	Long: `Burn shows you the total resource usage of each test in the specification. 
This information simplifies to calculation of burn rate.`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)

		tests, def, err := service.ProcessDefinitionFromFile(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}
		counts := map[string]map[string]int64{}

		for i := range tests {
			name := def.Spec.Tests[i].Name
			counts[name] = map[string]int64{}
			for _, instance := range tests[i].ProvisionCommand.Instances {
				counts[name]["cpus"] += instance.CPUs
				counts[name]["memory"] += instance.Memory
				counts[name]["storage"] += instance.Storage
			}
		}

		humanizedCounts := map[string]map[string]string{}
		for k, v := range counts {
			humanizedCounts[k] = map[string]string{}
			for k2, c := range v {
				switch k2 {
				case "cpus":
					humanizedCounts[k][k2] = fmt.Sprintf("%d vCPUs", c)
				default:
					humanizedCounts[k][k2] = fmt.Sprintf("%d MB", c)
				}

			}
		}

		for key := range humanizedCounts {
			util.Printf("%s:", key)
			for k, v := range humanizedCounts[key] {
				util.PrintKV(1, k, v)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(burnCmd)
}
