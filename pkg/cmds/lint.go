package cmds

import (
	"io/ioutil"
	"os"

	"github.com/whiteblock/genesis-cli/pkg/message"
	"github.com/whiteblock/genesis-cli/pkg/util"
	"github.com/whiteblock/genesis-cli/pkg/validate"

	"github.com/spf13/cobra"
	"github.com/whiteblock/definition"
)

var lintCmd = &cobra.Command{
	Use:     "lint <file>",
	Short:   "Check for errors without running a test",
	Long:    `Check for errors without running a test`,
	Aliases: []string{"validate"},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		f, err := os.Open(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}

		data, err := ioutil.ReadAll(f)
		if err != nil {
			util.ErrorFatal(err)
		}

		res, err := validate.Schema(data)
		if err != nil {
			util.ErrorFatal(err)
		}
		if !res.Valid() {
			for _, schemaErr := range res.Errors() {
				util.Error(schemaErr.String())
			}
			return
		}

		def, err := definition.SchemaYAML(data)
		if err != nil {
			util.ErrorFatal(err)
		}

		_, err = definition.GetTests(def)
		if err != nil {
			util.ErrorFatal(err)
		}
		util.Print(message.FilePassedValidation)
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)
}
