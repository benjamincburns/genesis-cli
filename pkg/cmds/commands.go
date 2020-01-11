package cmds

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var commandsCmd = &cobra.Command{
	Use:    "commands <file>",
	Hidden: true,
	Short:  "get commands",
	Long:   `get commands!`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)

		tests, def, err := service.ProcessDefinitionFromFile(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}

		for i := range tests {
			util.Printf("%s:", def.Spec.Tests[i].Name)
			for j := range tests[i].Commands {
				util.PrintKV(1, j, "")
				for k := range tests[i].Commands[j] {
					util.PrintKV(2, tests[i].Commands[j][k].Order.Type, tests[i].Commands[j][k].Order.Payload)
				}
			}
		}
	},
}

var localCmd = &cobra.Command{
	Use:    "local <file>",
	Hidden: true,
	Short:  "run a file locally",
	Long:   `run a file locally`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)

		tests, _, err := service.ProcessDefinitionFromFile(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}

		data, err := json.Marshal(tests[0])
		if err != nil {
			util.ErrorFatal(err)
		}

		resp, err := http.DefaultClient.Post(conf.LocalGenesisURL, "application/json", bytes.NewReader(data))
		if err != nil {
			util.ErrorFatal(err)
		}

		defer resp.Body.Close()

		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			util.ErrorFatal(err)
		}
		if resp.StatusCode != 200 {
			util.ErrorFatal(string(data))
		}
		util.Print(string(data))

	},
}

func init() {
	rootCmd.AddCommand(commandsCmd)
	rootCmd.AddCommand(localCmd)
}
