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
You can also see what these values might look like by using the --values flag. The names of the environment variables
will be the same as returned by this command. If you would like to have them include examples 
of the values of each variable, you may use the --values flag. Keep in mind that you should
not use these values in place of the environment variable.

Environment Variables for IPs
The Genesis platform provides environment variables to give you the IP addresses of services in the network.
All environment variables for IPs will be in all caps and also have '-' replaced with an underscore.

Services
The environment variables for the IP addresses of Services will be of the format 
"{service}_SERVICE{instance_no}_{network}". So, if you have a service "foo-baz" on the network 
"bar" then the first instance's IP address would be given in the environment variable
"FOO_BAZ_SERVICE0_BAR".

Sidecars
The naming of environment variables for sidecars is very similar to that of Services, with a few differences. 
The service's IP in the sidecar network will be the instance name of service, i.e., "{service}_SERVICE{instance_no}".
The sidecars' IP environment variables are formatted as though the service is their network. For example,
to find the IP of a sidecar "soap-bar" to the 0th service instance of "foo-baz", you would check the value of
"SOAP_BAR_FOO_BAZ_SERVICE0".
`,
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
