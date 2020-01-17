package cmds

import (
	"io/ioutil"
	"os"

	"github.com/whiteblock/genesis-cli/pkg/cmds/internal"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var composeCmd = &cobra.Command{
	Use:     "compose <file>",
	Short:   "Convert docker compose into a spec",
	Long:    `Convert docker compose into a spec`,
	Aliases: []string{},
	Hidden:  true,
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

		root, err := internal.SchemaFromCompose(data)
		if err != nil {
			util.ErrorFatal(err)
		}

		data, err = yaml.Marshal(root)
		if err != nil {
			util.ErrorFatal(err)
		}
		util.Print(string(data))

		internal.Lint(data)
	},
}

func init() {
	rootCmd.AddCommand(composeCmd)
}
