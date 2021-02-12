package plot

import "image"

// ClickCoords transforms from a click on the original (multi-row/column) image to the single image coordinates.
func ClickCoords(p []IPlotter, x, y, width, height, columns int) (xn, yn, n int) {
	g := newGrid(len(p), width, height, columns)
	pt := image.Point{x, y}
	rect, n := g.rect(0), 0
	for i := range p {
		if r := g.rect(i); pt.In(r) {
			rect = r
			n = i
			break
		}
	}
	im := p[n].image()
	rect = g.center(rect, im)
	pt = pt.Sub(rect.Min)
	return pt.X, pt.Y, n
}

// LineIPlotters is the callback routine for drawing a line with the mouse
// on the image of multiple plots.
func LineIPlotters(p []IPlotter, x, y, x1, y1, width, height, columns int) (complex128, bool) {
	if len(p) == 0 {
		return complex(0, 0), false
	}
	dx, dy, n := x1-x, y1-y, 0
	x, y, n = ClickCoords(p, x, y, width, height, columns)
	return p[n].line(x, y, x+dx, y+dy)
}

// ZoomIPlotters is the callback routine for dragging a rectange with a mouse
// on the image of multiple plots.
func ZoomIPlotters(p []IPlotter, x, y, dx, dy, width, height, columns int) (bool, int) {
	if len(p) == 0 {
		return false, 0
	}
	var n int
	x, y, n = ClickCoords(p, x, y, width, height, columns)
	return p[n].zoom(x, y, dx, dy), n
}

// PanIPlotters is the callback routine for dragging a rectange with a mouse
// on the image of multiple plots.
func PanIPlotters(p []IPlotter, x, y, dx, dy, width, height, columns int) (bool, int) {
	if len(p) == 0 {
		return false, 0
	}
	var n int
	x, y, n = ClickCoords(p, x, y, width, height, columns)
	return p[n].pan(x, y, dx, dy), n
}

// ClickIPlotters is the callback routine for clicking on a line or point
// on the image of multiple plots.
func ClickIPlotters(p []IPlotter, x, y, width, height, columns int, snapToPoint, deleteLine bool) (Callback, bool) {
	if len(p) == 0 {
		return Callback{}, false
	}
	var n int
	x, y, n = ClickCoords(p, x, y, width, height, columns)
	callback, ok := p[n].click(x, y, snapToPoint, deleteLine)
	callback.PlotIndex = n
	return callback, ok
}

func LimitsIPlotters(ip []IPlotter) []Limits {
	if len(ip) == 0 {
		return nil
	}
	limits := make([]Limits, len(ip))
	for i, p := range ip {
		limits[i] = p.limits()
	}
	return limits
}
