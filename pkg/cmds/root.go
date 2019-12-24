package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/config"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var conf = config.NewConfig()

var rootCmd = &cobra.Command{
	Use:     "genesis",
	Version: "",
	Short:   "",
	Long:    ``,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		util.PrintErrorFatal(err)
	}
}
