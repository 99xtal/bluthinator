package commands

import (
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

func newProgressBar(p *mpb.Progress, total int64, name string) *mpb.Bar {
	return p.New(total,
		mpb.BarStyle().Lbound("|").Filler("=").Tip(">").Padding("-").Rbound("|"),
		mpb.PrependDecorators(
			decor.Name(name),
			decor.CountersNoUnit("%d/%d"),
			decor.Name(" ("),
			decor.Percentage(),
			decor.Name(")"),
		),
		mpb.AppendDecorators(
			decor.Elapsed(decor.ET_STYLE_GO),
			decor.EwmaSpeed(0, " %.2f ops/s", 60),
			decor.Name(" (ETA: "),
			decor.EwmaETA(decor.ET_STYLE_GO, 60),
			decor.Name(")"),
		),
	)
}