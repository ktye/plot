package vg

func NewSvg(w, h int) *Svg{
}

type Svg struct {
	r image.Rectangle
	g []string
	svg *svg
}
type svg struct {
	fg color.Color
	bg color.Color

}

func (f *Svg) Reset() {}
func (f *Svg) rect() image.Rectangle { return image.Rectangle{Max:image.Point{f.r.Dx(), f.r.Dy()}} }
func (f *Svg) Size() (int, int) { return f.r.Dx(), f.r.Dy() }
func (f *Svg) SubImage(r image.Rect) Drawer {
	s := Svg{r:r.Add(f.r.min, svg:f.svg}
	return &s
}
func (f *Svg) Bounds() image.Rectangle { return f.r }
func (f *Svg) Clear(c color.Color) { f.svg.bg = c }
func (f *Svg) Color(c color.Color) { f.svg.fg = c }
func (f *Svg) rgb() string { r, g, b, _ := f.fg.RGBA(); return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8) }
func (f *Svg) fillStroke(lw int, fill bool) string {
	s := fmt.Sprintf("stroke-width='%d' strok='%s'", lw, f.rgb())
	if fill {
		s = fmt.Sprintf("fill='%s'", f.rgb())
	}
	return s
}
func (f *Svg) Line(l Line) { f.g = append(f.g, fmt.Sprintf("<line x1='%d' y1='%d' x2='%d' y2='%d' stroke-width='%d' stroke='%s'/>", l.X,l.Y,l.X+l.DX,l.Y+l.DY, l.LineWidth, f.rgb())) }
func (f *Svg) Circle(c Circle) { f.g = append(f.g, fmt.Sprintf("<circle cx='%d' cy='%d' r='%.1f' %s>", c.X, c.Y, float64(c.D)/2, f.fillStroke(c.LineWidth, c.Fill)) }
func (f *Svg) Rectangle(r Rectangle) { f.g = append(f.g, fmt.Sprintf("<rectangle x='%d' y='%d' width='%d' height='%d' %s>", r.X, r.Y, r.W, r.H, f.fillStroke(r.LineWidth, r.Fill))) }
func (f *Svg) Triangle(t Triangle) {/*todo*/}
func (f *Svg) Ray(r Ray) {/*todo*/}
func (f *Svg) Text(t Text) {/*todo*/}
func (f *Svg) Font(f bool) {/*todo*/}
func (f *Svg) ArrowHead(a ArrowHead) {/*todo*/}

func (f *Svg) FloatTics(t FloatTics) {/*todo*/}
func (f *Svg) FloatText(t FloatText) {/*todo*/}
func (f *Svg) FloatTextExtent(t FloatText) (int,int,int,int) {/*todo*/}
func (f *Svg) FloatBars(b FloatBars) {/*todo*/}
func (f *Svg) FloatCircles(t FloatCircles) {/*todo*/}
func (f *Svg) FloatEnvelope(t FloatEnvelope) {/*todo*/}
func (f *Svg) FloatPath(t FloatPath) {/*todo*/}

func (f *Svg) Paint() {
	var b strings.Builder
	id=fmt.Sprintf("r-%d-%d-%d-%d", f.r.Min.X, f.r.Min.Y, f.r.Max.X, f.r.Max.Y)
	fmt.Fprintf(&b, "<g transform='translate(%d,%d)'>\n", f.Min.X, f.Min.Y)
	fmt.Fprintf(&b, "<defs><clipPath id='%s'><rect x='%d' y='%d' width='%d' height='%d'/></clipPath></defs\n", id, f.r.Min.X,f.r.Min.Y,f.r.Dx(),f.r.Dy())
	fmt.Fprintf(&b, "<g clip-path='url(#%s)'>\n",id)
	for _, s := range f.g {
		fmt.Fprintln(&b, s)
	}
	fmt.Fprintf(&b, "</g></g>\n")
	f.g = []string{ b.String() }
}
func (f *Svg) File() []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "<svg viewBox='0 0 %d %d' xmlns='http://www.w3.org/2000/svg'\n>", f.r.Dx(), f.r.Dy())
	for _, s := range f.g {
		b.Write([]byte(s))
	}
	fmt.Fprintf(&b, "</svg>\n")
}

