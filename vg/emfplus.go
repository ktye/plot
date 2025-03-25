package vg

import (
	"image"
	"image/color"
	"math"

	"github.com/ktye/plot/vg/emfplus"
	"golang.org/x/image/math/fixed"
)

type Emf struct {
	rect image.Rectangle
	*emf
}

type emf struct {
	emfplus.File
	fg, bg     uint32
	fn, f1, f2 uint8
}

func NewEmf(w, h int) *Emf {
	f := Emf{rect: image.Rect(0, 0, w, h)}
	e := emf{File: *emfplus.New(w, h)}
	e.f1 = e.Font(16, "Calibri")
	e.f2 = e.Font(12, "Calibri")
	e.fn = e.f1
	f.emf = &e
	return &f
}

func (e *Emf) Reset()           {}
func (e *Emf) Size() (int, int) { return e.rect.Dx(), e.rect.Dy() }
func (e *Emf) SubImage(r image.Rectangle) Drawer {
	s := Emf{
		rect: r.Add(e.rect.Min),
		emf:  e.emf,
	}
	return &s
}
func (e *Emf) Bounds() image.Rectangle { return e.rect }
func (e *Emf) Paint()                  {}
func emfColor(c color.Color) uint32 {
	r, g, b, a := c.RGBA()
	return (a&0xff00)<<16 | (r&0xff00)<<8 | (g & 0xff00) | b>>8
}
func (e *Emf) Clear(c color.Color) { e.bg = emfColor(c) }
func (e *Emf) Color(c color.Color) { e.fg = emfColor(c) }
func (e *Emf) Line(l Line) {
	p := e.Pen(int16(l.LineWidth), e.fg)
	x0, y0 := int16(l.X), int16(l.Y)
	x1, y1 := x0+int16(l.DX), y0+int16(l.DY)
	e.DrawPolyline(p, false, []int16{x0, x1}, []int16{y0, y1})
}
func (e *Emf) Circle(c Circle) {
	left, top := int16(c.X), int16(c.Y)
	right, bottom := left+int16(c.D), top+int16(c.D)
	if c.Fill {
		p := e.Pen(int16(c.LineWidth), e.fg)
		e.DrawEllipse(p, left, top, right, bottom)
	} else {
		e.FillEllipse(e.fg, left, top, right, bottom)
	}
}
func (e *Emf) polygon(lw int, fill bool, x, y []int16) {
	if fill {
		p := e.Pen(int16(lw), e.fg)
		e.DrawPolyline(p, true, x, y)
	} else {
		e.FillPolygon(e.fg, x, y)
	}
}
func (e *Emf) Rectangle(r Rectangle) {
	x, y := int16(r.X), int16(r.Y)
	w, h := int16(r.W), int16(r.H)
	if r.Fill {
		e.FillRects(e.fg, []int16{x}, []int16{y}, []int16{w}, []int16{h})
	} else {
		p := e.Pen(int16(r.LineWidth), e.fg)
		e.DrawRects(p, []int16{x}, []int16{y}, []int16{w}, []int16{h})
	}
}
func (e *Emf) Triangle(t Triangle) {
	x0, y0 := e.rect.Min.X, e.rect.Min.Y
	x := []int16{int16(t.X0 + x0), int16(t.X1 + x0), int16(t.X2 + x0), int16(t.X0 + x0)}
	y := []int16{int16(t.Y0 + y0), int16(t.Y1 + y0), int16(t.Y2 + y0), int16(t.Y0 + y0)}
	e.polygon(t.LineWidth, t.Fill, x, y)
}
func (e *Emf) Ray(r Ray) {
	xo, yo := e.rect.Min.X, e.rect.Min.Y
	x0 := int16(float64(r.X+xo) + float64(r.R)*math.Cos(r.Phi) + 0.5)
	y0 := int16(float64(r.Y+yo) + float64(r.R)*math.Sin(r.Phi) + 0.5)
	x1 := int16(float64(r.X+xo) + float64(r.R+r.L)*math.Cos(r.Phi) + 0.5)
	y1 := int16(float64(r.Y+yo) + float64(r.R+r.L)*math.Sin(r.Phi) + 0.5)
	p := e.Pen(int16(r.LineWidth), e.fg)
	e.DrawPolyline(p, false, []int16{x0, x1}, []int16{y0, y1})
}
func (e *Emf) Text(t Text) {
	e.File.Text(int16(t.X), int16(t.Y), t.S, e.fn, t.Align, t.Vertical, e.fg)
}
func (e *Emf) Font(f1 bool) {
	if f1 {
		e.fn = e.f1
	} else {
		e.fn = e.f2
	}
}
func (e *Emf) ArrowHead(a ArrowHead) {
	if len(a.X) < 2 || len(a.Y) < 2 {
		return
	}
	x, y, xa, ya, xb, yb, xc, yc := a.points(rect26_6(e.rect))
	X := []int16{int16(x >> 6), int16(xb), int16(xa), int16(xc), int16(x >> 6)}
	Y := []int16{int16(y >> 6), int16(yb), int16(ya), int16(yc), int16(y >> 6)}
	e.polygon(0, true, X, Y)
}

