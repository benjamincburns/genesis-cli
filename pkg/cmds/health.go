package cmds

import (
	"io/ioutil"

	"github.com/whiteblock/genesis-cli/pkg/auth"
	"github.com/whiteblock/genesis-cli/pkg/config"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

func health(name, uri string) {
	client, err := auth.GetClient()
	if err != nil {
		util.ErrorFatal(err)
	}

	resp, err := client.Get(conf.APIEndpoint() + uri)
	if err != nil {
		util.Errorf(name+" unhealthy: %s", err.Error())
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		util.ErrorFatal(err)
	}
	if resp.StatusCode != 200 {
		util.Errorf(name+" unhealthy: %s", string(data))
		return
	}
	util.PrintKV(0, name, string(data))
}

var healthCmd = &cobra.Command{
	Use:     "health",
	Short:   "check the health of the backend services",
	Long:    `check the health of the backend services`,
	Aliases: []string{},
	Hidden:  true,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		health("registrar", "/api/v1/registrar/health")
		health("filehandler", "/api/v1/files/status")
		health("container-api", config.ContainerAPI+"/health")
		health("billing", conf.BillingHealthURI)

	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
}
