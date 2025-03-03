package plot

import "image"

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
func ClickIPlotters(p Iplots, x, y int, snapToPoint, deleteLine bool) (Callback, bool) {
	if len(p.p) == 0 {
		return Callback{}, false
	}
	var n int
	x, y, n = ClickCoords(p, x, y)
	callback, ok := p.p[n].click(x, y, snapToPoint, deleteLine)
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
