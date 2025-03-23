package vg

import (
	"fmt"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"math"
	"strings"
)

func NewSvg(w, h int) *Svg {
	s := svg{}
	S := Svg{r: image.Rectangle{Max: image.Point{w, h}}, svg: &s}
	return &S
}

type Svg struct {
	r   image.Rectangle
	svg *svg
}
type svg struct {
	fg color.Color
	bg color.Color
	f1 bool
	g  []string
}

func (f *Svg) push(s string)             { f.svg.g = append(f.svg.g, s) }
func (f *Svg) Reset()                    {}
func (f *Svg) rect() fixed.Rectangle26_6 { return fixed.R(0, 0, f.r.Dx(), f.r.Dy()) }
func (f *Svg) Size() (int, int)          { return f.r.Dx(), f.r.Dy() }
func (f *Svg) SubImage(r image.Rectangle) Drawer {
	s := Svg{r: r.Add(f.r.Min), svg: f.svg}
	return &s
}
func (f *Svg) Bounds() image.Rectangle { return f.r }
func (f *Svg) Clear(c color.Color)     { f.svg.bg = c }
func (f *Svg) Color(c color.Color)     { f.svg.fg = c }
func (f *Svg) rgb() string {
	r, g, b, _ := f.svg.fg.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}
func (f *Svg) fillStroke(lw int, fill bool) string {
	s := fmt.Sprintf("stroke-width='%d' strok='%s'", lw, f.rgb())
	if fill {
		s = fmt.Sprintf("fill='%s'", f.rgb())
	}
	return s
}
func (f *Svg) Line(l Line) {
	f.push(fmt.Sprintf("<line x1='%d' y1='%d' x2='%d' y2='%d' stroke-width='%d' stroke='%s'/>", l.X, l.Y, l.X+l.DX, l.Y+l.DY, l.LineWidth, f.rgb()))
}
func (f *Svg) Circle(c Circle) {
	f.push(fmt.Sprintf("<circle cx='%d' cy='%d' r='%.1f' %s>", c.X, c.Y, float64(c.D)/2, f.fillStroke(c.LineWidth, c.Fill)))
}
func (f *Svg) Rectangle(r Rectangle) {
	f.push(fmt.Sprintf("<rectangle x='%d' y='%d' width='%d' height='%d' %s>", r.X, r.Y, r.W, r.H, f.fillStroke(r.LineWidth, r.Fill)))
}
func (f *Svg) Triangle(t Triangle) {
	f.push(fmt.Sprintf("<polygon points='%d %d %d %d %d %d' %s>", t.X0, t.Y0, t.X1, t.Y1, t.X2, t.Y2, f.fillStroke(t.LineWidth, t.Fill)))
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
	b := "alignment-baseline='" + ([]string{"bottom", "bottom", "bottom", "center", "top", "top", "top", "center", "center"}[t.Align]) + "'"
	if a != "start" {
		a = "text-anchor='" + a + "'"
	}
	if t.Vertical {
		f.push(fmt.Sprintf("<g transform='translate(%d,%d) rotate(-90)'><text %s %s>%s</text></g>", t.X, t.Y, a, b, t.S))
	} else {
		f.push(fmt.Sprintf("<text x='%d' y='%d' %s %s>%s</text>", t.X, t.Y, a, b, t.S))
	}
}
func (f *Svg) Font(f1 bool) { f.svg.f1 = f1 }
func (f *Svg) ArrowHead(a ArrowHead) {
	x, y, xa, ya, xb, yb, xc, yc := a.points(f.rect())
	f.push(fmt.Sprintf("<polygon points='%d %d %d %d %d %d %d %d' %s/>", int(x>>6), int(y>>6), int(xa), int(ya), int(xb), int(yb), int(xc), int(yc), f.fillStroke(0, true)))
}

func (f *Svg) FloatTics(t FloatTics) {
	var p strings.Builder
	for _, v := range t.P {
		X, Y, dx, dy := t.Q, v, t.L, 0
		if t.Horizontal {
			X, Y, dx, dy = Y, X, dy, dx
		}
		x, y := transform(X, Y, t.CoordinateSystem, f.rect())
		fmt.Fprintf(&p, "M%d %d", int(x>>6)-dx, int(y>>6)-dy)
		if t.Horizontal {
			fmt.Fprintf(&p, "h%d", t.L)
		} else {
			fmt.Fprintf(&p, "v%d", t.L)
		}
	}
	f.push(fmt.Sprintf("<path d='%s' %s/>", p.String(), f.fillStroke(t.LineWidth, false)))
}
func (f *Svg) FloatText(t FloatText)                            { f.Text(t.toText(0, 0)) }
func (f *Svg) FloatTextExtent(t FloatText) (int, int, int, int) { return 0, 0, 0, 0 }
func (f *Svg) FloatBars(b FloatBars) {
	var p strings.Builder
	for i := 0; i < len(b.X); i += 2 {
		X0, Y0 := transform(b.X[i], b.Y[i], b.CoordinateSystem, f.rect())
		X1, Y1 := transform(b.X[i+1], b.Y[i+1], b.CoordinateSystem, f.rect())
		x0, y0, x1, y1 := int(X0<<6), int(Y0<<6), int(X1<<6), int(Y1<<6)
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
	for i := range e.X {
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
			fmt.Fprintf(&p, "M%d %d", int(x>>6), int(y>>6))
		} else {
			fmt.Fprintf(&p, "L%d %d", int(x>>6), int(y>>6))
		}
	}
	f.push(fmt.Sprintf("<path d='%s' %s/>", p.String(), f.fillStroke(q.LineWidth, false)))
}

func (f *Svg) Paint() {
	var b strings.Builder
	id := fmt.Sprintf("r-%d-%d-%d-%d", f.r.Min.X, f.r.Min.Y, f.r.Max.X, f.r.Max.Y)
	fmt.Fprintf(&b, "<g transform='translate(%d,%d)'>\n", f.r.Min.X, f.r.Min.Y)
	fmt.Fprintf(&b, "<defs><clipPath id='%s'><rect x='%d' y='%d' width='%d' height='%d'/></clipPath></defs\n", id, f.r.Min.X, f.r.Min.Y, f.r.Dx(), f.r.Dy())
	fmt.Fprintf(&b, "<g clip-path='url(#%s)'>\n", id)
	for _, s := range f.svg.g {
		fmt.Fprintln(&b, s)
	}
	fmt.Fprintf(&b, "</g></g>\n")
	f.svg.g = []string{b.String()}
}
func (f *Svg) File() string {
	var b strings.Builder
	fmt.Fprintf(&b, "<svg viewBox='0 0 %d %d' xmlns='http://www.w3.org/2000/svg'\n>", f.r.Dx(), f.r.Dy())
	for _, s := range f.svg.g {
		b.Write([]byte(s))
	}
	fmt.Fprintf(&b, "</svg>\n")
	return b.String()
}
