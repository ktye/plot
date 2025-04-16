package vg

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
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
	clip       image.Rectangle
	fg, bg     uint32
	fn, f1, f2 uint8
	h1         int //height of font1
}

func NewEmf(w, h int, font string, f1, f2 int) *Emf {
	if font == "" {
		font = "Calibri"
	}
	if f1 == 0 {
		f1 = 18
	}
	if f2 == 0 {
		f2 = 14
	}
	f := Emf{rect: image.Rect(0, 0, w, h)}
	e := emf{File: *emfplus.New(w, h)}
	e.f1 = e.Font(8*int16(f1), font)
	e.f2 = e.Font(8*int16(f2), font)
	e.fn = e.f1
	e.h1 = f1
	f.emf = &e
	return &f
}
func i16(x int) int16 { return int16(x << 3) }

func (e *Emf) Reset()           {}
func (e *Emf) Size() (int, int) { return e.rect.Dx(), e.rect.Dy() }
func (e *Emf) SubImage(r image.Rectangle) Drawer {
	s := Emf{
		rect: r.Add(e.rect.Min),
		emf:  e.emf,
	}
	return &s
}
func (e *Emf) Rgba() *image.RGBA { return nil }
func (e *Emf) clip() {
	if r := e.rect; r != e.emf.clip {
		e.Clip(i16(r.Min.X)-8, i16(r.Min.Y)-8, i16(r.Dx())+16, i16(r.Dy())+16)
		e.emf.clip = r
	}
}
func (e *Emf) translate(x, y int16) (int16, int16) {
	return x + i16(e.rect.Min.X), y + i16(e.rect.Min.Y)
}
func (e *Emf) rect26() fixed.Rectangle26_6 { return rect26_6(e.rect) }
func (e *Emf) Bounds() image.Rectangle     { return e.rect }
func (e *Emf) Paint()                      {}
func emfColor(c color.Color) uint32 {
	r, g, b, a := c.RGBA()
	return (a&0xff00)<<16 | (r&0xff00)<<8 | (g & 0xff00) | b>>8
}
func (e *Emf) Clear(c color.Color) { e.bg = emfColor(c) }
func (e *Emf) Color(c color.Color) { e.fg = emfColor(c) }
func (e *Emf) Line(l Line) {
	e.clip()
	p := e.Pen(i16(l.LineWidth), e.fg)
	x0, y0 := e.translate(i16(l.X), i16(l.Y))
	x1, y1 := x0+i16(l.DX), y0+i16(l.DY)
	hp := int16(0)
	if l.Floor == false {
		hp = 4
	}
	e.DrawPolyline(p, false, []int16{x0 + hp, x1 + hp}, []int16{y0 + hp, y1 + hp})
}
func (e *Emf) Circle(c Circle) {
	e.clip()
	x, y, w, h := i16(c.X), i16(c.Y), i16(c.D), i16(c.D)
	x, y = e.translate(x, y)
	if c.Fill {
		e.FillEllipse(e.fg, x, y, w, h)
	} else {
		p := e.Pen(i16(c.LineWidth), e.fg)
		e.DrawEllipse(p, x, y, w, h)
	}
}
func (e *Emf) polygon(lw int, fill bool, x, y []int16) {
	if fill {
		e.FillPolygon(e.fg, x, y)
	} else {
		p := e.Pen(i16(lw), e.fg)
		e.DrawPolyline(p, true, x, y)
	}
}
func (e *Emf) Rectangle(r Rectangle) {
	e.clip()
	x, y := i16(r.X), i16(r.Y)
	w, h := i16(r.W), i16(r.H)
	x, y = e.translate(x, y)
	if r.Fill {
		e.FillRects(e.fg, []int16{x}, []int16{y}, []int16{w}, []int16{h})
	} else {
		p := e.Pen(int16(r.LineWidth), e.fg)
		e.DrawRects(p, []int16{x}, []int16{y}, []int16{w}, []int16{h})
	}
}
func (e *Emf) Triangle(t Triangle) {
	e.clip()
	x0, y0 := e.rect.Min.X, e.rect.Min.Y
	x := []int16{i16(t.X0 + x0), i16(t.X1 + x0), i16(t.X2 + x0), i16(t.X0 + x0)}
	y := []int16{i16(t.Y0 + y0), i16(t.Y1 + y0), i16(t.Y2 + y0), i16(t.Y0 + y0)}
	e.polygon(t.LineWidth, t.Fill, x, y)
}
func (e *Emf) Ray(r Ray) {
	e.clip()
	xo, yo := e.rect.Min.X, e.rect.Min.Y
	x0 := int16(8 * (float64(r.X+xo) + float64(r.R)*math.Cos(r.Phi) /*+ 0.5 */))
	y0 := int16(8 * (float64(r.Y+yo) + float64(r.R)*math.Sin(r.Phi) /*+ 0.5 */))
	x1 := int16(8 * (float64(r.X+xo) + float64(r.R+r.L)*math.Cos(r.Phi) /*+ 0.5 */))
	y1 := int16(8 * (float64(r.Y+yo) + float64(r.R+r.L)*math.Sin(r.Phi) /*+ 0.5 */))
	p := e.Pen(i16(r.LineWidth), e.fg)
	e.DrawPolyline(p, false, []int16{x0, x1}, []int16{y0, y1})
}
func (e *Emf) Text(t Text) {
	e.clip()
	x, y := e.translate(i16(t.X), i16(t.Y))
	e.File.Text(x, y, t.S, e.fn, t.Align, t.Vertical, e.fg)
}
func (e *Emf) Font(f1 bool) {
	if f1 {
		e.fn = e.f1
	} else {
		e.fn = e.f2
	}
}
func (e *Emf) ArrowHead(a ArrowHead) {
	e.clip()
	if len(a.X) < 2 || len(a.Y) < 2 {
		return
	}
	x, y, xa, ya, xb, yb, xc, yc := a.points(e.rect26())
	X := []int16{int16(x >> 3), int16(8 * xb), int16(8 * xa), int16(8 * xc), int16(x >> 3)}
	Y := []int16{int16(y >> 3), int16(8 * yb), int16(8 * ya), int16(8 * yc), int16(y >> 3)}
	e.polygon(0, true, X, Y)
}

