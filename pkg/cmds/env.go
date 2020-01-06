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
	Long: `This command lists the environment variables provided to your container. These 
variables are useful to see the additional information available to you at runtime. 
You can also see what these values might look like by using the --values flag. However, 
the values you see may change slightly at runtime, ie, instead of an IP address of 192.168.2.1,
it might instead be 192.168.5.3`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		envs, def, err := service.DefinitionEnv(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}
		showIP, err := cmd.Flags().GetBool("values")
		for i := range envs {
			util.Printf("%s:", def.Spec.Tests[i].Name)
			out := []string{}
			for k, v := range envs[i] {
				envVar := fmt.Sprintf("\t%s", k)
				if showIP {
					envVar += fmt.Sprintf(": %s", v)
				}
				out = append(out, envVar)
			}
			sort.Strings(out)
			for j := range out {
				util.Print(out[j])
			}
		}
	},
}

func init() {
	envCmd.Flags().BoolP("values", "v", false, "show the expected values for each env var (might not be accurate)")
	rootCmd.AddCommand(envCmd)
}
