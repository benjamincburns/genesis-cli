package internal

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"
	"github.com/whiteblock/mpb"
	"github.com/whiteblock/mpb/decor"
)

func TrackRunStatusNoTTY(id string, total int64) {
	logger := log.New()
	log.SetOutput(os.Stdout)
	for {
		res, err := service.GetStatus(id)
		if err != nil {
			if strings.Contains(err.Error(), "could not find the status") {
				continue
			}
			util.ErrorFatal(err)
		}

		logger.WithFields(log.Fields{
			"test-id":          id,
			"steps-percentage": fmt.Sprintf("%.2f%%", (float64(total-int64(res.StepsLeft))/float64(total))*100),
		}).Info("Progress")
		if res.StepsLeft == 0 || res.Finished == true {
			// if the msg is not empty, there is an error
			if res.Message != "" {
				logger.WithFields(log.Fields{
					"test-id": id,
					"result":  res.Message,
				}).Error("Result")
				util.ErrorFatal(res.Message)
			}
			return
		}
		time.Sleep(5000 * time.Millisecond)
	}
}

func TrackRunStatus(p *mpb.Progress, bar *mpb.Bar, id string, info util.BarInfo) <-chan error {
	out := make(chan error)
	go func() {
		phase := ""
		for {
			res, err := service.GetStatus(id)
			if err != nil {
				if strings.Contains(err.Error(), "could not find the status") {
					continue
				}
				out <- err
				return
			}

			if phase != res.Phase {
				phase = res.Phase
				b2 := p.AddBar(info.Total,
					mpb.BarParkTo(bar),
					mpb.BarClearOnComplete(),
					mpb.PrependDecorators(
						decor.Name(info.Name),
						decor.Percentage(decor.WCSyncSpace),
					),
					mpb.AppendDecorators(
						decor.OnComplete(decor.Name(phase), "done!"),
					),
				)
				bar.SetCurrent(info.Total)
				bar = b2
			}
			bar.SetCurrent(info.Total - int64(res.StepsLeft))
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
