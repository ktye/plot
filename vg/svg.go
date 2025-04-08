package vg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"

	"golang.org/x/image/math/fixed"
)

func NewSvg(w, h int) *Svg {
	s := svg{}
	S := Svg{r: image.Rectangle{Max: image.Point{w, h}}, svg: &s}
	return &S
}

type Svg struct {
	r   image.Rectangle
	g   []string
	svg *svg
}
type svg struct {
	fg    color.Color
	bg    color.Color
	f1    bool
	clips map[string]string
	data  []string
}

func (f *Svg) push(s string)             { f.g = append(f.g, s) }
func (f *Svg) Reset()                    {}
func (f *Svg) rect() fixed.Rectangle26_6 { return fixed.R(0, 0, f.r.Dx(), f.r.Dy()) }
func (f *Svg) Size() (int, int)          { return f.r.Dx(), f.r.Dy() }
func (f *Svg) SubImage(r image.Rectangle) Drawer {
	if f.svg.clips == nil {
		f.svg.clips = make(map[string]string)
	}
	id := f.clipId(r)
	f.svg.clips[id] = fmt.Sprintf("<clipPath id='%s'><rect x='0' y='0' width='%d' height='%d'/></clipPath>", id, r.Dx(), r.Dy())
	s := Svg{r: r.Add(f.r.Min), svg: f.svg}
	return &s
}
func (f *Svg) Rgba() *image.RGBA { return nil }
func (f *Svg) Embed(x, y int, pngdata []byte) {
	s := "data:image/png;base64," + base64.StdEncoding.EncodeToString(pngdata)
	f.push(fmt.Sprintf("<image  xlink:href='%s'/>" /*x, y, m.Bounds().Dx(), m.Bounds().Dy(),*/, s))
}
func (f *Svg) clipId(r image.Rectangle) string { return fmt.Sprintf("r-%d-%d", r.Dx(), r.Dy()) }
func (f *Svg) Bounds() image.Rectangle         { return f.r }
func (f *Svg) Clear(c color.Color)             { f.svg.bg = c }
func (f *Svg) Color(c color.Color)             { f.svg.fg = c }
func (f *Svg) rgb() string {
	r, g, b, _ := f.svg.fg.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}
func (f *Svg) fillStroke(lw int, fill bool) string {
	s := fmt.Sprintf("stroke-width='%d' stroke='%s' fill='none'", lw, f.rgb())
	if fill {
		s = fmt.Sprintf("fill='%s'", f.rgb())
	}
	return s
}
func (f *Svg) Line(l Line) {
	f.push(fmt.Sprintf("<line x1='%d' y1='%d' x2='%d' y2='%d' stroke-width='%d' stroke='%s'/>", l.X, l.Y, l.X+l.DX, l.Y+l.DY, l.LineWidth, f.rgb()))
}
func (f *Svg) Circle(c Circle) {
	r := c.D / 2
	f.push(fmt.Sprintf("<circle cx='%d' cy='%d' r='%d' %s/>", c.X+r, c.Y+r, r, f.fillStroke(c.LineWidth, c.Fill)))
}
func (f *Svg) Rectangle(r Rectangle) {
	f.push(fmt.Sprintf("<rectangle x='%d' y='%d' width='%d' height='%d' %s/>", r.X, r.Y, r.W, r.H, f.fillStroke(r.LineWidth, r.Fill)))
}
func (f *Svg) Triangle(t Triangle) {
	f.push(fmt.Sprintf("<polygon points='%d %d %d %d %d %d' %s/>", t.X0, t.Y0, t.X1, t.Y1, t.X2, t.Y2, f.fillStroke(t.LineWidth, t.Fill)))
}
func (f *Svg) Ray(r Ray) {
	R := float64(r.R)
	x0 := r.X + int(math.Round(R*math.Cos(r.Phi)))
	y0 := r.Y + int(math.Round(R*math.Sin(r.Phi)))
	R += float64(r.L)
	x1 := r.X + int(math.Round(R*math.Cos(r.Phi)))
	y1 := r.Y + int(math.Round(R*math.Sin(r.Phi)))
	f.Line(NewLine(x0, y0, x1-x0, y1-y0, r.LineWidth))
}
func (f *Svg) Text(t Text) {
	a := []string{"start", "middle", "end", "end", "end", "middle", "start", "start", "middle"}[t.Align]
	b := "alignment-baseline='" + ([]string{"top", "top", "top", "middle", "hanging", "hanging", "hanging", "middle", "middle"}[t.Align]) + "'"
	Y := []int{3, 3, 3, 0, 0, 0, 0, 0, 0}[t.Align]
	if a != "start" {
		a = "text-anchor='" + a + "'"
	} else {
		a = ""
	}
	var class []string
	if f.svg.f1 == false {
		class = append(class, "s")
	}
	if t.Vertical {
		class = append(class, "v")
	}
	size := ""
	if len(class) > 0 {
		size = "class='" + strings.Join(class, " ") + "'"
	}
	//f.push(fmt.Sprintf("<g transform='translate(%d,%d) rotate(-90)'><text %s %s %s>%s</text></g>", t.X, t.Y, size, a, b, t.S))
	//f.push(fmt.Sprintf("<text x='%d' y='%d' %s %s %s style='writing-mode:sideways-lr'>%s</text>", t.X, t.Y, size, a, b, t.S))
	f.push(fmt.Sprintf("<text x='%d' y='%d' %s %s %s>%s</text>", t.X, t.Y-Y, size, a, b, t.S))
}
func (f *Svg) Font(f1 bool) { f.svg.f1 = f1 }
func (f *Svg) ArrowHead(a ArrowHead) {
	x, y, xa, ya, xb, yb, xc, yc := a.points(f.rect())
	f.push(fmt.Sprintf("<polygon points='%d %d %d %d %d %d %d %d' %s/>", int(x>>6), int(y>>6), int(xb), int(yb), int(xa), int(ya), int(xc), int(yc), f.fillStroke(0, true)))
}

