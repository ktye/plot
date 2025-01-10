package plot

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strconv"

	"github.com/ktye/plot/raster"
	"github.com/ktye/plot/vg"
	"github.com/ktye/plot/xmath"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
)

// axes defines the position of the axes in the plot image.
type axes struct {
	x, y, width, height int    // upper left point, width and height (inclusive).
	zSpace              int    // x,y axes are shorter by this amount for xyz axes(zSpace>0)
	limits              Limits // axis limits.
	fg, bg              color.Color
	inside              *image.RGBA
	parent              *image.RGBA
	plot                *Plot
}

func (p *Plot) newAxes(x, y, width, height int, limits Limits, parent *image.RGBA) axes {
	a := axes{x: x, y: y, width: width, height: height}
	a.limits = limits
	a.inside = image.NewRGBA(image.Rect(0, 0, width, height))
	a.parent = parent
	a.plot = p
	a.fg = a.plot.defaultForegroundColor()
	a.bg = a.plot.defaultBackgroundColor()
	return a
}

// Fill complete image background, not just the axis space.
func (a axes) fillParentBackground() {
	draw.Draw(a.parent, a.parent.Bounds(), image.NewUniform(a.bg), image.ZP, draw.Src)
}

// DrawXY draws an xy coordinate system with data.
// It draws everything what is within the axes and puts the axes in the parent figure.
// This is used for xy plot, amp plot and ang plot.
func (a axes) drawXY(xy xyer) {
	// Draw to inside image.
	// Fill background color.
	draw.Draw(a.inside, a.inside.Bounds(), image.NewUniform(a.bg), image.ZP, draw.Src)

	// Draw data.
	p := vg.NewPainter(a.inside)
	a.drawLines(p, xy)
	p.Paint()

	// Put the inside image to the parent, but keep it for highlighting.
	draw.Draw(a.parent, image.Rect(a.x, a.y, a.x+a.width, a.y+a.height), a.inside, image.Point{0, 0}, draw.Src)
}

func (a axes) drawImage() {
	// Draw to inside image.
	// Fill background color.
	draw.Draw(a.inside, a.inside.Bounds(), image.NewUniform(a.bg), image.ZP, draw.Src)

	if len(a.plot.Lines) > 0 {
		l := a.plot.Lines[0]
		data := l.Image
		cols := len(data)
		if cols < 1 {
			goto Noimage
		}
		rows := len(data[0])
		if rows < 1 {
			goto Noimage
		}
		// Build an image with the dimensions of the data
		pal256 := a.plot.Style.Map.Palette()
		palette := make(color.Palette, 256)
		for i := 0; i < 256; i++ {
			palette[i] = pal256[i]
		}
		im := image.NewPaletted(image.Rect(0, 0, cols, rows), palette)
		for i := 0; i < cols; i++ {
			for k := 0; k < rows; k++ {
				im.SetColorIndex(i, k, data[i][rows-1-k])
			}
		}
		// Stretch data image to the inside image.
		x0 := int(xmath.Scale(a.limits.Xmin, l.X[0], l.X[len(l.X)-1], 0.0, float64(cols)))
		x1 := int(xmath.Scale(a.limits.Xmax, l.X[0], l.X[len(l.X)-1], 0.0, float64(cols)))
		y1 := int(xmath.Scale(a.limits.Ymin, l.Y[0], l.Y[len(l.Y)-1], float64(rows), 0.0))
		y0 := int(xmath.Scale(a.limits.Ymax, l.Y[0], l.Y[len(l.Y)-1], float64(rows), 0.0))
		xdraw.NearestNeighbor.Scale(a.inside, a.inside.Bounds(), im, image.Rect(x0, y0, x1, y1), xdraw.Src, nil)
		// draw.Draw(a.inside, a.inside.Bounds(), im, image.ZP, draw.Src)
	}
Noimage:

	// Put the inside image to the parent.
	draw.Draw(a.parent, image.Rect(a.x, a.y, a.x+a.width, a.y+a.height), a.inside, image.Point{0, 0}, draw.Src)
}