func (e *Emf) FloatTics(t FloatTics) {
	e.clip()
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
		x0 = append(x0, int16(x>>3))
		y0 = append(y0, int16(y>>3))
		if t.Horizontal {
			y += off
		} else {
			x += off
		}
		x1 = append(x1, int16(x>>3))
		y1 = append(y1, int16(y>>3))
	}
	p := e.Pen(int16(t.LineWidth), e.fg)
	e.LineSegments(p, x0, x1, y0, y1)
}
func (e *Emf) FloatText(t FloatText) { e.Text(t.toText(e.rect.Min.X, e.rect.Min.Y)) }
func (e *Emf) FloatTextExtent(t FloatText) (int, int, int, int) {
	q := t.toText(e.rect.Min.X, e.rect.Min.Y)
	x, y := i16(q.X), i16(q.Y)
	w := (2 * e.h1 * len(t.S)) / 3 // constant fw/fh of 2/3 is guessed to approx text extent
	x -= 8 * int16(([]int{0, 1, 2, 2, 2, 1, 0, 0, 1}[t.Align]*w)/2)
	y -= 8 * int16(([]int{2, 2, 2, 1, 0, 0, 0, 1, 1}[t.Align]*e.h1)/2)
	return int(x >> 3), int(y >> 3), w, e.h1
}
func (e *Emf) FloatBars(b FloatBars) {

	/*
		e.clip()
		x, y := i16(r.X), i16(r.Y)
		w, h := i16(r.W), i16(r.H)
		x, y = e.translate(x, y)
		if r.Fill {
			e.FillRects(e.fg, []int16{x}, []int16{y}, []int16{w}, []int16{h})
		} else {
			p := e.Pen(int16(r.LineWidth), e.fg)
			e.DrawRects(p, []int16{x}, []int16{y}, []int16{w}, []int16{h})
		}
	*/

	e.clip()
	var x, y, w, h []int16
	for i := 0; i < len(b.X); i += 2 {
		x0, y0 := transform(b.X[i], b.Y[i], b.CoordinateSystem, e.rect26())
		x1, y1 := transform(b.X[i+1], b.Y[i+1], b.CoordinateSystem, e.rect26())
		x = append(x, int16(x0>>3))
		y = append(y, int16(y1>>3))
		w = append(w, int16((x1-x0)>>3))
		h = append(h, int16((y0-y1)>>3))
	}
	if b.Fill {
		e.FillRects(e.fg, x, y, w, h)
	} else {
		p := e.Pen(int16(b.LineWidth), e.fg)
		e.DrawRects(p, x, y, w, h)
	}
}
func (e *Emf) FloatCircles(c FloatCircles) {
	e.clip()
	for i := range c.X {
		x, y := transform(c.X[i], c.Y[i], c.CoordinateSystem, e.rect26())
		x += fixed.I(c.Z)
		y -= fixed.I(c.Z)
		x -= fixed.I(c.Radius)
		y -= fixed.I(c.Radius)
		xi := int16(x >> 3)
		yi := int16(y >> 3)
		d := i16(1 + 2*c.Radius)
		if c.Fill {
			e.FillEllipse(e.fg, xi, yi, d, d)
		} else {
			p := e.Pen(i16(c.LineWidth), e.fg)
			e.DrawEllipse(p, xi, yi, d, d)
		}
	}
}
func (e *Emf) FloatEnvelope(f FloatEnvelope) {
	e.clip()
	var X, Y []int16
	r := e.rect26()
	z := fixed.I(f.Z)
	for i := range f.X {
		x, y := transform(f.X[i], f.Y[i], f.CoordinateSystem, r)
		x += z
		y -= z
		X = append(X, int16(x>>3))
		Y = append(Y, int16(y>>3))
	}
	e.FillPolygon(e.fg, X, Y)
	e.DrawPolyline(e.Pen(i16(f.LineWidth), e.fg), true, X, Y)
}
func (e *Emf) FloatPath(p FloatPath) {
	e.clip()
	var X, Y []int16
	pen := e.Pen(i16(p.LineWidth), e.fg)
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
		x, y := transform(p.X[i], p.Y[i], p.CoordinateSystem, e.rect26())
		z := fixed.I(p.Z)
		x += z
		y -= z
		xi, yi := int16(x>>3), int16(y>>3)
		X = append(X, xi)
		Y = append(Y, yi)
	}
	flush()
}
func (e *Emf) Embed(x, y int, pngdata []byte) {
	e.clip()
	m, err := png.Decode(bytes.NewReader(pngdata))
	if err != nil {
		return
	}
	x0, y0 := int16(x+e.rect.Min.X), int16(y+e.rect.Min.Y)
	w, h := m.Bounds().Dx(), m.Bounds().Dy()
	p := e.Png(w, h, pngdata)
	e.DrawImage(p, x0, y0, int16(w), int16(h))
}