func (f *Svg) FloatTics(t FloatTics) {
	var p strings.Builder
	rect := rect26_6(t.Rect.Sub(f.r.Min))
	for _, v := range t.P {
		X, Y, dx, dy := t.Q, v, t.L, 0
		if t.Horizontal {
			X, Y, dx, dy = Y, X, dy, dx
		}
		x, y := transform(X, Y, t.CoordinateSystem, rect)
		if t.LeftTop == false {
			dx, dy = 0, 0
		}
		fmt.Fprintf(&p, "M%d %d", int(x>>6)-dx, int(y>>6)-dy)
		if t.Horizontal {
			fmt.Fprintf(&p, "v%d", t.L)
		} else {
			fmt.Fprintf(&p, "h%d", t.L)
		}
	}
	f.push(fmt.Sprintf("<path d='%s' %s/>", p.String(), f.fillStroke(t.LineWidth, false)))
}
func (f *Svg) FloatText(ft FloatText)                           { f.Text(ft.toText(f.r.Min.X, f.r.Min.Y)) }
func (f *Svg) FloatTextExtent(t FloatText) (int, int, int, int) { return 0, 0, 0, 0 }
func (f *Svg) FloatBars(b FloatBars) {
	var p strings.Builder
	for i := 0; i < len(b.X); i += 2 {
		X0, Y0 := transform(b.X[i], b.Y[i], b.CoordinateSystem, f.rect())
		X1, Y1 := transform(b.X[i+1], b.Y[i+1], b.CoordinateSystem, f.rect())
		x0, y0, x1, y1 := int(X0>>6), int(Y0>>6), int(X1>>6), int(Y1>>6)
		fmt.Fprintf(&p, "M%d %dL%d %dL%d %dL%d %dz", x0, y0, x0, y1, x1, y1, x1, y0)
	}
	f.push(fmt.Sprintf("<path d='%s' %s/>", p.String(), f.fillStroke(0, true)))
}
func (f *Svg) FloatCircles(c FloatCircles) {
	for i := range c.X {
		x, y := transform(c.X[i], c.Y[i], c.CoordinateSystem, f.rect())
		x += fixed.I(c.Z)
		y -= fixed.I(c.Z)
		f.push(fmt.Sprintf("<circle cx='%d' cy='%d' r='%d' %s/>", int(x>>6), int(y>>6), c.Radius, f.fillStroke(c.LineWidth, c.Fill)))
	}
}
func (f *Svg) FloatEnvelope(e FloatEnvelope) {
	var p strings.Builder
	for i := range e.X { // maybe in reverse, like emfplus.go?
		x, y := transform(e.X[i], e.Y[i], e.CoordinateSystem, f.rect())
		x += fixed.I(e.Z)
		y -= fixed.I(e.Z)
		fmt.Fprintf(&p, "%d %d ", int(x>>6), int(y>>6))
	}
	f.push(fmt.Sprintf("<polygon points='%s' %s/>", p.String(), f.fillStroke(e.LineWidth, true)))
}
func (f *Svg) FloatPath(q FloatPath) {
	var p strings.Builder
	st := true
	for i := range q.X {
		if math.IsNaN(q.X[i]) || math.IsNaN(q.Y[i]) {
			st = true
			continue
		}
		x, y := transform(q.X[i], q.Y[i], q.CoordinateSystem, f.rect())
		z := fixed.I(q.Z)
		x += z
		y -= z
		if st {
			st = false
			fmt.Fprintf(&p, "M%d %d", int(x>>6), int(y>>6))
		} else {
			fmt.Fprintf(&p, "L%d %d", int(x>>6), int(y>>6))
		}
	}
	f.push(fmt.Sprintf("<path d='%s' %s/>", p.String(), f.fillStroke(q.LineWidth, false)))
}

func (f *Svg) Paint() {
	var b strings.Builder
	fmt.Fprintf(&b, "<g transform='translate(%d,%d)'>\n", f.r.Min.X, f.r.Min.Y)
	clip := false
	if len(f.svg.clips) > 0 {
		id := f.clipId(f.r)
		_, clip = f.svg.clips[id]
		if clip {
			fmt.Fprintf(&b, "<g clip-path='url(#%s)'>\n", id)
		}
	}
	for _, s := range f.g {
		fmt.Fprintln(&b, s)
	}
	if clip {
		fmt.Fprintf(&b, "</g>")
	}
	fmt.Fprintf(&b, "</g>\n")
	f.svg.data = append(f.svg.data, b.String())
}
func (f *Svg) Bytes() []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "<svg viewBox='0 0 %d %d' width='%d' height='%d' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink'>\n", f.r.Dx(), f.r.Dy(), f.r.Dx(), f.r.Dy())
	fmt.Fprintf(&b, "<style>text{font-family:Tahoma,sans-serif;font-size:16px}.s{font-size:12px}.v{writing-mode:sideways-lr;font-size:17px}</style>\n")

	if len(f.svg.clips) > 0 {
		fmt.Fprintf(&b, "<defs>\n")
		for _, s := range f.svg.clips {
			b.Write([]byte(s))
		}
		fmt.Fprintf(&b, "</defs>\n")
	}
	fmt.Fprintf(&b, "<g transform='translate(0.5,0.5)'>\n")
	for _, s := range f.svg.data {
		b.Write([]byte(s))
	}
	fmt.Fprintf(&b, "</g></svg>\n")
	return b.Bytes()
}
