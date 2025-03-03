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
	m := Image{w: w, h: h}
	m.RGBA = image.NewRGBA(image.Rect(0, 0, m.w, m.h))
	m.Reset()
	return &m
}
func (m *Image) Reset()                  { m.p = NewPainter(m.RGBA) }
func (m *Image) Size() (int, int)        { return m.w, m.h }
func (m *Image) Bounds() image.Rectangle { return m.RGBA.Bounds() }
func (m *Image) SubImage(r image.Rectangle) Drawer {
	s := Image{w: r.Dx(), h: r.Dy(), RGBA: m.RGBA.SubImage(r).(*image.RGBA)}
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
func (m *Image) Rectangle(r Rectangle) { m.p.Add(r) }
func (m *Image) Triangle(t Triangle)   { m.p.Add(t) }
func (m *Image) Ray(r Ray)             { m.p.Add(r) }
func (m *Image) Text(t Text)           { m.p.Add(t) }
func (m *Image) Font(f font.Face)      { m.p.SetFont(f) }
func (m *Image) ArrowHead(a ArrowHead) { m.p.Add(a) }

func (m *Image) FloatTics(t FloatTics)                            { m.p.Add(t) }
func (m *Image) FloatText(t FloatText)                            { m.p.Add(t) }
func (m *Image) FloatTextExtent(t FloatText) (int, int, int, int) { return m.p.FloatTextExtent(t) }
func (m *Image) FloatBars(b FloatBars)                            { m.p.Add(b) }
func (m *Image) FloatCircles(c FloatCircles)                      { m.p.Add(c) }
func (m *Image) FloatEnvelope(e FloatEnvelope)                    { m.p.Add(e) }
func (m *Image) FloatPath(p FloatPath)                            { m.p.Add(p) }
