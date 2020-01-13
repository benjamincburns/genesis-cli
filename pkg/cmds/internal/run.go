package internal

import (
	"strings"
	"time"

	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"

	"github.com/whiteblock/mpb"
)

func TrackRunStatusNoTTY(id string, total int64) {
	for {
		res, err := service.GetStatus(id)
		if err != nil {
			if strings.Contains(err.Error(), "could not find the status") {
				continue
			}
			util.ErrorFatal(err)
		}
		util.Printf("%d%%", total-int64(res.StepsLeft))
		if res.StepsLeft == 0 || res.Finished == true {
			return
		}
		time.Sleep(5000 * time.Millisecond)
	}
}

func TrackRunStatus(p *mpb.Progress, bar *mpb.Bar, id string, total int64) {
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

func AwaitStatus(id string, total int64) <-chan error {
	out := make(chan error)
	go func() {
		for {
			res, err := service.GetStatus(id)
			if err != nil {
				if strings.Contains(err.Error(), "could not find the status") {
					continue
				}
				out <- err
				return
			}
			if res.StepsLeft == 0 || res.Finished == true {
				out <- nil
				return
			}

			time.Sleep(500 * time.Millisecond)
		}
		out <- nil
	}()
	return out
}
