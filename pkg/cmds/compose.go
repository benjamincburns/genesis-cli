package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/cmds/internal"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var composeCmd = &cobra.Command{
	Use:     "compose <file>",
	Short:   "Convert docker compose into a spec",
	Long:    `Convert docker compose into a spec`,
	Aliases: []string{},
	Hidden:  true,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)

		isJSON, err := cmd.Flags().GetBool("json")
		if err != nil {
			util.ErrorFatal(err)
		}
		var data []byte
		if isJSON {
			data = internal.MustSchemaJSONFromCompose(util.MustReadFile(args[0]))
		} else {
			data = internal.MustSchemaYAMLFromCompose(util.MustReadFile(args[0]))
		}

		util.Print(string(data))
		internal.Lint(data)
	},
}

func init() {
	composeCmd.Flags().Bool("json", false, "file is a docker compose file")
	rootCmd.AddCommand(composeCmd)
}