// DrawXYTics draws x and y tics with labels on the axes parent image.
func (a axes) drawXYTics(X, Y []float64, xlabels, ylabels []string) {
	p := vg.NewPainter(a.parent)
	p.SetColor(a.fg)

	// draw axis border lines
	boxLw := a.plot.defaultAxesGridLineWidth()
	aoff := a.plot.defaultTicLength()
	zw := a.zSpace
	p.Add(vg.Line{vg.LineCoords{a.x + zw, a.y - aoff, a.width - 1 - zw, 0}, boxLw, true})           //top
	p.Add(vg.Line{vg.LineCoords{a.x, a.y + a.height + aoff, a.width - 1 - zw, 0}, boxLw, true})     //bottom
	p.Add(vg.Line{vg.LineCoords{a.x - aoff, a.y + zw, 0, a.height - 1 - zw}, boxLw, true})          //left
	p.Add(vg.Line{vg.LineCoords{a.x + a.width + aoff - 1, a.y, 0, a.height - 1 - zw}, boxLw, true}) //right
	if zw > 0 {                                                                                     //diag (todo correct aoff)
		p.Add(vg.Line{vg.LineCoords{a.x, a.y + a.height + aoff, zw, -zw}, boxLw, true})                    //lower left
		p.Add(vg.Line{vg.LineCoords{a.x, a.y - aoff + zw, zw, -zw}, boxLw, true})                          //upper left
		p.Add(vg.Line{vg.LineCoords{a.x + a.width - 1 - zw, a.y + a.height + aoff, zw, -zw}, boxLw, true}) //lower right
	}

	// x and y tics on all 4 borders.
	L := a.plot.defaultTicLength()
	lw := a.plot.defaultAxesGridLineWidth()
	lim := a.limits
	cs := vg.CoordinateSystem{lim.Xmin, lim.Ymax, lim.Xmax, lim.Ymin}
	rect0 := image.Rect(a.x, a.y+zw, a.x+a.width-zw, a.y+a.height) //lower right
	rect1 := image.Rect(a.x+zw, a.y, a.x+a.width, a.y+a.height-zw) //upper left
	p.Add(vg.FloatTics{P: X, Q: a.limits.Ymin, Horizontal: true, LeftTop: false, L: L, LineWidth: lw, CoordinateSystem: cs, Rect: rect0})
	p.Add(vg.FloatTics{P: X, Q: a.limits.Ymax, Horizontal: true, LeftTop: true, L: L, LineWidth: lw, CoordinateSystem: cs, Rect: rect1})
	p.Add(vg.FloatTics{P: Y, Q: a.limits.Xmin, Horizontal: false, LeftTop: true, L: L, LineWidth: lw, CoordinateSystem: cs, Rect: rect0})
	p.Add(vg.FloatTics{P: Y, Q: a.limits.Xmax, Horizontal: false, LeftTop: false, L: L, LineWidth: lw, CoordinateSystem: cs, Rect: rect1})

	// Draw x tic labels if requested (todo aoff)
	textWidth := func(s string) int {
		return 7 * len(s) // Hard-coded.
	}
	p.SetFont(font2)
	var stop int
	for i, s := range xlabels {
		// Skip label, if it does not fit.
		if i > 0 {
			start, _ := cs.Pixel(X[i], lim.Ymin, rect0)
			start -= textWidth(s) / 2
			if start-stop < 3 {
				continue
			}
		}
		stop, _ = cs.Pixel(X[i], lim.Ymin, rect0)
		stop += textWidth(s) / 2
		yoff := font2.Metrics().Height.Ceil()
		p.Add(vg.FloatText{X: X[i], Y: lim.Ymin, S: s, Yoff: yoff, Align: 5, CoordinateSystem: cs, Rect: rect0})
	}

	// Draw y tic labels if requested
	xoff := -2 * L
	for i, s := range ylabels {
		p.Add(vg.FloatText{X: lim.Xmin, Y: Y[i], S: s, Yoff: 2, Xoff: xoff, Align: 3, CoordinateSystem: cs, Rect: rect0})
	}

	p.Paint()
}

// drawPolarDataOnly draws data for a polar diagram, but does not draw the diagram itself.
// It is used for interactive replot after zoom/pan.
func (a axes) drawPolarDataOnly(ccw bool) {
	draw.Draw(a.inside, a.inside.Bounds(), image.NewUniform(a.bg), image.ZP, draw.Src)
	p := vg.NewPainter(a.inside)
	a.drawLines(p, a.xyRing(ccw))
	p.Paint()
	draw.Draw(a.parent, image.Rect(a.x, a.y, a.x+a.width, a.y+a.height), a.inside, image.Point{0, 0}, draw.Src)
}

