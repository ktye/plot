package vg

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/ktye/plot/vg/wmf"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

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

	f.CreatePen(wmf.Pen{Style: 5})     //invisible
	f.CreateBrush(wmf.Brush{Style: 1}) //hollow
	//f.CreateFont(wmf.Font{Height: 20, Face: "Arial"})
	//f.CreateFont(wmf.Font{Height: 16, Face: "Arial"})
	f.currentPen = -1
	f.currentBrush = -1
	f.objects = 2

	return &f
}
func (f *Wmf) setFillStroke(fill bool, lw int) {
	//fmt.Println("setFillStroke", "fill", fill, "fg", f.fg, "lw", lw)
	if fill {
		brush := wmf.Brush{Color: wmf.Blue}
		b, o := f.brushes[brush]
		if o == false {
			b = f.objects
			f.objects++
			f.brushes[brush] = b
			f.CreateBrush(brush)
			fmt.Println("new brush", b, f.fg)
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
			fmt.Println("new pen", p, f.fg)
		}
		f.setCurrentBrush(1)
		f.setCurrentPen(p)
	}
}
func (f *Wmf) setCurrentPen(x int) {
	if f.currentPen != x {
		f.currentPen = x
		fmt.Println("select pen", x)
		f.Select(x)
	}
}
func (f *Wmf) setCurrentBrush(x int) {
	if f.currentBrush != x {
		f.currentBrush = x
		fmt.Println("select brush", x)
		f.Select(x)
	}
}
func (f *Wmf) setFont(x int) {
	if f.currentFont != x {
		f.currentFont = x
		//f.Select(x)
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
func (f *Wmf) Reset() {
	fmt.Println(":reset")
}
func (f *Wmf) Size() (int, int) {
	//	fmt.Println(":size", f.rect.Dx(), f.rect.Dy())
	return f.rect.Dx(), f.rect.Dy()
}
func (f *Wmf) SubImage(r image.Rectangle) Drawer {
	s := Wmf{
		rect:      r.Add(f.rect.Min),
		wmfcanvas: f.wmfcanvas,
	}
	fmt.Println(":subimage", s.rect)
	return &s
}
func (f *Wmf) Bounds() image.Rectangle {
	//	fmt.Println(":bounds", f.rect)
	return f.rect
}
func (f *Wmf) Paint() {
	// fmt.Println(":paint")
}
func wmfColor(c color.Color) wmf.Color {
	r, g, b, _ := c.RGBA()
	r >>= 8
	g >>= 8
	b >>= 8
	return wmf.Color(uint32(b) | uint32(g)<<8 | uint32(r)<<16)
}
func (f *Wmf) Clear(c color.Color) {
	//	fmt.Println(":clear", c)
	f.bg = wmfColor(c)
}
func (f *Wmf) Color(c color.Color) {
	//	fmt.Println(":color", c)
	f.fg = wmfColor(c)
}
func (f *Wmf) Line(l Line) {
	fmt.Printf(":line %+v\n", l)
	f.setFillStroke(false, l.LineWidth)
	x, y, dx, dy := l.X+f.rect.Min.X, l.Y+f.rect.Min.Y, l.DX, l.DY
	f.MoveTo(x, y)
	f.LineTo(x+dx, y+dy)
}
func (f *Wmf) Circle(c Circle) {
	f.setFillStroke(c.Fill, c.LineWidth)
	fmt.Printf(":circle %+v\n", c)
	x, y, d := c.X+f.rect.Min.X, c.Y+f.rect.Min.Y, c.D
	f.Ellipse(int16(x), int16(y), int16(x+d), int16(y+d))
}
func (f *Wmf) Rectangle(r Rectangle) {
	// fmt.Printf(":rectangle %+v\n", r)
}
func (f *Wmf) Triangle(t Triangle) {
	// fmt.Printf(":triangle %+v\n", t)
}
func (f *Wmf) Ray(r Ray) {
	fmt.Printf(":ray %+v\n", r)
	f.setFillStroke(false, r.LineWidth)
	xo, yo := f.rect.Min.X, f.rect.Min.Y
	x0 := float64(r.X+xo) + float64(r.R)*math.Cos(r.Phi) - 0.5
	y0 := float64(r.Y+yo) + float64(r.R)*math.Sin(r.Phi) - 0.5
	x1 := float64(r.X+xo) + float64(r.R+r.L)*math.Cos(r.Phi) - 0.5
	y1 := float64(r.Y+yo) + float64(r.R+r.L)*math.Sin(r.Phi) - 0.5
	f.MoveTo(int(x0), int(y0))
	f.LineTo(int(x1), int(y1))
}
func (f *Wmf) Text(t Text) {
	fmt.Printf(":text %+v\n", t)
	f.setAlign([]uint16{8, 6 + 8, 2 + 8, 2 + 24, 2, 6, 0, 24, 6 + 24}[t.Align])
	f.setTextColor()
	x, y := t.X+f.rect.Min.X, t.Y+f.rect.Min.Y
	f.File.Text(int16(x), int16(y), t.S)
}
func (f *Wmf) Font(fn font.Face) {
	fmt.Printf(":font\n")
	f.setFont(2)
	//f.setFont(3) todo
}
func (f *Wmf) ArrowHead(a ArrowHead) {
	// fmt.Printf(":arrowhead %+v\n", a)
}
func (f *Wmf) FloatTics(t FloatTics) {
	// fmt.Printf(":floattics %+v\n", t)
}
func (f *Wmf) FloatText(t FloatText) {
	// fmt.Printf(":floattext %+v\n", t)
}
func (f *Wmf) FloatTextExtent(t FloatText) (int, int, int, int) {
	//	fmt.Printf(":floattextextent %+v\n", t)
	return 0, 0, 0, 0
}
func (f *Wmf) FloatBars(b FloatBars) {
	// fmt.Printf(":floatbars %+v\n", b)
}
func (f *Wmf) FloatCircles(c FloatCircles) {
	f.setFillStroke(c.Fill, c.LineWidth)
	fmt.Printf(":floatcicles #%d\n", len(c.X))
	for i := range c.X {
		x, y := transform(c.X[i], c.Y[i], c.CoordinateSystem, rect26_6(f.rect))
		x += fixed.I(c.Z)
		y -= fixed.I(c.Z)
		x -= fixed.I(c.Radius)
		y -= fixed.I(c.Radius)
		x >>= 6
		y >>= 6
		d := int16(2 * c.Radius)
		f.Ellipse(int16(x), int16(y), int16(x)+d, int16(y)+d)
	}
}
func (f *Wmf) FloatEnvelope(e FloatEnvelope) {
	// fmt.Printf(":floatenvelope %+v\n", e)
}
func (f *Wmf) FloatPath(p FloatPath) {
	fmt.Printf(":floatpath %+v\n", p)
}
