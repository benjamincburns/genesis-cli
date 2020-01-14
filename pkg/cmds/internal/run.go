package internal

import (
	"errors"
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
		util.Printf("%f%%", float64(total-int64(res.StepsLeft))/float64(total))
		if res.StepsLeft == 0 || res.Finished == true {
			if res.Message != "" {
				util.PrintKV(0, id, res.Message)
			}
			return
		}
		time.Sleep(5000 * time.Millisecond)
	}
}

func TrackRunStatus(p *mpb.Progress, bar *mpb.Bar, id string, total int64) <-chan error {
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
			bar.SetCurrent(total - int64(res.StepsLeft))
			if res.StepsLeft == 0 || res.Finished == true {
				if res.Message != "" {
					out <- errors.New(res.Message)
				} else {
					out <- nil
				}
				return
			}

			time.Sleep(500 * time.Millisecond)
		}
	}()
	return out
}
