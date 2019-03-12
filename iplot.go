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
func (p Plots) IPlots(width, height int) ([]IPlotter, error) {
	var err error
	if len(p) == 0 {
		return nil, fmt.Errorf("there are no plots.")
	}

	plotters := make([]IPlotter, len(p))
	for i := 0; i < len(p); i++ {
		w := width / len(p)
		switch p[i].Type {
		case "":
			plotters[i], err = p[i].NewEmpty(w, height)
		case XY, Raster:
			plotters[i], err = p[i].NewXY(w, height)
		case Polar:
			plotters[i], err = p[i].NewPolar(w, height)
		case AmpAng:
			plotters[i], err = p[i].NewAmpAng(w, height)
		case Foto:
			plotters[i], err = p[i].NewFoto(w, height)
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
func Image(h []IPlotter, ids []HighlightID, width, height int) image.Image {
	if len(h) < 1 {
		return nil
	}
	m := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(m, m.Bounds(), &image.Uniform{h[0].background()}, image.ZP, draw.Src)

	w := width / len(h)
	for i := range h {
		im := h[i].image()
		im = h[i].highlight(ids)

		// center the sub image in it's columns.
		bounds := im.Bounds()
		imwidth := bounds.Max.X - bounds.Min.X
		imheight := bounds.Max.Y - bounds.Min.Y
		yoff := (height - imheight) / 2
		xoff := i*w + (w-imwidth)/2
		rect := image.Rect(xoff, yoff, imwidth+xoff, imheight+yoff)
		draw.Draw(m, rect, im, image.Point{0, 0}, draw.Src)
	}
	return m
}

type imageResult struct {
	im    *image.RGBA
	index int
}
