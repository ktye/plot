package plot

import (
	"image"
	"image/color"
	"math"
)

type waterfallPlot struct {
	plot *Plot
	Limits
	waterfallDimension
	xtics, ytics Tics
	im           *image.RGBA
	ax           *axes
}
type waterfallDimension struct {
	plotAreaHeight int
	plotAreaWidth  int
	rightXYWidth   int
	titleHeight    int
	xlabelHeight   int
	ylabelWidth    int
	ticLabelWidth  int
	ticLabelHeight int
	ticLength      int
}

func (plt *Plot) NewWaterfall(width, height int) (p waterfallPlot, err error) {
	p.plot = plt
	p.Limits = plt.getXYLimits()
	xtics := getXTics(p.Limits)
	ytics := getYTics(p.Limits)

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

	p.im = image.NewRGBA(image.Rect(0, 0, width, height))
	ax := plt.newAxes(
		p.ylabelWidth+p.ticLabelWidth+2*p.ticLength+border,
		p.titleHeight+p.ticLength+border,
		p.plotAreaWidth,
		p.plotAreaHeight,
		p.Limits,
		p.im,
	)
	ax.zSpace = zSpace
	p.ax = &ax

	p.xtics = xtics
	p.ytics = ytics
	p.draw()
	return p, nil

	return p, nil
}
func (p waterfallPlot) draw() {
	p.ax.fillParentBackground()
	p.ax.drawXY(xyXY{})
	p.ax.drawXYTics(p.xtics.Pos, p.ytics.Pos, p.xtics.Labels, p.ytics.Labels)
	p.ax.drawTitle(p.plot.defaultTicLength())
	p.ax.drawXlabel()
	p.ax.drawYlabel()
}

func (p waterfallPlot) background() color.Color { return p.plot.defaultBackgroundColor() }
func (p waterfallPlot) image() *image.RGBA      { return p.im }
func (p waterfallPlot) zoom(x, y, dx, dy int) bool {
	X0, Y0 := p.ax.toFloats(x, y+dy)
	X1, Y1 := p.ax.toFloats(x+dx, y)
	p.ax.limits.Xmin = X0
	p.ax.limits.Xmax = X1
	p.ax.limits.Ymin = Y0
	p.ax.limits.Ymax = Y1
	p.xtics = getXTics(p.ax.limits)
	p.ytics = getYTics(p.ax.limits)
	p.draw()
	return true
}
func (p waterfallPlot) pan(x, y, dx, dy int) bool {
	X0, Y0 := p.ax.toFloats(x, y+dy)
	X1, Y1 := p.ax.toFloats(x+dx, y)
	DX := X1 - X0
	DY := Y1 - Y0
	p.ax.limits.Xmin -= DX
	p.ax.limits.Xmax -= DX
	p.ax.limits.Ymin += DY
	p.ax.limits.Ymax += DY
	p.xtics = getXTics(p.ax.limits)
	p.ytics = getYTics(p.ax.limits)
	p.draw()
	return true
}
func (p waterfallPlot) limits() Limits { return p.ax.limits }
func (p waterfallPlot) highlight(id []HighlightID) *image.RGBA {
	p.ax.highlight(id, xyXY{})
	return p.im
}
func (p waterfallPlot) line(x0, y0, x1, y1 int) (complex128, bool) { //todo..
	if !p.ax.isInside(x0, y0) {
		return complex(0, 0), false
	}
	vec, X0, Y0, X1, Y1 := p.ax.line(x0, y0, x1, y1)
	p.plot.Lines = append(p.plot.Lines, Line{
		Id:    p.plot.nextNegativeLineId(),
		X:     []float64{X0, X1},
		Y:     []float64{Y0, Y1},
		Style: DataStyle{Line: LineStyle{Width: 1, Color: -1}},
	})
	p.draw()
	return vec, true
}
func (p waterfallPlot) click(x, y int, snapToPoint, deleteLine bool) (Callback, bool) { //todo..
	if !p.ax.isInside(x, y) {
		limits := p.ax.limits
		if x < p.ax.x {
			if x < p.ax.x-p.ticLabelWidth {
				return Callback{Type: UnitCallback}, true
			} else {
				return Callback{
					Type:   AxisCallback,
					Limits: limits,
				}, true
			}
		} else if y > p.ax.y+p.ax.height && y < p.ax.y+p.ax.height+p.ticLabelHeight {
			return Callback{
				Type:   AxisCallback,
				Limits: limits,
			}, true
		}
		return Callback{}, false
	}
	if p.plot.Type == Raster {
		pi, ok := p.ax.clickImage(x, y)
		return Callback{PointInfo: pi}, ok
	} else {
		pi, ok := p.ax.click(x, y, xyXY{}, snapToPoint)
		if ok && snapToPoint == false {
			p.plot.Lines = append(p.plot.Lines, Line{
				Id: p.plot.nextNegativeLineId(),
				X:  []float64{pi.X},
				Y:  []float64{pi.Y},
				Style: DataStyle{
					Marker: MarkerStyle{
						Marker: CrossMarker,
						Color:  -1,
						Size:   3,
					},
				},
			})
			p.draw()
			return Callback{Type: MeasurePoint, PointInfo: pi}, ok
		}
		return Callback{PointInfo: pi}, ok
	}
}
