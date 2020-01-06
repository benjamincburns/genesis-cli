package cmds

import (
	"os"
	"path/filepath"

	"github.com/whiteblock/genesis-cli/pkg/parser"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var filesCmd = &cobra.Command{
	Use:   "files <file> ",
	Short: "List the local files referenced by a definition",
	Long:  `List the local files referenced by a definition, and optionally checks if they exist`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		basePath := filepath.Dir(args[0])
		files, err := parser.ExtractFiles(args[0])
		if err != nil {
			util.ErrorFatal(err)
		}
		max := 0
		filenames := make([]string, len(files))
		for i, file := range files {
			filenames[i] = filepath.Join(basePath, file)
			if len(filenames[i]) > max {
				max = len(filenames[i])
			}
		}

		noCheck, err := cmd.Flags().GetBool("no-check")
		if err != nil {
			util.ErrorFatal(err)
		}

		if noCheck {
			for _, name := range filenames {
				util.Print(name)
			}
			return
		}
		formatString := "%s %s"
		for _, name := range filenames {
			_, err := os.Lstat(name)
			if err == nil {
				util.Printf(formatString, name, "OK")
			} else {
				util.Errorf(formatString, name, "NOT FOUND")
			}
		}
	},
}

func init() {
	filesCmd.Flags().BoolP("no-check", "u", false, "do not check for file existance")
	rootCmd.AddCommand(filesCmd)
}
