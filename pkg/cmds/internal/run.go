package internal

import (
	"fmt"
	"os"
	"errors"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/whiteblock/genesis-cli/pkg/service"
	"github.com/whiteblock/genesis-cli/pkg/util"
	"github.com/whiteblock/mpb"
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
			if res.Message != "" {
				logger.WithFields(log.Fields{
					"test-id": id,
					"result":  res.Message,
				}).Info("Result")
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
