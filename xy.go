package plot

import (
	"image"
	"image/color"
)

// xyPlot is a standard rectangular x-y plot.
// It implements the IPlotter interface.
type xyPlot struct {
	plot *Plot
	Limits
	xyDimension
	xtics, ytics Tics
	im           *image.RGBA
	ax           *axes
}

// Vertical layout
// - (border)
// - (vFill)
// - titleHeight
// - plotAreaHeight
// - ticLabelHeight
// - titleHeight
// - (vFill)
// - (border)
//
// Horizontal layout, left to right
// - (border)
// - (hFill)
// - titleHeight (y label is rotated)
// - ticLabelWidth
// - plotAreaWidth
// - rightXYWidth
// - (hFill)
// - (border)
type xyDimension struct {
	plotAreaHeight int
	plotAreaWidth  int
	rightXYWidth   int
	titleHeight    int
	ticLabelWidth  int
	ticLabelHeight int
}

// Create a new xy plot.
// Width, height is the available space, the image will be smaller.
func (plt *Plot) NewXY(width, height int) (p xyPlot, err error) {
	p.plot = plt
	p.Limits = plt.getXYLimits()
	xtics := getXTics(p.Limits)
	ytics := getYTics(p.Limits)

	// Calculate Dimensions.
	border := plt.defaultBorder()
	p.titleHeight = plt.defaultTitleHeight()
	p.ticLabelHeight = plt.defaultTicLabelHeight()
	p.ticLabelWidth = plt.defaultTicLabelWidth(ytics.Labels)
	if len(xtics.Labels) > 0 {
		p.rightXYWidth = plt.defaultRightXYWidth(xtics.Labels[len(xtics.Labels)-1])
	}

	// Needed space for decorations.
	hFix := 2*border + p.titleHeight + p.ticLabelWidth + p.rightXYWidth
	vFix := 2*border + 2*p.titleHeight + p.ticLabelHeight

	// Available space for plotArea.
	hSpace := width - hFix
	vSpace := height - vFix

	// Plot may be wide but not too slender.
	if vSpace > 2*hSpace {
		vSpace = 2 * hSpace
	}

	p.plotAreaWidth = hSpace
	p.plotAreaHeight = vSpace

	// Calculate (smaller) image dimensions.
	width = p.titleHeight + p.ticLabelWidth + p.plotAreaWidth + p.rightXYWidth
	height = 2*p.titleHeight + p.plotAreaHeight + p.ticLabelHeight

	// Create the image.
	p.im = image.NewRGBA(image.Rect(0, 0, width, height))

	ax := plt.newAxes(
		p.titleHeight+p.ticLabelWidth,
		p.titleHeight,
		p.plotAreaWidth,
		p.plotAreaHeight,
		p.Limits,
		p.im,
	)
	p.ax = &ax

	p.xtics = xtics
	p.ytics = ytics
	p.draw()
	return p, nil
}

// drawAxis draws the axis, data and ticks.
func (p xyPlot) draw() {
	p.ax.fillParentBackground()
	if p.plot.Type == Raster {
		p.ax.drawImage()
	} else {
		p.ax.drawXY(xyXY{})
	}
	p.ax.drawXYTics(p.xtics.Pos, p.ytics.Pos, p.xtics.Labels, p.ytics.Labels)
	p.ax.drawTitle(p.plot.defaultTicLength())
	p.ax.drawXlabel()
	p.ax.drawYlabel()
}

func (p xyPlot) background() color.Color {
	return p.plot.defaultBackgroundColor()
}

func (p xyPlot) zoom(x, y, dx, dy int) bool {
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

func (p xyPlot) pan(x, y, dx, dy int) bool {
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

func (p xyPlot) limits() Limits {
	return p.ax.limits
}

func (p xyPlot) image() *image.RGBA {
	return p.im
}

func (p xyPlot) line(x0, y0, x1, y1 int) (complex128, bool) {
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

func (p xyPlot) click(x, y int, snapToPoint bool) (Callback, bool) {
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

func (p xyPlot) highlight(id []HighlightID) *image.RGBA {
	p.ax.highlight(id, xyXY{})
	return p.im
}
