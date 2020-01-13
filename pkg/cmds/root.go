package cmds

import (
	"fmt"

	"github.com/whiteblock/genesis-cli/pkg/config"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	conf = config.NewConfig()
)

var rootCmd = &cobra.Command{
	Use:     "genesis",
	Version: "",
	Short:   "A utility for accessing Whiteblock Genesis",
	Long: `Whiteblock Genesis is the only fully automated platform that helps development teams 
quickly design and run end-to-end distributed systems tests. Whiteblock Genesis is here to 
accelerate your testing process to maturity.`,
}

var version string

func Execute(buildTime, commitHash string) {
	version = commitHash
	rootCmd.Version = fmt.Sprintf("%s-%s", buildTime, commitHash)
	if err := rootCmd.Execute(); err != nil {
		util.ErrorFatal(err)
	}
}
func init() {
	if conf.GenesisBanner != "" {
		util.Print(conf.GenesisBanner)
	}
	rootCmd.PersistentFlags().Bool("no-colors", false, "disable terminal colors")
	rootCmd.PersistentFlags().Bool("dev", false, "disable terminal colors")
	rootCmd.PersistentFlags().MarkHidden("dev")

	viper.BindPFlag("no-colors", rootCmd.PersistentFlags().Lookup("no-colors"))
}
