package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/config"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		util.ErrorFatal(err)
	}
}
func init() {
	rootCmd.PersistentFlags().Bool("no-colors", false, "disable terminal colors")

	viper.BindPFlag("no-colors", rootCmd.PersistentFlags().Lookup("no-colors"))
}
