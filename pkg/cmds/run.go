package cmds

import (
	"os"
	"strings"

	"github.com/whiteblock/genesis-cli/pkg/cmds/internal"
	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/Pallinder/go-randomdata"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <file> [org]",
	Short: "Run a test",
	Long: `Run the tests specified in the given test definition file. If it is your first time running the
	command, you will need to specify the org to deploy to. Subsequent runs will use this value until you specify it
	once again.`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 2)
		org := ""
		if len(args) > 1 {
			org = args[1]
		}

		tests, def, err := service.ProcessDefinitionFromFile(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}

		var dns []string
		dnsDisabled, err := cmd.Flags().GetBool("no-dns")
		if err != nil {
			util.ErrorFatal(err)
		}
		awaitDisabled, err := cmd.Flags().GetBool("no-await")
		if err != nil {
			util.ErrorFatal(err)
		}

		if !dnsDisabled {
			for range tests {
				dns = append(dns, strings.ToLower(randomdata.SillyName()))
			}
		}

		res, ids, err := service.TestExecute(args[0], org, dns)
		if err != nil {
			util.Error(err)
			if len(res) > 0 && !strings.Contains(res, "<!DOCTYPE html>") {
				util.Errorf("Response : %s", res)
			}
			os.Exit(1)
		}

		util.Print(res)

		if !dnsDisabled {
			for i := range tests {
				util.Printf("Test: %s", def.Spec.Tests[i].Name)
				for j := range tests[i].ProvisionCommand.Instances {
					util.Printf("\t%s-%d.%s", dns[i], j, conf.BiomeDNSZone)
				}
			}
		}

		if !awaitDisabled {
			infos := []util.BarInfo{}
			for i := range ids {
				infos = append(infos, util.BarInfo{
					Name:  def.Spec.Tests[i].Name,
					Total: tests[i].GuessSteps(),
				})
			}
			awaiter, bars := util.SetupBars(infos)

			for i := range ids {
				go internal.TrackRunStatus(bars[i], ids[i], infos[i].Total)
			}
			awaiter.Wait()
		}

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("no-dns", "d", false, "disable assigning a DNS name to your deployment")
	runCmd.Flags().BoolP("no-await", "a", false, "don't wait for completion, exit immediately after sending test")
}
