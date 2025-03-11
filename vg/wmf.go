package vg

import (
	"image"
	"image/color"

	"github.com/ktye/plot/vg/wmf"
	"golang.org/x/image/font"
)

type Wmf struct {
	rect image.Rectangle
	wmf.File
	fg           color.Color
	bg           color.Color
	fill, stroke bool
	objects      int
	pens         map[wmf.Pen]int
	brushes      map[wmf.Brush]int
	currentPen   int
	currentBrush int
}

func NewWmf(w, h int) *Wmf {
	var f Wmf
	f.rect = image.Rect(0, 0, w, h)
	f.File = *wmf.New(w, h)
	return &f
}
func (f *Wmf) setFillStroke(fill bool, lw int) {}
func (f *Wmf) Reset()                          {}
func (f *Wmf) Size() (int, int)                { return 0, 0 }
func (f *Wmf) SubImage(r image.Rectangle) Drawer {
	f.rect = r.Add(f.rect.Min)
	return f
}
func (f *Wmf) Bounds() image.Rectangle { return f.rect }
func (f *Wmf) Paint()                  {}
func (f *Wmf) Clear(c color.Color)     { f.bg = c }
func (f *Wmf) Color(c color.Color)     { f.fg = c }
func (f *Wmf) Line(Line)               {}
func (f *Wmf) Circle(c Circle) {
	f.setFillStroke(c.Fill, c.LineWidth)
	r := c.D / 2
	f.Ellipse(int16(c.X-r), int16(c.Y-r), int16(c.X+r), int16(c.Y-r))
}
func (f *Wmf) Rectangle(Rectangle)                            {}
func (f *Wmf) Triangle(Triangle)                              {}
func (f *Wmf) Ray(Ray)                                        {}
func (f *Wmf) Text(Text)                                      {}
func (f *Wmf) Font(font.Face)                                 {}
func (f *Wmf) ArrowHead(ArrowHead)                            {}
func (f *Wmf) FloatTics(FloatTics)                            {}
func (f *Wmf) FloatText(FloatText)                            {}
func (f *Wmf) FloatTextExtent(FloatText) (int, int, int, int) { return 0, 0, 0, 0 }
func (f *Wmf) FloatBars(FloatBars)                            {}
func (f *Wmf) FloatCircles(FloatCircles)                      {}
func (f *Wmf) FloatEnvelope(FloatEnvelope)                    {}
func (f *Wmf) FloatPath(FloatPath)                            {}
