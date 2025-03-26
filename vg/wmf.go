package vg

import (
	"image"
	"image/color"
	"math"

	"github.com/ktye/plot/vg/wmf"
	"golang.org/x/image/math/fixed"
)

/*
 * no clip, no embedded png
 */

type Wmf struct {
	rect image.Rectangle
	*wmfcanvas
}

type wmfcanvas struct {
	wmf.File
	fg           wmf.Color
	bg           wmf.Color
	objects      int
	pens         map[wmf.Pen]int
	brushes      map[wmf.Brush]int
	currentPen   int
	currentBrush int
	currentFont  int
	fontSize     int
	textColor    wmf.Color
	textAlign    uint16
}

func NewWmf(w, h int) *Wmf {
	var f Wmf
	var c wmfcanvas
	f.rect = image.Rect(0, 0, w, h)
	c.File = *wmf.New(w, h)
	f.wmfcanvas = &c
	f.pens = make(map[wmf.Pen]int)
	f.brushes = make(map[wmf.Brush]int)
	f.fontSize = 20

	f.SetBkMode(1)                                                                    //transparent
	f.CreatePen(wmf.Pen{Style: 5})                                                    //invisible
	f.CreateBrush(wmf.Brush{Style: 1})                                                //hollow
	f.CreateFont(wmf.Font{Height: -16, Face: "Calibri", OutPrecision: 7, Quality: 4}) //OutPrecision:7, Quality:4
	f.CreateFont(wmf.Font{Height: -12, Face: "Calibri"})
	f.currentPen = -1
	f.currentBrush = -1
	f.objects = 4

	return &f
}
func (f *Wmf) setFillStroke(fill bool, lw int) {
	if fill {
		brush := wmf.Brush{Color: f.fg}
		b, o := f.brushes[brush]
		if o == false {
			b = f.objects
			f.objects++
			f.brushes[brush] = b
			f.CreateBrush(brush)
		}
		f.setCurrentPen(0)
		f.setCurrentBrush(b)
	} else {
		pen := wmf.Pen{Width: uint16(lw), Color: f.fg}
		p, o := f.pens[pen]
		if o == false {
			p = f.objects
			f.objects++
			f.pens[pen] = p
			f.CreatePen(pen)
		}
		f.setCurrentBrush(1)
		f.setCurrentPen(p)
	}
}
func (f *Wmf) setCurrentPen(x int) {
	if f.currentPen != x {
		f.currentPen = x
		f.Select(x)
	}
}
func (f *Wmf) setCurrentBrush(x int) {
	if f.currentBrush != x {
		f.currentBrush = x
		f.Select(x)
	}
}
func (f *Wmf) setFont(x int) {
	if f.currentFont != x {
		f.currentFont = x
		if x == 2 {
			f.fontSize = 20
		} else {
			f.fontSize = 16
		}
		f.Select(x)
	}
}
func (f *Wmf) setTextColor() {
	if f.fg != f.textColor {
		f.textColor = f.fg
		f.SetTextColor(f.textColor)
	}
}
func (f *Wmf) setAlign(a uint16) {
	if a != f.textAlign {
		f.textAlign = a
		f.SetTextAlign(a)
	}
}
func (f *Wmf) Reset()           {}
func (f *Wmf) Size() (int, int) { return f.rect.Dx(), f.rect.Dy() }
func (f *Wmf) SubImage(r image.Rectangle) Drawer {
	s := Wmf{
		rect:      r.Add(f.rect.Min),
		wmfcanvas: f.wmfcanvas,
	}
	return &s
}
func (f *Wmf) Bounds() image.Rectangle { return f.rect }
func (f *Wmf) Paint()                  {}
func wmfColor(c color.Color) wmf.Color {
	r, g, b, _ := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	return wmf.Color(uint32(b) | uint32(g)<<8 | uint32(r)<<16)
}
func (f *Wmf) Clear(c color.Color) { f.bg = wmfColor(c) }
func (f *Wmf) Color(c color.Color) { f.fg = wmfColor(c) }
func (f *Wmf) Line(l Line) {
	f.setFillStroke(false, l.LineWidth)
	x, y, dx, dy := l.X+f.rect.Min.X, l.Y+f.rect.Min.Y, l.DX, l.DY
	f.MoveTo(int16(x), int16(y))
	f.LineTo(int16(x+dx), int16(y+dy))
}
func (f *Wmf) Circle(c Circle) {
	f.setFillStroke(c.Fill, c.LineWidth)
	x, y, d := c.X+f.rect.Min.X, c.Y+f.rect.Min.Y, c.D
	f.Ellipse(int16(x), int16(y), int16(x+d), int16(y+d))
}
func (f *Wmf) Rectangle(r Rectangle) {
	f.setFillStroke(r.Fill, r.LineWidth)
	x, y := r.X+f.rect.Min.X, r.Y+f.rect.Min.Y
	f.File.Rectangle(int16(x), int16(y), int16(x+r.W), int16(y+r.H))
}
func (f *Wmf) Triangle(t Triangle) {
	f.setFillStroke(t.Fill, t.LineWidth)
	x0, y0 := f.rect.Min.X, f.rect.Min.Y
	x := []int16{int16(t.X0 + x0), int16(t.X1 + x0), int16(t.X2 + x0), int16(t.X0 + x0)}
	y := []int16{int16(t.Y0 + y0), int16(t.Y1 + y0), int16(t.Y2 + y0), int16(t.Y0 + y0)}
	f.Polygon(x, y)
}
func (f *Wmf) Ray(r Ray) {
	f.setFillStroke(false, r.LineWidth)
	xo, yo := f.rect.Min.X, f.rect.Min.Y
	x0 := float64(r.X+xo) + float64(r.R)*math.Cos(r.Phi) + 0.5
	y0 := float64(r.Y+yo) + float64(r.R)*math.Sin(r.Phi) + 0.5
	x1 := float64(r.X+xo) + float64(r.R+r.L)*math.Cos(r.Phi) + 0.5
	y1 := float64(r.Y+yo) + float64(r.R+r.L)*math.Sin(r.Phi) + 0.5
	f.MoveTo(int16(x0), int16(y0))
	f.LineTo(int16(x1), int16(y1))
}
func (f *Wmf) Text(t Text) {
	//fmt.Printf(":text %+v\n", t)
	//f.setAlign([]uint16{8, 6 + 8, 2 + 8, 2 + 24, 2, 6, 0, 24, 6 + 24}[t.Align])
	h, h2 := int16(f.fontSize), int16(f.fontSize/2)
	f.setAlign([]uint16{0, 6, 2, 2, 2, 6, 0, 0, 6}[t.Align])
	y0 := []int16{h, h, h, h2, 0, 0, 0, h2, h2}[t.Align]
	f.setTextColor()
	x, y := t.X+f.rect.Min.X, t.Y+f.rect.Min.Y
	f.File.Text(int16(x), int16(y)-y0, t.S)
}
func (f *Wmf) Font(font1 bool) {
	if font1 {
		f.setFont(2)
	} else {
		f.setFont(3)
	}
}
func (f *Wmf) ArrowHead(a ArrowHead) {
	if len(a.X) < 2 || len(a.Y) < 2 {
		return
	}
	x, y, xa, ya, xb, yb, xc, yc := a.points(rect26_6(f.rect))
	X := []int16{int16(x >> 6), int16(xb), int16(xa), int16(xc), int16(x >> 6)}
	Y := []int16{int16(y >> 6), int16(yb), int16(ya), int16(yc), int16(y >> 6)}
	f.setFillStroke(true, 0)
	f.Polygon(X, Y)
}
func (f *Wmf) FloatTics(t FloatTics) {
	f.setFillStroke(false, t.LineWidth)
	for _, v := range t.P {
		X, Y := t.Q, v
		if t.Horizontal {
			X, Y = Y, X
		}
		rect := rect26_6(t.Rect)
		x, y := transform(X, Y, t.CoordinateSystem, rect)
		off := fixed.I(t.L)
		if t.LeftTop == false {
			off = -off
		}
		if t.Horizontal {
			y -= off
		} else {
			x -= off
		}
		f.MoveTo(int16(x>>6), int16(y>>6))
		if t.Horizontal {
			y += off
		} else {
			x += off
		}
		f.LineTo(int16(x>>6), int16(y>>6))
	}
}
func (f *Wmf) FloatText(t FloatText)                            { f.Text(t.toText(f.rect.Min.X, f.rect.Min.Y)) }
func (f *Wmf) FloatTextExtent(t FloatText) (int, int, int, int) { return 0, 0, 0, 0 } //only used interactively
func (f *Wmf) FloatBars(b FloatBars) {
	f.setFillStroke(true, 0)
	for i := 0; i < len(b.X); i += 2 {
		x0, y0 := transform(b.X[i], b.Y[i], b.CoordinateSystem, rect26_6(f.rect))
		x1, y1 := transform(b.X[i+1], b.Y[i+1], b.CoordinateSystem, rect26_6(f.rect))
		f.File.Rectangle(int16(x0>>6), int16(y0>>6), int16(x1>>6), int16(y1>>6))
	}
}
func (f *Wmf) FloatCircles(c FloatCircles) {
	f.setFillStroke(c.Fill, c.LineWidth)
	for i := range c.X {
		x, y := transform(c.X[i], c.Y[i], c.CoordinateSystem, rect26_6(f.rect))
		x += fixed.I(c.Z)
		y -= fixed.I(c.Z)
		x -= fixed.I(c.Radius)
		y -= fixed.I(c.Radius)
		xi := int16(x >> 6)
		yi := int16(y >> 6)
		d := int16(1 + 2*c.Radius)
		f.Ellipse(xi, yi, xi+d, yi+d)
	}
}
func (f *Wmf) FloatEnvelope(e FloatEnvelope) {
	f.setFillStroke(true, e.LineWidth)
	var X, Y []int16
	z := fixed.I(e.Z)
	r := rect26_6(f.rect)
	for i := range e.X {
		x, y := transform(e.X[i], e.Y[i], e.CoordinateSystem, r)
		x += z
		y -= z
		X = append(X, int16(x>>6))
		Y = append(Y, int16(y>>6))
	}
	f.Polygon(X, Y)
}
func (f *Wmf) FloatPath(p FloatPath) {
	f.setFillStroke(false, p.LineWidth)
	var X, Y []int16
	flush := func() {
		if len(X) > 1 {
			f.Polyline(X, Y)
		}
		X, Y = nil, nil
	}
	for i := range p.X {
		if math.IsNaN(p.X[i]) || math.IsNaN(p.Y[i]) {
			flush()
			continue
		}
		x, y := transform(p.X[i], p.Y[i], p.CoordinateSystem, rect26_6(f.rect))
		z := fixed.I(p.Z)
		x += z
		y -= z
		X = append(X, int16(x>>6))
		Y = append(Y, int16(y>>6))
	}
	flush()
}
