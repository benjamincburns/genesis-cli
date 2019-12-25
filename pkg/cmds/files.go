package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/parser"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var filesCmd = &cobra.Command{
	Use:   "files <file> ",
	Short: "List the local files referenced by a definition",
	Long:  `List the local files referenced by a definition`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)

		files, err := parser.ExtractFiles(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}
		for _, file := range files {
			util.Print(file)
		}
	},
}

func init() {
	rootCmd.AddCommand(filesCmd)
}
