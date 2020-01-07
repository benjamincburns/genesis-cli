package cmds

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const CommandName = "genesis"

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update the cli",
	Long:    `Updates the cli to the latest version`,
	Aliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 0)
		defer os.Exit(0) //ensure that the program exits when this function returns

		update := <-service.CheckForUpdates(version)
		if !update {
			util.Print("already at latest version")
			return
		}

		endpoint := fmt.Sprintf(conf.CLIURL, runtime.GOOS, runtime.GOARCH)
		log.WithFields(log.Fields{"ep": endpoint}).Trace("fetching the binary data")
		resp, err := http.Get(endpoint)
		if err != nil {
			util.ErrorFatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			util.ErrorFatal("got back a non-200 status code")
		}

		binary, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			util.ErrorFatal(err)
		}
		log.WithFields(log.Fields{"size": len(binary)}).Trace("fetched the binary data")

		location, err := exec.LookPath(CommandName)
		if err != nil {
			util.ErrorFatal(err)
		}
		log.WithFields(log.Fields{"loc": location}).Trace("got the binary location")
		location, err = filepath.Abs(location)
		if err != nil {
			util.ErrorFatal(err)
		}

		location, err = filepath.EvalSymlinks(location)
		if err != nil {
			util.ErrorFatal(err)
		}

		fi, err := os.Lstat(location)
		if err != nil {
			util.ErrorFatal(err)
		}
		err = ioutil.WriteFile(filepath.Join(location, ".tmp"), binary, fi.Mode())
		if err != nil {
			util.ErrorFatal(err)
		}
		err = os.Rename(location+".tmp", location)
		if err != nil {
			util.ErrorFatal(err)
		}
		util.Print("updated successfully.")

	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
