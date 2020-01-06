package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var manCmd = &cobra.Command{
	Hidden: true,
	Use:    "man <output dir>",
	Short:  "Generate man pages for the Whiteblock Genesis CLI",
	Long: `This command automatically generates up-to-date man pages for genesis's command
line interface into the specified directory.`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		header := &doc.GenManHeader{
			Section: "1",
			Title:   "whiteblock",
		}

		cmd.Root().DisableAutoGenTag = true

		err := doc.GenManTree(cmd.Root(), header, args[0])
		if err != nil {
			util.ErrorFatal(err)
		}
	},
}

func init() {
	manCmd.Flags().BoolP("no-check", "u", false, "do not check for file existance")
	rootCmd.AddCommand(manCmd)
}