// DrawPolar draws the polar diagram background (x/y-lines), the data and the circle around.
// It draws everything what is within the axes and puts the axes in the parent figure.
func (a axes) drawPolar(ring, ccw bool) {
	// Draw to inside image.
	// Fill transparent background.
	draw.Draw(a.inside, a.inside.Bounds(), image.NewUniform(a.bg), image.ZP, draw.Src)

	// Draw data
	p := vg.NewPainter(a.inside)
	p.SetColor(a.fg)
	a.drawLines(p, a.xyRing(ccw))

	p.Paint()

	// Put the inside image to the parent, but keep it for highlighting.
	draw.Draw(a.parent, image.Rect(a.x, a.y, a.x+a.width, a.y+a.height), a.inside, image.Point{0, 0}, draw.Src)
	a.drawPolarCircle(ring)
}

func (a axes) xyRing(ccw bool) xyPolar {
	return xyPolar{rmin: a.limits.Zmin, rmax: a.limits.Zmax, ccw: ccw}
}

func (a axes) drawPolarCircle(ring bool) {
	// Draw x/y axis.
	p := vg.NewPainter(a.parent)
	p.SetColor(a.fg)
	r := a.width / 2
	lw := a.plot.defaultAxesGridLineWidth()
	if ring == false { // grid lines
		p.Add(vg.Line{vg.LineCoords{X: a.x + r, Y: a.y, DX: 0, DY: 2*r + 1}, lw, true})
		p.Add(vg.Line{vg.LineCoords{X: a.x, Y: a.y + r, DX: 2*r + 1, DY: 0}, lw, true})
	}

	// Draw grid circles
	asmin := 0.0
	if ring {
		asmin = a.limits.Zmin
	}
	as := autoscale{min: asmin, max: a.limits.Ymax}
	min, _, spacing := as.niceLimits()
	p.SetColor(color.Gray{128})
	rr := min
	if math.IsNaN(rr) == false { // rr may be NaN which results in an endless loop.
		for {
			if rr >= a.limits.Ymax {
				break
			}
			if rr > 0 && ring == false {
				w := int(float64(a.width) * rr / a.limits.Ymax)
				off := (a.width - w) / 2
				p.Add(vg.Circle{a.x + off, a.y + off, w, lw, false})
			} else if ring && rr > a.limits.Zmin {
				rrr := xmath.Scale(rr, a.limits.Zmin, a.limits.Zmax, innerRing, 1.0)
				w := int(float64(a.width) * rrr)
				off := (a.width - w) / 2
				p.Add(vg.Circle{a.x + off, a.y + off, w, lw, false})
			}
			rr += spacing
		}
	}

	p.SetColor(a.fg)
	lw = 2
	p.Add(vg.Circle{a.x, a.y, a.width, lw, false})
	if ring {
		off := int((float64(a.width) * innerRing) / 2.0)
		p.Add(vg.Circle{a.x + off, a.y + off, int(float64(a.width) * innerRing), lw, false})
	}
	p.Paint()
}

// DrawPolarTics draws tics, tic labels and the polar scale and unit.
func (a axes) drawPolarTics(ring, ccw bool) {
	// Draw Tics.
	p := vg.NewPainter(a.parent)
	p.SetColor(a.fg)
	r := a.width / 2
	if 2*r < a.width {
		r++ // Correct round off.
	}
	l := a.plot.defaultTicLength()
	aligns := []int{1, 0, 0, 7, 6, 6, 5, 4, 4, 3, 2, 2}
	p.SetFont(font2)
	phi0 := math.Pi / 2.0
	for i := 0; i < 360; i += 30 {
		phi := float64(i) * math.Pi / 180.0
		p.Add(vg.Ray{a.x + r, a.y + r, r - l/2, l, phi, a.plot.defaultAxesGridLineWidth()})
		s := strconv.Itoa(i)
		if ccw {
			s = strconv.Itoa((360 + 90 - i) % 360)
		}
		tx := float64(a.x+r) + float64(r+l/2)*math.Cos(phi-phi0)
		ty := float64(a.y+r) + float64(r+l/2)*math.Sin(phi-phi0)
		p.Add(vg.Text{int(tx + 0.5), int(ty + 0.5), s, aligns[i/30]})
	}

	// Draw scale/unit marker.
	p.SetFont(font1)
	phi := 130.0 * math.Pi / 180.0
	s := strconv.FormatFloat(a.limits.Ymax, 'g', 4, 64)
	p.Add(vg.Ray{a.x + r, a.y + r, r, 3 * l, phi - phi0, a.plot.defaultAxesGridLineWidth()})
	tx := float64(a.x+r) + float64(r+3*l)*math.Cos(phi-phi0) - float64(2*l)
	ty := float64(a.y+r) + float64(r+3*l)*math.Sin(phi-phi0)
	singleLine := false
	if ls := s + string(a.plot.Yunit); len(ls) < 6 {
		singleLine = true
		s = ls
	}
	p.Add(vg.Text{int(tx + 0.5), int(ty + 0.5), s, 6})
	if singleLine == false {
		ty += float64(3 + font2.Metrics().Height.Ceil()) // 3 should be line gap
		p.Add(vg.Text{int(tx + 0.5), int(ty + 0.5), string(a.plot.Yunit), 6})
	}
	if ring {
		s = strconv.FormatFloat(a.limits.Zmin, 'g', 4, 64)
		rmin := int(innerRing*float64(r)) - 3*l - 1
		p.Add(vg.Ray{a.x + r, a.y + r, rmin, 3 * l, phi - phi0, a.plot.defaultAxesGridLineWidth()})
		tx := float64(a.x+r) + float64(rmin)*math.Cos(phi-phi0)
		ty := float64(a.y+r) + float64(rmin)*math.Sin(phi-phi0)
		p.Add(vg.Text{int(tx - 0.5), int(ty - 0.5), s, 2})
	}
	p.Paint()
}

