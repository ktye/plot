package plot

// LineIPlotters is the callback routine for drawing a line with the mouse
// on the image of multiple plots.
func LineIPlotters(p []IPlotter, x, y, x1, y1, width, height int) (complex128, bool) {
	if len(p) == 0 {
		return complex(0, 0), false
	}
	w := width / len(p)
	n := x / w // n is single plot index.
	if n < 0 || n >= len(p) {
		return complex(0, 0), false
	}
	x %= w // x is now the offset within a single plot frame.
	im := p[n].image()
	bounds := im.Bounds()
	imwidth := bounds.Max.X - bounds.Min.X
	imheight := bounds.Max.Y - bounds.Min.Y
	yoff := (height - imheight) / 2
	xoff := (w - imwidth) / 2
	x -= xoff
	y -= yoff
	x1 %= w
	x1 -= xoff
	y1 -= yoff
	// x and y are now offsets within the single plot image.
	return p[n].line(x, y, x1, y1)
}

// ZoomIPlotters is the callback routine for dragging a rectange with a mouse
// on the image of multiple plots.
func ZoomIPlotters(p []IPlotter, x, y, dx, dy, width, height int) (bool, int) {
	if len(p) == 0 {
		return false, 0
	}
	w := width / len(p)
	n := x / w // n is single plot index.
	if n < 0 || n >= len(p) {
		return false, 0
	}
	if (x+dx)/w != n {
		return false, 0
	}
	x %= w // x is now the offset within a single plot frame.
	im := p[n].image()
	bounds := im.Bounds()
	imwidth := bounds.Max.X - bounds.Min.X
	imheight := bounds.Max.Y - bounds.Min.Y
	yoff := (height - imheight) / 2
	xoff := (w - imwidth) / 2
	x -= xoff
	y -= yoff
	// x and y are now offsets within the single plot image.
	return p[n].zoom(x, y, dx, dy), n
}

// PanIPlotters is the callback routine for dragging a rectange with a mouse
// on the image of multiple plots.
func PanIPlotters(p []IPlotter, x, y, dx, dy, width, height int) (bool, int) {
	if len(p) == 0 {
		return false, 0
	}
	w := width / len(p)
	n := x / w // n is single plot index.
	if n < 0 || n >= len(p) {
		return false, 0
	}
	x %= w // x is now the offset within a single plot frame.
	im := p[n].image()
	bounds := im.Bounds()
	imwidth := bounds.Max.X - bounds.Min.X
	imheight := bounds.Max.Y - bounds.Min.Y
	yoff := (height - imheight) / 2
	xoff := (w - imwidth) / 2
	x -= xoff
	y -= yoff
	// x and y are now offsets within the single plot image.
	return p[n].pan(x, y, dx, dy), n
}

// ClickIPlotters is the callback routine for clicking on a line or point
// on the image of multiple plots.
func ClickIPlotters(p []IPlotter, x, y, width, height int, snapToPoint bool) (Callback, bool) {
	if len(p) == 0 {
		return Callback{}, false
	}
	w := width / len(p)
	n := x / w // n is single plot index.
	if n < 0 || n >= len(p) {
		return Callback{}, false
	}
	x %= w // x is now the offset within a single plot frame.
	im := p[n].image()
	bounds := im.Bounds()
	imwidth := bounds.Max.X - bounds.Min.X
	imheight := bounds.Max.Y - bounds.Min.Y
	yoff := (height - imheight) / 2
	xoff := (w - imwidth) / 2
	x -= xoff
	y -= yoff
	// x and y are now offsets within the single plot image.
	callback, ok := p[n].click(x, y, snapToPoint)
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