func (e *Emf) FloatTics(t FloatTics) {
	var x0, y0, x1, y1 []int16
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
		x0 = append(x0, int16(x>>6))
		y0 = append(y0, int16(y>>6))
		if t.Horizontal {
			y += off
		} else {
			x += off
		}
		x1 = append(x1, int16(x>>6))
		y1 = append(y1, int16(y>>6))
	}
	p := e.Pen(int16(t.LineWidth), e.fg)
	e.LineSegments(p, x0, x1, y0, y1)
}
func (e *Emf) FloatText(t FloatText)                            { e.Text(t.toText(e.rect.Min.X, e.rect.Min.Y)) }
func (e *Emf) FloatTextExtent(t FloatText) (int, int, int, int) { return 0, 0, 0, 0 }
func (e *Emf) FloatBars(b FloatBars) {
	var x, y, w, h []int16
	for i := 0; i < len(b.X); i += 2 {
		x0, y0 := transform(b.X[i], b.Y[i], b.CoordinateSystem, rect26_6(e.rect))
		x1, y1 := transform(b.X[i+1], b.Y[i+1], b.CoordinateSystem, rect26_6(e.rect))
		x = append(x, int16(x0>>6))
		y = append(y, int16(y0>>6))
		w = append(w, int16((x1-x0)>>6))
		h = append(h, int16((y1-y0)>>6))
	}
	e.FillRects(e.fg, x, y, w, h)
}
func (e *Emf) FloatCircles(c FloatCircles) {
	for i := range c.X {
		x, y := transform(c.X[i], c.Y[i], c.CoordinateSystem, rect26_6(e.rect))
		x += fixed.I(c.Z)
		y -= fixed.I(c.Z)
		x -= fixed.I(c.Radius)
		y -= fixed.I(c.Radius)
		xi := int16(x >> 6)
		yi := int16(y >> 6)
		d := int16(1 + 2*c.Radius)
		if c.Fill {
			e.FillEllipse(e.fg, xi, yi, xi+d, yi+d)
		} else {
			p := e.Pen(int16(c.LineWidth), e.fg)
			e.DrawEllipse(p, xi, yi, xi+d, yi+d)
		}
	}
}
func (e *Emf) FloatEnvelope(f FloatEnvelope) {
	var X, Y []int16
	z := fixed.I(f.Z)
	r := rect26_6(e.rect)
	for i := range f.X {
		x, y := transform(f.X[i], f.Y[i], f.CoordinateSystem, r)
		x += z
		y -= z
		X = append(X, int16(x>>6))
		Y = append(Y, int16(y>>6))
	}
	e.FillPolygon(e.fg, X, Y)
}
func (e *Emf) FloatPath(p FloatPath) {
	var X, Y []int16
	pen := e.Pen(int16(p.LineWidth), e.fg)
	flush := func() {
		if len(X) > 1 {
			e.DrawPolyline(pen, false, X, Y)
		}
		X, Y = nil, nil
	}
	for i := range p.X {
		if math.IsNaN(p.X[i]) || math.IsNaN(p.Y[i]) {
			flush()
			continue
		}
		x, y := transform(p.X[i], p.Y[i], p.CoordinateSystem, rect26_6(e.rect))
		z := fixed.I(p.Z)
		x += z
		y -= z
		X = append(X, int16(x>>6))
		Y = append(Y, int16(y>>6))
	}
	flush()
}
