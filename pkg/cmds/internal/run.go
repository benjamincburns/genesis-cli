package internal

import (
	"strings"
	"time"

	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/whiteblock/mpb"
)

func TrackRunStatus(bar *mpb.Bar, id string, total int64) {
	for {
		res, err := service.GetStatus(id)
		if err != nil {
			if strings.Contains(err.Error(), "could not find the status") {
				continue
			}
			util.ErrorFatal(err)
		}
		bar.SetCurrent(total - int64(res.StepsLeft))
		if res.StepsLeft == 0 || res.Finished == true {
			return
		}

		time.Sleep(500 * time.Millisecond)
	}
}
