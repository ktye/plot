package plot

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"

	"github.com/ktye/plot/vg"
)

// IPlotter is an interactive plotter, that is a plot
// which can highlight a specific data series and
// has zoom and pan methods.
type Iplotter interface {
	image() *image.RGBA
	highlight([]HighlightID) *image.RGBA
	background() color.Color
	zoom(int, int, int, int) bool
	pan(int, int, int, int) bool
	click(int, int, bool, bool, bool) (Callback, bool)
	line(int, int, int, int) (complex128, bool)
	measure(int, int, int, int) (MeasureInfo, bool)
	limits() Limits
}

type Iplots struct {
	p []Iplotter
	d vg.Drawer
	g grid
}

func (ip Iplots) SetLimitsTo(p *Plots) {
	if len(ip.p) == len(*p) {
		for i := range ip.p {
			(*p)[i].Limits = ip.p[i].limits()
			//todo: force these limits for polar plots (nonzero origin)
		}
	}
}

// IPlots returns a slice of IPlotters, one for each plot.
// The subplots are shown next to each other.
func (p Plots) Iplots(d vg.Drawer, columns int) (Iplots, error) {
	width, height := d.Size()
	var err error
	if len(p) == 0 {
		return Iplots{}, fmt.Errorf("there are no plots.")
	}
	g := newGrid(len(p), width, height, columns)
	r := Iplots{
		p: make([]Iplotter, len(p)),
		d: d,
		g: g,
	}
	if len(p) > 0 {
		d.Clear(p[0].defaultBackgroundColor())
	}
	for i := 0; i < len(p); i++ {
		rect := g.rect(i)
		switch p[i].Type {
		case "":
			r.p[i], err = p[i].NewEmpty(d.SubImage(rect))
		case XY, Raster:
			r.p[i], err = p[i].NewXY(d.SubImage(rect))
		case Polar:
			r.p[i], err = p[i].NewPolar(d.SubImage(rect), false)
		case Ring:
			r.p[i], err = p[i].NewPolar(d.SubImage(rect), true)
		case AmpAng:
			r.p[i], err = p[i].NewAmpAng(d.SubImage(rect))
		case Waterfall:
			r.p[i], err = p[i].NewWaterfall(d.SubImage(rect))
		case Foto:
			r.p[i], err = p[i].NewFoto(d.SubImage(rect))
		case Text:
			r.p[i], err = p[i].NewTextPlot(d.SubImage(rect))
		default:
			return r, fmt.Errorf("plot type: '%s' is not implemented.", p[i].Type)
		}
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

func (p Iplots) Image(ids []HighlightID) image.Image {
	if len(p.p) < 1 {
		return nil
	}
	for i := range p.p {
		p.p[i].highlight(ids)
	}

	im, ok := p.d.(*vg.Image)
	if !ok {
		return nil
	}
	return im.RGBA
}

func (p Plots) Png(width, height, columns int, idx []HighlightID) ([]byte, error) {
	ip, e := p.Iplots(vg.NewImage(width, height, font1, font2), columns)
	if e != nil {
		return nil, e
	}
	m := ip.Image(idx)
	if m == nil {
		return nil, fmt.Errorf("no image")
	}
	var b bytes.Buffer
	e = png.Encode(&b, m)
	return b.Bytes(), e
}
func (p Plots) Wmf(width, height, columns int, idx []HighlightID) ([]byte, error) {
	ip, e := p.Iplots(vg.NewWmf(width, height), columns)
	if e != nil {
		return nil, e
	}
	//todo: for i := range p.p { p.p[i].highlight(ids) }
	return ip.d.(*vg.Wmf).MarshallBinary(), nil
}
func (p Plots) Emf(width, height, columns int, idx []HighlightID, font string, f1, f2 int) ([]byte, error) { //font f1 f2 maybe empty
	ip, e := p.Iplots(vg.NewEmf(width, height, font, f1, f2), columns)
	if e != nil {
		return nil, e
	}
	for i := range ip.p {
		ip.p[i].highlight(idx)
	}
	return ip.d.(*vg.Emf).MarshallBinary(), nil
}
func (p Plots) Svg(width, height, columns int, idx []HighlightID) ([]byte, error) {
	ip, e := p.Iplots(vg.NewSvg(width, height), columns)
	if e != nil {
		return nil, e
	}
	//todo: for i := range p.p { p.p[i].highlight(ids) }
	return ip.d.(*vg.Svg).Bytes(), nil
}

func (p Iplots) IsPolar(n int) bool {
	var ok bool
	if n >= 0 && n < len(p.p) && p.p != nil {
		_, ok = p.p[n].(polarPlot)
	}
	return ok
}

//// Image creates an Image from a slice of iplotters.
//// If ids is nil, no lines will be highlighted.
//// The image dimensions may be differnt than the in the initial call to IPlots,
//// e.g. after a resize.
//func Image(h []IPlotter, ids []HighlightID, width, height, columns int) image.Image {
//	if len(h) < 1 {
//		return nil
//	}
//	/*
//		m := image.NewRGBA(image.Rect(0, 0, width, height))
//		draw.Draw(m, m.Bounds(), &image.Uniform{h[0].background()}, image.ZP, draw.Src)
//		g := newGrid(len(h), width, height, columns)
//		// w := width / len(h)
//		for i := range h {
//			im := h[i].image()
//			im = h[i].highlight(ids)
//			rect := g.center(g.rect(i), im)
//			draw.Draw(m, rect, im, image.Point{0, 0}, draw.Src)
//		}
//		return m
//	*/
//	return h[0].image() //todo..
//}

type imageResult struct {
	im    *image.RGBA
	index int
}

type grid struct {
	plots, width, height, maxcols int
	rows, cols, w, h              int
	colmajor                      bool
}

func RowsCols(nplots, maxcols int) (rows, cols int) { //used externally
	g := newGrid(nplots, 100, 100, maxcols)
	return g.rows, g.cols
}
func newGrid(plots, width, height, maxcols int) grid {
	g := grid{plots: plots, width: width, height: height, maxcols: maxcols}
	if g.maxcols == 0 {
		//         0  1  2  3  4  5  6  7  8  9 10 11 12
		q := []int{4, 4, 4, 4, 4, 3, 3, 4, 4, 5, 5, 4, 4}
		if plots >= 0 && plots < len(q) {
			g.maxcols = q[plots]
		} else {
			g.maxcols = 5
		}
	}
	if g.maxcols < 0 {
		g.maxcols = -g.maxcols
		g.colmajor = true
	}
	if plots <= g.maxcols {
		g.rows, g.cols = 1, plots
	} else {
		g.rows = plots / g.maxcols
		g.cols = g.maxcols
		if g.rows*g.cols < plots {
			g.rows++
		}
	}
	g.w = width / g.cols
	g.h = height / g.rows
	return g
}
func (g grid) rect(n int) image.Rectangle {
	i, k := n/g.cols, n%g.cols
	if g.colmajor {
		i, k = k, i
	}
	x := 0
	if i == (g.plots-1)/g.cols { // center plots on last row (space on the left and right)
		plots := 1 + ((g.plots - 1) % g.cols)
		x = (g.width - plots*g.w) / 2
	}
	x += k * g.w
	y := i * g.h
	return image.Rect(x, y, x+g.w, y+g.h)
}

func (g grid) center(dst image.Rectangle, im image.Image) image.Rectangle {
	bounds := im.Bounds()
	pt := image.Point{}
	if dw := dst.Dx() - bounds.Dx(); dw > 0 {
		pt.X = dw / 2
	}
	if dh := dst.Dy() - bounds.Dy(); dh > 0 {
		pt.Y = dh / 2
	}
	return dst.Add(pt)
}