// drawTitle draws the title centered over the axes.
// between the axes top and the title bottom, the space is vSpace, which
// is usually 0 but needed for polar coordinate angle labels.
func (a axes) drawTitle(vSpace int) {
	p := vg.NewPainter(a.parent)
	p.SetColor(a.fg)
	p.SetFont(font1)
	// Center the label
	x := a.x + a.width/2
	y := a.y - vSpace - 1 // vSpace is defaultTicLength and more for polar.
	p.Add(vg.Text{x, y, a.plot.Title, 1})
	p.Paint()
}

// drawXLabel draws the x axis label.
func (a axes) drawXlabel() {
	p := vg.NewPainter(a.parent)
	p.SetColor(a.fg)
	p.SetFont(font1)
	x := a.x + a.width/2
	y := a.y + a.height + a.plot.defaultTicLabelHeight() //+ 3
	t := a.plot.Xlabel
	if a.plot.Xunit != "" {
		t += " " + string(a.plot.Xunit)
	}
	p.Add(vg.Text{x, y, t, 5})
	p.Paint()
}

// drawYLabel draws the rotated y axis label.
func (a axes) drawYlabel() {
	// Draw title to temporary image.
	f := font1
	t := a.plot.Ylabel
	if a.plot.Yunit != "" {
		t += " " + string(a.plot.Yunit)
	}

	bounds, _ := font.BoundString(f, t)
	d := bounds.Max.Sub(bounds.Min)
	width := d.X.Ceil() + 2
	metrics := f.Metrics()
	height := (metrics.Ascent + metrics.Descent).Ceil()

	tmp := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(tmp, tmp.Bounds(), image.NewUniform(a.bg), image.ZP, draw.Src)
	p := vg.NewPainter(tmp)
	p.SetColor(a.fg)
	p.SetFont(f)
	p.Add(vg.Text{1, height, t, 0})
	p.Paint()

	tmp = rotate(tmp)

	// Draw the temporary image into the parent.
	yoff := a.y + a.height/2 - width/2
	draw.Draw(a.parent, image.Rect(0, yoff, height, yoff+width), tmp, image.Point{0, 0}, draw.Src)
}

func scale3d(cs vg.CoordinateSystem) vg.CoordinateSystem {
	dx := (cs.X1 - cs.X0) * math.Sqrt2
	dy := (cs.Y1 - cs.Y0) * math.Sqrt2
	return vg.CoordinateSystem{cs.X0, cs.Y0 - dy, cs.X0 + dx, cs.Y1}
}

// drawLines draws plot's line data to the axes.
// It uses default styles, if no styles are defined.
func (a axes) drawLines(p *vg.Painter, xy xyer) {
	lim := a.limits
	cs := vg.CoordinateSystem{lim.Xmin, lim.Ymax, lim.Xmax, lim.Ymin} // upper left, lower right corner.
	if a.zSpace > 0 {
		cs = scale3d(cs)
	}
	z := 0.0
	for i, l := range a.plot.Lines {
		// Append vertical line data to lines separated by NaNs.
		for _, t := range l.V {
			l.X = append(l.X, math.NaN(), t, t)
			l.Y = append(l.Y, math.NaN(), lim.Ymin, lim.Ymax)
		}
		if n := len(a.plot.Lines); n > 0 {
			z = float64(i) / float64(n-1)
		}
		a.drawLine(p, xy, cs, l, z, false)
	}
}

