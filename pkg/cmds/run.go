package cmds

import (
	"io/ioutil"
	"os"

	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/Pallinder/go-randomdata"
	"github.com/spf13/cobra"
	"github.com/whiteblock/definition"
)

var runCmd = &cobra.Command{
	Use:   "run <file> [org]",
	Short: "Run a test",
	Long:  `Run a test`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 2)
		org := ""
		if len(args) > 1 {
			org = args[1]
		}

		file, err := os.Open(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}
		defer file.Close()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			util.ErrorFatal(err)
		}

		def, err := definition.SchemaYAML(data)
		if err != nil {
			util.ErrorFatal(err)
		}

		tests, err := definition.GetTests(def)
		if err != nil {
			util.ErrorFatal(err)
		}

		var dns []string
		dnsEnabled, err := cmd.Flags().GetBool("dns")
		if err != nil {
			util.ErrorFatal(err)
		}
		if dnsEnabled {
			for range tests {
				dns = append(dns, randomdata.SillyName())
			}
		}

		res, err := service.TestExecute(args[0], org, dns)
		if err != nil {
			util.Error(err)
			if len(res) > 0 {
				util.Errorf("Response : %s", res)
			}
			os.Exit(1)
		}

		util.Print(res)

		if dnsEnabled {
			for i := range tests {
				for j := range tests[i].ProvisionCommand.Instances {
					util.Printf("%s-%d.%s", dns[i], j, conf.BiomeDNSZone)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("dns", "d", false, "use dns so you can access your deployments")
}
