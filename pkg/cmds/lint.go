package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/cmds/internal"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint <file>",
	Short: "Check for errors without running a test",
	Long: `Check for errors without running a test.
Will first check to make sure that the file exists and matches the test definition schema. 
Then will run an internal mock of the test, to check for other issues.`,
	Aliases: []string{"validate"},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)

		internal.Lint(util.MustReadFile(args[0]))
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)
}