func (a axes) drawLine(p *vg.Painter, xy xyer, cs vg.CoordinateSystem, l Line, z float64, isHighlight bool) {
	x, y, isEnvelope := xy.XY(l)

	// Set default style if unset.
	if l.Style.Marker.Size == 0 && l.Style.Line.Width == 0 {
		switch xy.(type) {
		case xyPolar:
			l.Style.Marker.Size = 3
		default:
			l.Style.Line.Width = 2
		}
	}
	var c color.Color = color.Black
	if size := l.Style.Marker.Size; size > 0.0 {
		c = a.plot.Style.Order.Get(l.Style.Marker.Color, l.Id+1).Color()
		if isHighlight {
			size *= 3
		}
		p.SetColor(c)
		if l.Style.Marker.Marker == Bar {
			p.Add(vg.FloatBars{X: x, Y: y, CoordinateSystem: cs})
		} else {
			p.Add(vg.FloatCircles{X: x, Y: y, CoordinateSystem: cs, Radius: size, Fill: true})
		}
	}
	if width := l.Style.Line.Width; width > 0.0 {
		c = a.plot.Style.Order.Get(l.Style.Line.Color, l.Id+1).Color()
		if isHighlight {
			width *= 3
		}
		if isEnvelope {
			p.SetColor(c)
			p.Add(vg.FloatEnvelope{X: x, Y: y, CoordinateSystem: cs, LineWidth: width})
		} else {
			if len(x) > 1000 {
				// Use a fast (non-antialiased) version if
				// there many points.
				im := raster.Image{
					Image: p.Image().(*image.RGBA),
					Color: c,
				}
				raster.FloatLines(im, x, y, raster.CoordinateSystem(cs))
			} else {
				p.SetColor(c)
				p.Add(vg.FloatPath{X: x, Y: y, CoordinateSystem: cs, LineWidth: width})
				if l.Style.Line.Arrow != 0 {
					p.Add(vg.ArrowHead{X: x, Y: y, CoordinateSystem: cs, LineWidth: width, Arrow: l.Style.Line.Arrow})
				}
			}
		}
	}
}

// drawSegment draws the given line segment.
func (a axes) drawSegment(p *vg.Painter, xy xyer, cs vg.CoordinateSystem, l Line, segment int) {
	// we modify l.X, and l.Y and restore it later.
	saveX := l.X
	saveY := l.Y
	saveC := l.C
	defer func() {
		l.X = saveX
		l.Y = saveY
		l.C = saveC
	}()

	// Get slice range for the given segment.
	x, _, _ := xy.XY(l)
	start, stop := 0, len(x)
	n := 0
	for i, f := range x {
		if math.IsNaN(f) {
			n++
			if n == segment {
				start = i + 1
			} else if n == segment+1 {
				stop = i
			}
		}
	}

	// What we acutally need to cut depends on the xyer.
	if start < len(l.X) && stop <= len(l.X) {
		l.X = l.X[start:stop]
	}
	if start < len(l.Y) && stop <= len(l.Y) {
		l.Y = l.Y[start:stop]
	}
	if start < len(l.C) && stop <= len(l.C) {
		l.C = l.C[start:stop]
	}

	a.drawLine(p, xy, cs, l, 0, false)
}

// drawDirectLine draws a single line directly on the axes parent image.
// It is used by StreamPlotter.
func (a axes) drawDirectLine(l Line, xy xyer, cs vg.CoordinateSystem) {
	r := image.Rectangle{image.Point{a.x, a.y}, image.Point{a.x + a.width, a.y + a.height}}
	im := a.parent.SubImage(r)
	/*
		m := raster.Image{
			Image: im.(*image.RGBA),
			Color: a.plot.Style.Color.Order.Get(l.Style.Line.Color, l.Id+1).Color(),
		}
		x, y, _ := xy.XY(l)
		raster.FloatLines(m, x, y, raster.CoordinateSystem(cs))
	*/
	p := vg.NewPainter(im.(*image.RGBA))
	a.drawLine(p, xy, cs, l, 0, false)
	p.Paint()
}

