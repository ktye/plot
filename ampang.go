package plot

import (
	"image"
	"image/color"

	"github.com/ktye/plot/vg"
)

// ampAngPlot is an implementation of the HiPlotter interface.
type ampAngPlot struct {
	plot *Plot
	Limits
	ampAngDimension
	drawer   vg.Drawer
	amp, ang *axes
}

// Vertical layout
// - (border)
// - (vFill)
// - titleHeight
// - ticLength
// - ampAreaHeight
// - ampAngSpace
// - angAreaHeight
// - ticLabelHeight
// - xlabelHeight
// - (vFill)
// - (border)
//
// Horizontal layout, left to right
// - (border)
// - (hFill)
// - ylabelWidth (y label is rotated)
// - ticLabelWidth
// - 2*ticLength
// - plotAreaWidth
// - ticLength
// - rightXYWidth
// - (hFill)
// - (border)
type ampAngDimension struct {
	ampAreaHeight  int
	ampAngSpace    int
	angAreaHeight  int
	plotAreaWidth  int
	rightXYWidth   int
	titleHeight    int
	xlabelHeight   int
	ylabelWidth    int
	ticLabelWidth  int
	ticLabelHeight int
	ticLength      int
}

// Create a new ampang plot.
// width height is the available space, the image will be smaller.
func (plt *Plot) NewAmpAng(d vg.Drawer) (p ampAngPlot, err error) {
	width, height := d.Size()
	p.drawer = d
	p.plot = plt
	p.Limits = plt.getAmpAngLimits()
	xtics := getXTics(p.Limits)
	ytics := getYTics(p.Limits)

	// Calculate Dimensions.
	border := plt.defaultBorder()
	p.titleHeight = plt.defaultTitleHeight()
	p.ticLength = plt.defaultTicLength()
	p.ampAngSpace = plt.defaultAmpAngSpace()
	p.ticLabelHeight = plt.defaultTicLabelHeight()
	p.ticLabelWidth = plt.defaultTicLabelWidth(append(ytics.Labels, "-180"))
	p.xlabelHeight = plt.defaultXlabelHeight()
	p.ylabelWidth = plt.defaultYlabelWidth()
	if p.rightXYWidth = 0; len(xtics.Labels) > 0 {
		p.rightXYWidth = plt.defaultRightXYWidth(xtics.Labels[len(xtics.Labels)-1])
	}

	// Needed space for decorations.
	hFix := func() int { return 2*border + 3*p.ticLength + p.ylabelWidth + p.ticLabelWidth + p.rightXYWidth }
	vFix := func() int {
		return 2*border + 2*p.ticLength + p.titleHeight + p.ampAngSpace + p.ticLabelHeight + p.xlabelHeight
	}

	// Available space for plotArea.
	hSpace := width - hFix()
	vSpace := height - vFix()

	// Make sure the plot is not too wide or too slender.
	if hSpace > 3*vSpace/2 {
		hSpace = 3 * vSpace / 2
	} else if vSpace > 2*hSpace {
		vSpace = 2 * hSpace
	}

	p.plotAreaWidth = hSpace
	p.ampAreaHeight = vSpace * 2 / 3
	p.angAreaHeight = vSpace - p.ampAreaHeight
	// p.angAreaHeight should be odd.
	if p.angAreaHeight%2 != 1 {
		p.angAreaHeight--
	}

	// Calculate (smaller) image dimensions.
	swidth := hFix() + p.plotAreaWidth
	sheight := vFix() + p.ampAreaHeight + p.angAreaHeight
	x0, y0 := (width-swidth)/2, (height-sheight)/2

	amp := plt.newAxes(
		x0+p.ylabelWidth+p.ticLabelWidth+2*p.ticLength+border,
		y0+p.titleHeight+p.ticLength+border,
		p.plotAreaWidth,
		p.ampAreaHeight,
		p.Limits,
		d,
	)
	amp.x0 = x0
	p.amp = &amp

	ang := plt.newAxes(
		x0+p.ylabelWidth+p.ticLabelWidth+2*p.ticLength+border,
		y0+p.titleHeight+p.ticLength+border+p.ampAreaHeight+p.ampAngSpace,
		p.plotAreaWidth,
		p.angAreaHeight,
		Limits{false, p.Limits.Xmin, p.Limits.Xmax, -180.0, 180.0, 0, 0},
		d,
	)
	ang.x0 = x0
	p.ang = &ang

	p.draw()
	return p, nil
}
func (p ampAngPlot) draw() {
	xtics := getXTics(p.Limits)
	ytics := getYTics(p.Limits)
	atics := Tics{Pos: []float64{-180, -90, 0, 90, 180}, Labels: []string{"-180", "-90", "0", "90", "180"}}
	p.amp.fillParentBackground()
	p.amp.drawXY(xyAmp{})
	p.ang.drawXY(xyAng{})
	p.amp.drawXYTics(xtics.Pos, ytics.Pos, nil, ytics.Labels)
	p.ang.drawXYTics(xtics.Pos, atics.Pos, xtics.Labels, atics.Labels)
	p.amp.drawTitle(p.plot.defaultTicLength())
	p.ang.drawXlabel()
	p.amp.drawYlabel()
	p.ang.inside.Paint()
	p.amp.inside.Paint()
	//p.amp.reset()
	p.drawer.Paint()
	p.amp.store()
}
func (p ampAngPlot) background() color.Color { return p.plot.defaultBackgroundColor() }
func (p ampAngPlot) zoom(x, y, dx, dy int) bool {
	if p.ang.isInside(x, y) {
		X0, _ := p.ang.toFloats(x, y+dy)
		X1, _ := p.ang.toFloats(x+dx, y)
		p.amp.limits.Xmin = X0
		p.amp.limits.Xmax = X1
		p.ang.limits.Xmin = X0
		p.ang.limits.Xmax = X1
	} else {
		X0, Y0 := p.amp.toFloats(x, y+dy)
		X1, Y1 := p.amp.toFloats(x+dx, y)
		p.amp.limits.Xmin = X0
		p.amp.limits.Xmax = X1
		p.amp.limits.Ymin = Y0
		p.amp.limits.Ymax = Y1
		p.ang.limits.Xmin = X0
		p.ang.limits.Xmax = X1
	}
	p.amp.reset()
	p.ang.reset()
	p.draw()
	return true
}
func (p ampAngPlot) pan(x, y, dx, dy int) bool {
	if p.ang.isInside(x, y) {
		X0, _ := p.ang.toFloats(x, y+dy)
		X1, _ := p.ang.toFloats(x+dx, y)
		DX := X1 - X0
		p.amp.limits.Xmin -= DX
		p.amp.limits.Xmax -= DX
		p.ang.limits.Xmin -= DX
		p.ang.limits.Xmax -= DX
	} else {
		X0, Y0 := p.amp.toFloats(x, y+dy)
		X1, Y1 := p.amp.toFloats(x+dx, y)
		DX := X1 - X0
		DY := Y1 - Y0
		p.amp.limits.Xmin -= DX
		p.amp.limits.Xmax -= DX
		p.amp.limits.Ymin += DY
		p.amp.limits.Ymax += DY
		p.ang.limits.Xmin -= DX
		p.ang.limits.Xmax -= DX
	}
	p.amp.reset()
	p.ang.reset()
	p.draw()
	return true
}
func (p ampAngPlot) limits() Limits     { return p.amp.limits }
func (p ampAngPlot) image() *image.RGBA { return p.drawer.(*vg.Image).RGBA }
func (p ampAngPlot) line(x0, y0, x1, y1 int) (complex128, bool) {
	if !p.amp.isInside(x0, y0) {
		return complex(0, 0), false
	}
	vec, X0, Y0, X1, Y1 := p.amp.line(x0, y0, x1, y1)
	p.plot.Lines = append(p.plot.Lines, Line{
		Id:    p.plot.nextNegativeLineId(),
		X:     []float64{X0, X1},
		Y:     []float64{Y0, Y1},
		C:     []complex128{complex(Y0, 0), complex(Y1, 0)},
		Style: DataStyle{Line: LineStyle{Width: 1, Color: -1}},
	})
	p.draw()
	return vec, true
}
func (p ampAngPlot) click(x, y int, snapToPoint, deleteLine bool) (Callback, bool) {
	if p.amp.isInside(x, y) {
		pi, ok := p.amp.click(x, y, xyAmp{}, snapToPoint)
		if ok == true && snapToPoint == false {
			p.plot.Lines = append(p.plot.Lines, Line{
				Id: p.plot.nextNegativeLineId(),
				X:  []float64{pi.X},
				Y:  []float64{pi.Y},
				C:  []complex128{complex(pi.Y, 0)},
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
	} else if p.ang.isInside(x, y) {
		pi, ok := p.ang.click(x, y, xyAng{}, snapToPoint)
		return Callback{PointInfo: pi}, ok
	}

	limits := p.amp.limits
	if x < p.amp.x && y < p.amp.y+p.amp.height {
		if x < p.amp.x-p.ticLabelWidth {
			return Callback{Type: UnitCallback}, true
		} else {
			return Callback{
				Type:   AxisCallback,
				Limits: limits,
			}, true
		}
	} else if y > p.ang.y+p.ang.height && y < p.ang.y+p.ang.height+p.ticLabelHeight {
		return Callback{
			Type:   AxisCallback,
			Limits: limits,
		}, true
	}

	return Callback{}, false
}
func (p ampAngPlot) highlight(id []HighlightID) *image.RGBA {
	p.amp.restore()
	p.amp.highlight(id, xyAmp{})
	p.ang.highlight(id, xyAng{})
	return p.drawer.(*vg.Image).RGBA
}
