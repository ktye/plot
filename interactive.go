package plot

import (
	"image"
)

// ClickCoords transforms from a click on the original (multi-row/column) image to the single image coordinates.
func ClickCoords(p Iplots, x, y int) (xn, yn, n int) {
	pt := image.Point{x, y}
	rect, n := p.g.rect(0), 0
	for i := range p.p {
		if r := p.g.rect(i); pt.In(r) {
			rect = r
			n = i
			break
		}
	}
	im := p.p[n].image()
	rect = p.g.center(rect, im)
	pt = pt.Sub(rect.Min)
	return pt.X, pt.Y, n
}

// LineIPlotters is the callback routine for drawing a line with the mouse
// on the image of multiple plots.
func LineIPlotters(p Iplots, x, y, x1, y1 int) (complex128, bool) {
	if len(p.p) == 0 {
		return complex(0, 0), false
	}
	dx, dy, n := x1-x, y1-y, 0
	x, y, n = ClickCoords(p, x, y)
	return p.p[n].line(x, y, x+dx, y+dy)
}

func Measure(p Iplots, x, y, x1, y1 int) (MeasureInfo, bool) { //MeasureInfo:callback.go
	if len(p.p) == 0 {
		return MeasureInfo{}, false
	}
	abs := func(x int) int { return max(x, -x) }
	dx, dy, n := x1-x, y1-y, 0
	x, y, n = ClickCoords(p, x, y)
	mi, ok := p.p[n].measure(x, y, x+dx, y+dy)
	mi.X0, mi.Y0, mi.X1, mi.Y1 = x, y, x1, y1
	mi.N = n
	mi.Vertical = abs(dy) > abs(dx)
	if mi.Polar == false { //force h/v
		if mi.Vertical {
			mi.B = complex(real(mi.A), imag(mi.B))
			dx = 0
		} else {
			mi.B = complex(real(mi.B), imag(mi.A))
			dy = 0
		}
	}
	if c, ok := p.p[n].click(x, y, true, false, false); ok {
		mi.AA = c.PointInfo
	}
	if c, ok := p.p[n].click(x+dx, y+dy, true, false, false); ok {
		mi.BB = c.PointInfo
	}
	return mi, ok
}
func Annotate(p Iplots, m MeasureInfo, label string, circle, color, lw int) {
	type drawer interface {
		draw()
	}
	var plt *Plot
	var drw drawer
	var aa bool
	switch v := p.p[m.N].(type) {
	case polarPlot:
		plt, drw = v.plot, v
	case xyPlot:
		plt, drw = v.plot, v
	case ampAngPlot:
		plt, drw, aa = v.plot, v, true
	default:
		return
	}
	if m.Polar {
		arrow := 3 + 2*lw
		if circle > 0 {
			arrow = 0
		}
		plt.Lines = append(plt.Lines, Line{
			Id:    plt.nextNegativeLineId(),
			C:     []complex128{m.A, m.B},
			Style: DataStyle{Line: LineStyle{Width: lw, Color: color, Arrow: arrow, Circle: circle}},
			Label: label,
		})
	} else {
		var C []complex128
		Y := []float64{imag(m.A), imag(m.B)}
		if aa {
			C, Y = []complex128{complex(imag(m.A), 0), complex(imag(m.B), 0)}, nil
		}
		plt.Lines = append(plt.Lines, Line{
			Id:    plt.nextNegativeLineId(),
			X:     []float64{real(m.A), real(m.B)},
			Y:     Y,
			C:     C,
			Style: DataStyle{Line: LineStyle{Width: lw, Color: color, EndMarks: 2 + lw}},
			Label: label,
		})
	}
	drw.draw()
}

// ZoomIPlotters is the callback routine for dragging a rectange with a mouse
// on the image of multiple plots.
func ZoomIPlotters(p Iplots, x, y, dx, dy int) (bool, int) {
	if len(p.p) == 0 {
		return false, 0
	}
	var n int
	x, y, n = ClickCoords(p, x, y)
	return p.p[n].zoom(x, y, dx, dy), n
}

// PanIPlotters is the callback routine for dragging a rectange with a mouse
// on the image of multiple plots.
func PanIPlotters(p Iplots, x, y, dx, dy int) (bool, int) {
	if len(p.p) == 0 {
		return false, 0
	}
	var n int
	x, y, n = ClickCoords(p, x, y)
	return p.p[n].pan(x, y, dx, dy), n
}

// ClickIPlotters is the callback routine for clicking on a line or point
// on the image of multiple plots.
func ClickIPlotters(p Iplots, x, y int, snapToPoint, deleteLine, dodraw bool) (Callback, bool) {
	if len(p.p) == 0 {
		return Callback{}, false
	}
	var n int
	x, y, n = ClickCoords(p, x, y)
	callback, ok := p.p[n].click(x, y, snapToPoint, deleteLine, dodraw)
	callback.PlotIndex = n
	return callback, ok
}

func LimitsIPlotters(ip Iplots) []Limits {
	if len(ip.p) == 0 {
		return nil
	}
	limits := make([]Limits, len(ip.p))
	for i, p := range ip.p {
		limits[i] = p.limits()
	}
	return limits
}
