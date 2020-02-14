package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/cmds/admin"
)

func init() {
	rootCmd.AddCommand(admin.Command)
}
