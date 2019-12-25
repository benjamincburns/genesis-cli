package cmds

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/whiteblock/definition"
	"github.com/whiteblock/go-prettyjson"
)

var burnCmd = &cobra.Command{
	Use:   "burn <file>",
	Short: "Check the resource usage of a spec",
	Long:  `Check the resource usage of a spec`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)

		file, err := os.Open(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}
		defer file.Close()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			util.ErrorFatal(err)
		}
		v := viper.New()
		v.Set("verbosity", "FATAL")
		definition.ConfigureGlobalFromViper(v) //shush library

		def, err := definition.SchemaYAML(data)
		if err != nil {
			util.ErrorFatal(err)
		}

		tests, err := definition.GetTests(def)
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
		out, _ := prettyjson.Marshal(humanizedCounts)
		util.Print(string(out))
	},
}

func init() {
	rootCmd.AddCommand(burnCmd)
}
