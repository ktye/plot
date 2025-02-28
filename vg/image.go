package vg

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
)

// Image implements Drawer and renders on a go image.
type Image struct {
	RGBA *image.RGBA
	w, h int
	p    *Painter
}

func NewImage(w, h int) *Image {
	m := Image{w: w, h: h, RGBA: image.NewRGBA(image.Rect(0, 0, w, h))}
	m.p = NewPainter(m.RGBA)
	return &m
}
func (m *Image) Size() (int, int) { return m.w, m.h }
func (m *Image) SubImage(x, y, w, h int) Drawer {
	s := Image{w: w, h: h, RGBA: m.RGBA.SubImage(image.Rect(x, y, x+w, y+h)).(*image.RGBA)}
	s.p = NewPainter(s.RGBA)
	return &s
}
func (m *Image) Paint() { m.p.Paint() }
func (m *Image) Clear(bg color.Color) {
	draw.Draw(m.RGBA, m.RGBA.Bounds(), image.NewUniform(bg), image.ZP, draw.Src)
}
func (m *Image) Color(c color.Color)   { m.p.SetColor(c) }
func (m *Image) Line(l Line)           { m.p.Add(l) }
func (m *Image) Circle(c Circle)       { m.p.Add(c) }
func (m *Image) Triangle(t Triangle)   { m.p.Add(t) }
func (m *Image) Ray(r Ray)             { m.p.Add(r) }
func (m *Image) Text(t Text)           { m.p.Add(t) }
func (m *Image) Font(f font.Face)      { m.p.SetFont(f) }
func (m *Image) ArrowHead(a ArrowHead) { m.p.Add(a) }

func (m *Image) FloatTics(t FloatTics)         { m.p.Add(t) }
func (m *Image) FloatText(t FloatText)         { m.p.Add(t) }
func (m *Image) FloatBars(b FloatBars)         { m.p.Add(b) }
func (m *Image) FloatCircles(c FloatCircles)   { m.p.Add(c) }
func (m *Image) FloatEnvelope(e FloatEnvelope) { m.p.Add(e) }
func (m *Image) FloatPath(p FloatPath)         { m.p.Add(p) }
