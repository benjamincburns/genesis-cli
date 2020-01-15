package cmds

import (
	"sort"

	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var portsCmd = &cobra.Command{
	Use:     "ports <file>",
	Short:   "show the exposed ports",
	Long:    `show the exposed ports`,
	Aliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)

		dists, _, err := service.DefinitionDist(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}
		out := [][][]int{}
		for _, dist := range dists {
			resourceDist := (*dist)[len(*dist)-1]
			portDist := [][]int{}
			for _, d := range resourceDist {
				ports := []int{}
				for port := range d.Resource.Ports {
					ports = append(ports, port)
				}
				sort.Ints(ports)
				portDist = append(portDist, ports)
			}
			out = append(out, portDist)
		}
		util.PrintS(0, out)
	},
}

func init() {
	rootCmd.AddCommand(portsCmd)
}