// drawPixel draws a single pixel directly on the parent image.
// It is used by StreamPlotter.
func (a axes) drawPixel(x, y float64, cs vg.CoordinateSystem, c color.Color) {
	r := image.Rectangle{image.Point{a.x, a.y}, image.Point{a.x + a.width, a.y + a.height}}
	xp, yp := cs.Pixel(x, y, r)
	a.parent.Set(xp, yp, c)
}

// drawPoint draws a highlighted point.
// pointNumber does not count NaN's.
func (a axes) drawPoint(p *vg.Painter, xy xyer, cs vg.CoordinateSystem, l Line, pointNumber int) {
	x, y, isEnvelope := xy.XY(l)

	// add number of NaNs leading pointNumber to pointNumber.
	targetNumber := pointNumber
	for i, v := range x {
		if i > targetNumber {
			break
		}
		if math.IsNaN(v) {
			pointNumber++
		}
	}

	if len(x) <= pointNumber || len(y) <= pointNumber || pointNumber < 0 {
		return
	}
	p.SetFont(font1)
	labels := make([]vg.FloatText, 2)
	if isEnvelope {
		if n := len(x); n != len(y) || pointNumber+2 > n {
			return
		} else {
			xp, yp := x[pointNumber], y[pointNumber]
			xp2, yp2 := x[n-pointNumber-2], y[n-pointNumber-2]
			x = []float64{xp, xp2}
			y = []float64{yp, yp2}
			labels[0] = vg.FloatText{X: xp, Y: yp, S: fmt.Sprintf("(%.4g, %.4g)", xp, yp), Align: 5}
			labels[1] = vg.FloatText{X: xp2, Y: yp2, S: fmt.Sprintf("(%.4g, %.4g)", xp2, yp2), Align: 1}
		}
	} else {
		xp, yp := x[pointNumber], y[pointNumber]
		x = []float64{xp}
		y = []float64{yp}
		var s string
		if xyp, ok := xy.(xyPolar); ok {
			xstr := ""
			if xyp.rmin == 0 && xyp.rmax == 0 { // polar
				if len(l.X) > pointNumber && pointNumber >= 0 {
					xstr = fmt.Sprintf("%.4g, ", l.X[pointNumber])
				}
				z := complex(yp, xp)
				if xyp.ccw {
					z = complex(xp, yp)
				}
				s = xstr + xmath.Absang(z, "%.4g@%.0f")
			} else { // ring
				s = fmt.Sprintf("%.4g@%.1f", l.X[pointNumber], 180.0*l.Y[pointNumber]/math.Pi)
			}
		} else {
			s = fmt.Sprintf("(%.4g, %.4g)", xp, yp)
		}
		labels[0] = vg.FloatText{X: xp, Y: yp, S: s, Align: 1}
		labels = labels[:1]
	}

	size := l.Style.Marker.Size
	if size == 0 {
		size = l.Style.Line.Width
	}
	if size == 0 {
		size = 9
	} else {
		size *= 3
	}
	c := a.plot.Style.Order.Get(l.Style.Marker.Color, l.Id+1).Color()
	p.SetColor(c)
	p.Add(vg.FloatCircles{X: x, Y: y, CoordinateSystem: cs, Radius: size, Fill: true})
	rect := a.inside.Bounds()
	for _, l := range labels {
		l.CoordinateSystem = cs
		l.Rect = rect

		// Change the alignment, if the label would be placed at a picture boundary.
		x0, y0 := cs.Pixel(l.X, l.Y, rect)
		if l.Align == 1 && y0 < 30 {
			l.Align = 5
		} else if l.Align == 5 && y0 > rect.Max.Y-30 {
			l.Align = 1
		}
		if x0 < 50 {
			if l.Align == 1 {
				l.Align = 0
			} else if l.Align == 5 {
				l.Align = 6
			}
		} else if x0 > rect.Max.X-50 {
			if l.Align == 1 {
				l.Align = 2
			} else if l.Align == 5 {
				l.Align = 4
			}
		}

		// Place the label above or below with the offset of the marker's radius.
		if l.Align <= 2 { // Label is above point.
			l.Yoff = -size
		} else if l.Align >= 4 { // Label is below point
			l.Yoff = size
		}

		// Fill background rectangle of the label.
		x, y, w, h := l.Extent(p)
		saveColor := p.GetColor()
		p.SetColor(a.bg)
		p.Add(vg.Rectangle{X: x, Y: y, W: w, H: h, Fill: true})
		p.SetColor(saveColor)
		p.Add(l)
	}
}

