package plot

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/cmplx"

	"github.com/ktye/plot/vg"
)

// polarPlot is an implementation of the HiPlotter interface.
type polarPlot struct {
	plot   *Plot // underlying plot structure
	Limits       // computed axis limits
	ring   bool
	polarDimension
	axes   *axes
	drawer vg.Drawer
}

// Vertical layout
// - (border)
// - (vFill)
// - titleHeight
// - ticLabelHeight
// - polarDiameter
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
}

// Create a new polar coordinate system in the subimage.
// Width and height are available areas on input, the image (p.im) will be smaller.
func (plt *Plot) NewPolar(d vg.Drawer, isRing bool) (p polarPlot, err error) {
	width, height := d.Size()
	p.drawer = d
	p.plot = plt
	p.ring = isRing
	p.Limits = plt.getPolarLimits(p.ring)
	if p.Limits.Ymax-p.Limits.Ymin == 0 || math.IsNaN(p.Limits.Ymax) {
		p.Limits.Ymin, p.Limits.Ymax = 0, 1
		//return p, fmt.Errorf("cannot calculate polar limits (no data?)")
	}

	// Calculate Dimensions.
	border := plt.defaultBorder()
	p.ticLabelHeight = plt.defaultTicLabelHeight()
	p.ticLabelWidth = plt.defaultPolarTicLabelWidth()
	p.titleHeight = plt.defaultTitleHeight()

	// Needed space for decorations.
	hFix := func() int { return 2*border + 2*p.ticLabelWidth }
	vFix := func() int { return 2*border + p.titleHeight + 2*p.ticLabelHeight }

	// Available space for plotArea.
	hSpace := width - hFix()
	vSpace := height - vFix()

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
	if p.polarDiameter < 1 {
		return p, fmt.Errorf("image space is too small")
	}

	// Calculate (smaller) image dimensions.
	swidth := hFix() + p.polarDiameter
	sheight := vFix() + p.polarDiameter // p.titleHeight + 2*p.ticLabelHeight + 2*p.ticLength + p.polarDiameter
	x0 := (width - swidth) / 2
	y0 := (height - sheight) / 2
	ax := plt.newAxes(
		x0+p.ticLabelWidth+border,
		y0+p.titleHeight+p.ticLabelHeight+border,
		p.polarDiameter,
		p.polarDiameter,
		p.Limits,
		d,
	)
	p.axes = &ax

	p.draw()
	p.axes.store()
	return p, nil
}

func (p polarPlot) draw() {
	p.axes.reset()
	ccw := p.plot.Style.Counterclockwise
	p.axes.fillParentBackground()
	p.axes.drawPolar(p.ring, ccw)
	p.axes.drawTitle(p.ticLabelHeight)
	p.axes.inside.Paint()
	p.drawer.Paint()
}
func (p polarPlot) background() color.Color { return p.plot.defaultBackgroundColor() }
func (p polarPlot) image() *image.RGBA      { return p.drawer.(*vg.Image).RGBA }
func (p polarPlot) zoom(x, y, dx, dy int) bool {
	if dx > dy { // Keep it square.
		dy = dx
	} else {
		dx = dy
	}
	X0, Y0 := p.axes.toFloats(x, y+dy)
	X1, Y1 := p.axes.toFloats(x+dx, y)
	{ //remain centered, if offset is small
		r := max(math.Abs(Y1-Y0), math.Abs(X1-X0))
		o := cmplx.Abs(complex(0.5*(X1+X0), 0.5*(Y1+Y0)))
		if o < 0.1*r {
			X0, Y0, X1, Y1 = -r, -r, r, r
		}
	}
	p.axes.limits.Xmin = X0
	p.axes.limits.Xmax = X1
	p.axes.limits.Ymin = Y0
	p.axes.limits.Ymax = Y1
	p.axes.reset()
	p.draw()
	p.axes.store()
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
	p.axes.reset()
	p.draw()
	p.axes.store()
	return true
}

func (p polarPlot) limits() Limits {
	return p.axes.limits
}

func (p polarPlot) measure(x0, y0, x1, y1 int) (MeasureInfo, bool) {
	if !p.axes.isInside(x0, y0) || p.ring {
		return MeasureInfo{}, false
	}
	_, X0, Y0, X1, Y1 := p.axes.line(x0, y0, x1, y1)

	if p.plot.Style.Counterclockwise == false {
		X0, Y0 = Y0, X0
		X1, Y1 = Y1, X1
	}
	return MeasureInfo{A: complex(X0, Y0), B: complex(X1, Y1), Polar: true, Yunit: p.plot.Yunit}, true
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
		Style: DataStyle{Line: LineStyle{Width: 1, Color: -1, Arrow: 5}},
	})
	p.draw()
	return vec, true
}

func (p polarPlot) click(x, y int, snapToPoint, deleteLine, dodraw bool) (Callback, bool) {
	if !p.axes.isInside(x, y) {
		return Callback{
			Type:   AxisCallback,
			Limits: Limits{Xmax: p.axes.limits.Xmax, Ymax: p.axes.limits.Ymax},
		}, true
	}
	ccw := p.plot.Style.Counterclockwise
	pi, ok := p.axes.click(x, y, p.axes.xyRing(ccw), snapToPoint)
	if ok && snapToPoint == false {
		if p.ring {
			pi.C = complex(0, 0)
			pi.X = pi.X
			pi.Y = pi.Y
		} else {
			if ccw {
				pi.C = complex(pi.X, pi.Y)
			} else {
				pi.C = complex(pi.Y, pi.X)
			}
			pi.X = 0
			pi.Y = 0
		}
		if dodraw {
			if deleteLine {
				if n := len(p.plot.Lines); n > 0 {
					p.plot.Lines = p.plot.Lines[:n-1]
				}
			} else {
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
			}
			p.draw()
		}
		return Callback{Type: MeasurePoint, PointInfo: pi}, ok
	}
	return Callback{PointInfo: pi}, ok
}

func (p polarPlot) highlight(id []HighlightID) *image.RGBA {
	if id != nil {
		ccw := p.plot.Style.Counterclockwise
		a := p.axes
		a.restore()
		a.highlight(id, a.xyRing(ccw))
		//a.drawPolarTics(p.ring, ccw, a.limits.isPolarLimits() == false)
		//a.drawPolarCircle(p.ring)
	}
	return p.drawer.Rgba()
}
