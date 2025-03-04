package plot

import (
	"math"

	"github.com/ktye/plot/vg"
)

type waterfallPlot struct {
	xyPlot
}

func (plt *Plot) NewWaterfall(d vg.Drawer) (p waterfallPlot, err error) {
	width, height := d.Size()
	p.drawer = d
	p.plot = plt
	p.Limits = plt.getXYLimits()
	xtics := getXTics(p.Limits)
	ytics := getYTics(p.Limits)
	ztics := getZTics(p.Limits)

	border := plt.defaultBorder()
	p.ticLength = plt.defaultTicLength()
	p.titleHeight = plt.defaultTitleHeight()
	p.ticLabelHeight = plt.defaultTicLabelHeight()
	p.ticLabelWidth = plt.defaultTicLabelWidth(ytics.Labels)
	p.xlabelHeight = plt.defaultXlabelHeight()
	p.ylabelWidth = plt.defaultYlabelWidth()
	if len(xtics.Labels) > 0 {
		p.rightXYWidth = plt.defaultRightXYWidth(xtics.Labels[len(xtics.Labels)-1])
	}

	hFix := func() int { return 2*border + 3*p.ticLength + p.ylabelWidth + p.ticLabelWidth + p.rightXYWidth }
	vFix := func() int { return 2*border + 2*p.ticLength + p.titleHeight + p.ticLabelHeight + p.xlabelHeight }
	hSpace := width - hFix()
	vSpace := height - vFix()
	if vSpace > 2*hSpace {
		vSpace = 2 * hSpace
	}
	zSpace := int(math.Min(float64(hSpace), float64(vSpace)) * (1 - math.Sqrt2/2))

	p.plotAreaWidth = hSpace
	p.plotAreaHeight = vSpace
	width = hFix() + p.plotAreaWidth
	height = vFix() + p.plotAreaHeight

	ax := plt.newAxes(
		p.ylabelWidth+p.ticLabelWidth+2*p.ticLength+border,
		p.titleHeight+p.ticLength+border,
		p.plotAreaWidth,
		p.plotAreaHeight,
		p.Limits,
		d,
	)
	ax.zSpace = zSpace
	p.ax = &ax

	p.xtics = xtics
	p.ytics = ytics
	p.ztics = ztics
	p.draw()
	return p, nil
}