// line converts the line endpoints to coordinates.
func (a axes) line(x0, y0, x1, y1 int) (vec complex128, X0, Y0, X1, Y1 float64) {
	lim := a.limits
	cs := vg.CoordinateSystem{lim.Xmin, lim.Ymax, lim.Xmax, lim.Ymin}
	bounds := image.Rect(a.x, a.y, a.x+a.width, a.y+a.height)
	X0, Y0 = cs.Point(x0, y0, bounds)
	X1, Y1 = cs.Point(x1, y1, bounds)
	vec = complex(X1-X0, Y1-Y0)
	return
}

// Click returns the point info for the clicked point.
// If snapToPoint is true, it returns the info for the closest data point, otherwise
// for the clicked coordinate.
func (a axes) click(x, y int, xy xyer, snapToPoint bool) (PointInfo, bool) {
	// x, y := a.toFloats(xClick, yClick)
	lim := a.limits
	cs := vg.CoordinateSystem{lim.Xmin, lim.Ymax, lim.Xmax, lim.Ymin}
	bounds := image.Rect(a.x, a.y, a.x+a.width, a.y+a.height)

	if snapToPoint == false {
		px, py := cs.Point(x, y, bounds)
		return PointInfo{
			LineID:      -1,
			PointNumber: -1,
			NumPoints:   0,
			X:           px,
			Y:           py,
		}, true
	}

	dist := math.Inf(1)
	pIdx := -1
	lIdx := -1
	numPoints := 0
	isEnvelope := false
	maxSegment := 0
	isSegment := false
	for i, l := range a.plot.Lines {
		X, Y, isEnv := xy.XY(l)
		nNotNaN := -1
		segmentIdx := 0
		for n := range X {
			xi, yi := cs.Pixel(X[n], Y[n], bounds)
			// We only increase the index, if the data point is valid.
			nNotNaN++
			if math.IsNaN(X[n]) || math.IsNaN(Y[n]) {
				segmentIdx++
				if segmentIdx > maxSegment {
					maxSegment = segmentIdx
				}
				nNotNaN--
			}
			if d := float64((xi-x)*(xi-x) + (yi-y)*(yi-y)); d < dist {
				lIdx = i
				pIdx = nNotNaN
				isEnvelope = isEnv
				if isEnvelope {
					if n > len(X)/2 {
						pIdx = len(X) - n - 2
					}
				}
				dist = d

				numPoints = len(X)
				if l.Segments {
					pIdx = segmentIdx
					isSegment = true
				}
			}
		}
	}
	if lIdx < 0 || pIdx < 0 {
		return PointInfo{}, false
	}
	var px, py float64
	var pc complex128
	l := a.plot.Lines[lIdx]
	if len(l.X) > pIdx {
		px = l.X[pIdx]
	}
	if len(l.Y) > pIdx {
		py = l.Y[pIdx]
	}
	if len(l.C) > pIdx {
		pc = l.C[pIdx]
	}
	if isSegment {
		px = 0
		py = 0
		pc = complex(0, 0)
		numPoints = maxSegment + 1
	}
	return PointInfo{
		LineID:      l.Id,
		PointNumber: pIdx,
		NumPoints:   numPoints,
		IsEnvelope:  isEnvelope,
		X:           px,
		Y:           py,
		C:           pc,
	}, true
}

func (a axes) clickImage(xClick, yClick int) (PointInfo, bool) {
	x, y := a.toFloats(xClick, yClick)
	if len(a.plot.Lines) > 0 {
		i, _, z := a.plot.Lines[0].imageValueAt(x, y)
		var c complex128
		if i < len(a.plot.Lines[0].C) {
			c = a.plot.Lines[0].C[i]
		}
		zmin, zmax := a.plot.Lines[0].ImageMin, a.plot.Lines[0].ImageMax
		return PointInfo{X: x, Y: y, Z: z, C: c, IsImage: true, Zmin: zmin, Zmax: zmax}, true
	} else {
		return PointInfo{}, false
	}
}

func (a axes) isInside(x, y int) bool {
	if x < a.x || y < a.y {
		return false
	}
	if x > a.x+a.width || y > a.y+a.height {
		return false
	}
	return true
}

// convert from pixel coordinates to axis coordinates.
func (a axes) toFloats(x, y int) (X float64, Y float64) {
	X = xmath.Scale(float64(x), float64(a.x), float64(a.x+a.width), a.limits.Xmin, a.limits.Xmax)
	Y = xmath.Scale(float64(y), float64(a.y), float64(a.y+a.height), a.limits.Ymax, a.limits.Ymin)
	return
}

