package plot

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

// IPlotter is an interactive plotter, that is a plot
// which can highlight a specific data series and
// has zoom and pan methods.
type IPlotter interface {
	image() *image.RGBA
	highlight([]HighlightID) *image.RGBA
	background() color.Color
	zoom(int, int, int, int) bool
	pan(int, int, int, int) bool
	click(int, int, bool) (Callback, bool)
	line(int, int, int, int) (complex128, bool)
	limits() Limits
}

// IPlots returns a slice of IPlotters, one for each plot.
// The subplots are shown next to each other.
func (p Plots) IPlots(width, height, columns int) ([]IPlotter, error) {
	var err error
	if len(p) == 0 {
		return nil, fmt.Errorf("there are no plots.")
	}

	g := newGrid(len(p), width, height, columns)
	w, h := g.w, g.h

	plotters := make([]IPlotter, len(p))
	for i := 0; i < len(p); i++ {
		switch p[i].Type {
		case "":
			plotters[i], err = p[i].NewEmpty(w, h)
		case XY, Raster:
			plotters[i], err = p[i].NewXY(w, h)
		case Polar:
			plotters[i], err = p[i].NewPolar(w, h, false)
		case Ring:
			plotters[i], err = p[i].NewPolar(w, h, true)
		case AmpAng:
			plotters[i], err = p[i].NewAmpAng(w, h)
		case Foto:
			plotters[i], err = p[i].NewFoto(w, h)
		case Text:
			plotters[i], err = p[i].NewTextPlot(w, h)
		default:
			return plotters, fmt.Errorf("plot type: '%s' is not implemented.", p[i].Type)
		}
		if err != nil {
			return nil, err
		}
	}
	return plotters, nil
}

// Image creates an Image from a slice of iplotters.
// If ids is nil, no lines will be highlighted.
// The image dimensions may be differnt than the in the initial call to IPlots,
// e.g. after a resize.
func Image(h []IPlotter, ids []HighlightID, width, height, columns int) image.Image {
	if len(h) < 1 {
		return nil
	}
	m := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(m, m.Bounds(), &image.Uniform{h[0].background()}, image.ZP, draw.Src)
	g := newGrid(len(h), width, height, columns)
	// w := width / len(h)
	for i := range h {
		im := h[i].image()
		im = h[i].highlight(ids)
		rect := g.center(g.rect(i), im)
		draw.Draw(m, rect, im, image.Point{0, 0}, draw.Src)
	}
	return m
}

type imageResult struct {
	im    *image.RGBA
	index int
}

type grid struct {
	plots, width, height, maxcols int
	rows, cols, w, h              int
}

func newGrid(plots, width, height, maxcols int) grid {
	g := grid{plots: plots, width: width, height: height, maxcols: maxcols}
	if g.maxcols == 0 {
		g.maxcols = 4
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
