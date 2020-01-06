package cmds

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Hidden: true,
	Use:    "completion",
	Short:  "Generates bash completion scripts",
	Long: `To load completion run
. <(genesis completion)
`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