// highlight copies the initial inside image,
// plots over the original inside image with a highlighted line or data point.
// draws it over the parent and restores the inside image.
func (a *axes) highlight(ids []HighlightID, xy xyer) {
	// Backup axis inside image.
	backup := image.NewRGBA(a.inside.Bounds())
	draw.Draw(backup, a.inside.Bounds(), a.inside, image.Point{0, 0}, draw.Src)

	// Draw the highlighted line into the inside image.
	p := vg.NewPainter(a.inside)
	lim := a.limits
	cs := vg.CoordinateSystem{lim.Xmin, lim.Ymax, lim.Xmax, lim.Ymin} // upper left, lower right corner.

	if a.plot.Type == Raster {
		a.highlightImage(ids, xy, p)
	} else {
		// Drawing segments need to clear the background.
		// TODO: It would be nice to preserve the axis.
		clearBackground := false
		for _, id := range ids {
			if id.Point >= 0 {
				for _, l := range a.plot.Lines {
					if l.Segments {
						clearBackground = true
					}
				}
			}
		}
		if clearBackground {
			draw.Draw(a.inside, a.inside.Bounds(), image.NewUniform(a.bg), image.ZP, draw.Src)
		}

		z := 0.0
		for i, l := range a.plot.Lines {
			for _, id := range ids {
				if l.Id == -1 || l.Id == id.Line {
					if id.Point == -1 {
						if n := len(a.plot.Lines); n > 0 {
							z = float64(i) / float64(n-1)
						}
						a.drawLine(p, xy, cs, l, z, true)
					} else if l.Segments == true {
						a.drawSegment(p, xy, cs, l, id.Point)
					} else {
						a.drawPoint(p, xy, cs, l, id.Point)
					}
				}
			}
		}
	}
	p.Paint()

	// Draw the inside image into the parent.
	draw.Draw(a.parent, image.Rect(a.x, a.y, a.x+a.width, a.y+a.height), a.inside, image.Point{0, 0}, draw.Src)

	// Restore the non-highlighted inside image from the backup.
	a.inside = backup
}

func (a *axes) highlightImage(ids []HighlightID, xy xyer, ptr *vg.Painter) {
	if len(ids) != 1 {
		return
	}
	if len(a.plot.Lines) != 1 {
		return
	}
	if len(a.plot.Lines[0].Image) < 2 {
		return
	}

	id := ids[0]
	lim := a.limits
	cs := vg.CoordinateSystem{lim.Xmin, lim.Ymax, lim.Xmax, lim.Ymin} // upper left, lower right corner.
	line0 := a.plot.Lines[0]
	xf, yf := id.XImage, id.YImage
	xi, yi, _ := line0.imageValueAt(xf, yf)
	if xi < 0 || xi >= len(line0.Image) {
		return
	}
	if yi < 0 || yi >= len(line0.Image[0]) {
		return
	}
	l := Line{
		X: line0.X,
	}
	for i := 0; i < len(line0.Image); i++ {
		l.Y = append(l.Y, xmath.Scale(float64(line0.Image[i][yi]), 0, 255, lim.Ymin, 0.5*(lim.Ymin+lim.Ymax)))
	}

	// Draw horizontal line throught selection.
	a.drawLine(ptr, xy, cs, Line{
		X: []float64{lim.Xmin, lim.Xmax},
		Y: []float64{yf, yf},
	}, 0, false)

	// Draw vertical line throught selection.
	style2 := DataStyle{Line: LineStyle{Color: 2}}
	a.drawLine(ptr, xy, cs, Line{
		X:     []float64{xf, xf},
		Y:     []float64{lim.Ymin, lim.Ymax},
		Style: style2,
	}, 0, false)

	// Draw horizontal spectral line from the bottom to half the height.
	a.drawLine(ptr, xy, cs, l, 0, false)

	// Draw vertical spectral line from the right edge to half the width.
	l = Line{
		Y:     line0.Y,
		Style: style2,
	}
	for i := 0; i < len(line0.Image[xi]); i++ {
		l.X = append(l.X, xmath.Scale(float64(line0.Image[xi][i]), 0, 255, lim.Xmax, 0.5*(lim.Xmin+lim.Xmax)))
	}
	a.drawLine(ptr, xy, cs, l, 0, false)
}
