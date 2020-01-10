package util

import (
	"github.com/whiteblock/mpb"
	"github.com/whiteblock/mpb/decor"
)

type BarInfo struct {
	Name  string
	Total int64
}

type Awaiter interface {
	Wait()
}

func SetupBars(bars []BarInfo) (Awaiter, []*mpb.Bar) {
	p := mpb.New()
	out := make([]*mpb.Bar, len(bars))

	for i := range bars {
		bar := p.AddBar(int64(bars[i].Total),
			mpb.PrependDecorators(
				// simple name decorator
				decor.Name(bars[i].Name),
				// decor.DSyncWidth bit enables column width synchronization
				decor.Percentage(decor.WCSyncSpace),
			),
			mpb.AppendDecorators(
				// replace ETA decorator with "done" message, OnComplete event
				decor.OnComplete(
					decor.Name(bars[i].Name), "done",
				),
			),
		)
		out[i] = bar
	}
	return p, out
}
