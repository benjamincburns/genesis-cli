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

		data := internal.MustSchemaYAMLFromCompose(util.MustReadFile(args[0]))

		util.Print(string(data))
		internal.Lint(data)
	},
}

func init() {
	rootCmd.AddCommand(composeCmd)
}
