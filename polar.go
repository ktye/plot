package plot

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

// polarPlot is an implementation of the HiPlotter interface.
type polarPlot struct {
	plot   *Plot // underlying plot structure
	Limits       // computed axis limits
	ring   bool
	polarDimension
	im   *image.RGBA
	axes *axes
}

// Vertical layout
// - (border)
// - (vFill)
// - ticLength
// - titleHeight
// - ticLabelHeight
// - polarDiameter
// - ticLength
// - ticLabelHeight
// - (vFill)
// - (border)
//
// Horizontal layout, left to right
// - (border)
// - (hFill)
// - ticLabelWidth
// - polarDiameter
// - ticLabelWidth
// - (hFill)
// - (border)
type polarDimension struct {
	polarDiameter  int
	titleHeight    int
	ticLabelWidth  int
	ticLabelHeight int
	ticLength      int
}

// Create a new polar coordinate system in the subimage.
// Width and height are available areas on input, the image (p.im) will be smaller.
func (plt *Plot) NewPolar(width, height int, isRing bool) (p polarPlot, err error) {
	p.plot = plt
	p.ring = isRing
	p.Limits = plt.getPolarLimits(p.ring)
	if p.Limits.Ymax == 0 || math.IsNaN(p.Limits.Ymax) {
		return p, fmt.Errorf("cannot calculate polar limits (no data?)")
	}

	// Calculate Dimensions.
	border := plt.defaultBorder()
	p.ticLabelHeight = plt.defaultTicLabelHeight()
	p.ticLabelWidth = plt.defaultPolarTicLabelWidth()
	p.titleHeight = plt.defaultTitleHeight()
	p.ticLength = plt.defaultTicLength()

	// Needed space for decorations.
	hFix := 2*border + 2*p.ticLabelWidth
	vFix := 2*border + p.titleHeight + 2*p.ticLabelHeight + 2*p.ticLength

	// Available space for plotArea.
	hSpace := width - hFix
	vSpace := height - vFix

	if hSpace < 0 || vSpace < 0 {
		p.polarDiameter = 0
	} else {
		if hSpace < vSpace {
			p.polarDiameter = hSpace
		} else {
			p.polarDiameter = vSpace
		}
	}
	// Make sure polarDiameter is odd, then the axis lines are on a pixel and
	// not in between (distributed over 2px).
	if p.polarDiameter%2 != 1 {
		p.polarDiameter--
	}

	// Calculate (smaller) image dimensions.
	width = 2*p.ticLabelWidth + p.polarDiameter
	height = p.titleHeight + 2*p.ticLabelHeight + 2*p.ticLength + p.polarDiameter
	if p.polarDiameter < 1 {
		return p, fmt.Errorf("image space is too small")
	}

	// Create Image.
	p.im = image.NewRGBA(image.Rect(0, 0, width, height))
	ax := plt.newAxes(
		p.ticLabelWidth,
		p.titleHeight+p.ticLabelHeight,
		p.polarDiameter,
		p.polarDiameter,
		p.Limits,
		p.im,
	)
	p.axes = &ax

	p.draw(false)
	return p, nil
}

func (p polarPlot) draw(noTics bool) {
	p.axes.fillParentBackground()
	p.axes.drawPolar(p.ring)
	if noTics == false {
		p.axes.drawPolarTics(p.ring)
	}
	p.axes.drawTitle(p.ticLabelHeight + 3) // The title needs extra spacing because of the angule label.
}

func (p polarPlot) background() color.Color {
	return p.plot.defaultBackgroundColor()
}

func (p polarPlot) image() *image.RGBA {
	return p.im
}

func (p polarPlot) zoom(x, y, dx, dy int) bool {
	// Keep it square.
	if dx > dy {
		dy = dx
	} else {
		dx = dy
	}
	X0, Y0 := p.axes.toFloats(x, y+dy)
	X1, Y1 := p.axes.toFloats(x+dx, y)
	p.axes.limits.Xmin = X0
	p.axes.limits.Xmax = X1
	p.axes.limits.Ymin = Y0
	p.axes.limits.Ymax = Y1
	p.axes.drawPolarDataOnly()
	return true
}

func (p polarPlot) pan(x, y, dx, dy int) bool {
	X0, Y0 := p.axes.toFloats(x, y+dy)
	X1, Y1 := p.axes.toFloats(x+dx, y)
	DX := X1 - X0
	DY := Y1 - Y0
	p.axes.limits.Xmin -= DX
	p.axes.limits.Xmax -= DX
	p.axes.limits.Ymin += DY
	p.axes.limits.Ymax += DY
	p.axes.drawPolarDataOnly()
	return true
}

func (p polarPlot) limits() Limits {
	return p.axes.limits
}

func (p polarPlot) line(x0, y0, x1, y1 int) (complex128, bool) {
	if !p.axes.isInside(x0, y0) {
		return complex(0, 0), false
	}
	vec, X0, Y0, X1, Y1 := p.axes.line(x0, y0, x1, y1)
	vec = complex(imag(vec), real(vec))
	p.plot.Lines = append(p.plot.Lines, Line{
		Id:    p.plot.nextNegativeLineId(),
		C:     []complex128{complex(Y0, X0), complex(Y1, X1)},
		Style: DataStyle{Line: LineStyle{Width: 1, Color: -1}},
	})
	if p.axes.limits.isPolarLimits() {
		p.draw(true)
	} else { // axes are zoomed.
		p.axes.drawPolarDataOnly()
	}
	return vec, true
}

func (p polarPlot) click(x, y int, snapToPoint bool) (Callback, bool) {
	if !p.axes.isInside(x, y) {
		return Callback{
			Type:   AxisCallback,
			Limits: Limits{Xmax: p.axes.limits.Xmax, Ymax: p.axes.limits.Ymax},
		}, true
	}
	pi, ok := p.axes.click(x, y, p.axes.xyRing(), snapToPoint)
	if ok && snapToPoint == false {
		if p.ring {
			pi.C = complex(0, 0)
			pi.X = pi.X
			pi.Y = pi.Y
		} else {
			pi.C = complex(pi.Y, pi.X)
			pi.X = 0
			pi.Y = 0
		}
		p.plot.Lines = append(p.plot.Lines, Line{
			Id: p.plot.nextNegativeLineId(),
			C:  []complex128{pi.C},
			Style: DataStyle{
				Marker: MarkerStyle{
					Marker: CrossMarker,
					Color:  -1,
					Size:   3,
				},
			},
		})
		if p.axes.limits.isPolarLimits() {
			p.draw(true)
		} else { // axes are zoomed.
			p.axes.drawPolarDataOnly()
		}
		return Callback{Type: MeasurePoint, PointInfo: pi}, ok
	}
	return Callback{PointInfo: pi}, ok
}

func (p polarPlot) highlight(id []HighlightID) *image.RGBA {
	if id != nil {
		a := p.axes
		a.highlight(id, a.xyRing())
		if a.limits.isPolarLimits() {
			a.drawPolarTics(p.ring)
			a.drawPolarCircle(p.ring)
		}
	}
	return p.im
}
