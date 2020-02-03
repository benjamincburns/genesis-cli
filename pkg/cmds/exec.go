package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whiteblock/utility/common"
)

var execCmd = &cobra.Command{
	Use:   "exec <target> <command>",
	Short: "exec",
	Long:  `exec`,
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 2, util.NoMaxArgs)
		testID := util.GetStringFlagValue(cmd, "test-id")
		if len(testID) == 0 {
			test, err := service.GetMostRecentTest("")
			if err != nil {
				util.ErrorFatal(err)
			}
			testID = test.ID
		}
		info, err := service.PrepareExec(common.Exec{
			Test:        testID,
			Target:      args[0],
			Command:     args[1:],
			Privileged:  util.GetBoolFlagValue(cmd, "privileged"),
			Interactive: util.GetBoolFlagValue(cmd, "interactive"),
			TTY:         util.GetBoolFlagValue(cmd, "tty"),
			Detach:      util.GetBoolFlagValue(cmd, "detach"),
		})
		if err != nil {
			util.ErrorFatal(err)
		}
		log.WithField("info", info).Debug("got the info")
		attach := common.ExecAttach{
			ExecInfo:    info,
			Detach:      util.GetBoolFlagValue(cmd, "detach"),
			Interactive: util.GetBoolFlagValue(cmd, "interactive"),
			TTY:         util.GetBoolFlagValue(cmd, "tty"),
		}
		if util.GetBoolFlagValue(cmd, "detach") {
			err = service.RunDetach(attach)
		} else {
			err = service.Attach(attach)
		}
		if err != nil {
			util.ErrorFatal(err)
		}
	},
}

func init() {
	execCmd.Flags().BoolP("interactive", "i", false, "")
	execCmd.Flags().BoolP("tty", "t", false, "")
	execCmd.Flags().BoolP("detach", "d", false, "")
	execCmd.Flags().BoolP("privileged", "p", false, "")
	execCmd.Flags().String("test-id", "", "get from a specify test")
	rootCmd.AddCommand(execCmd)
}
