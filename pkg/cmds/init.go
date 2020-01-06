package cmds

import (
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/spf13/cobra"
)

var specFormatStr = `services:
  - name: my-service
    image: %s
    resources:
      cpus: 1
      memory: 500 MB
      storage: 1 GiB
task-runners:
  - name: wait-5-minutes
    script:
      inline: sleep 600
tests:
  - name: my-test
    description: You can put a description here
    system:
      - name: my-system
        type: my-service
        count: %s
    phases:
      - name: your-phase
        tasks:
          - type: wait-5-minutes
            timeout: 600s
`

var initCmd = &cobra.Command{
	Use:   "init <docker image> <node count>",
	Short: "Easy quickstart, generates a boiler plate spec for your software",
	Long: `Automatically generates a test spec which will deploy the given container the given number
of times, and wait for 5 minutes before teardown. The test spec will be output to stdout by default.
For further instructions, please refer to our documentation https://docs.genesis.whiteblock.io/docs/intro/.
`,
	Aliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 2, 2)
		util.PlainPrintf(specFormatStr, args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
