package plot

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/ktye/plot/vg"
)

// emptyPlot implements an IPlotter that returns an empty image.
// It can be usefull when arranging plots in a grid with missing items.
type emptyPlot struct {
	plot          *Plot
	width, height int
	color         color.Color
	im            *image.RGBA
}

func (plt *Plot) NewEmpty(d vg.Drawer) (p *emptyPlot, err error) {
	w, h := d.Size()
	return &emptyPlot{plot: plt, width: w, height: h}, nil
}

func (e *emptyPlot) image() *image.RGBA {
	if e.im == nil {
		r := image.Rectangle{Max: image.Point{e.width, e.height}}
		e.im = image.NewRGBA(r)
		if e.color == nil {
			if e.plot.Style.Transparent {
				e.color = color.Transparent
			} else if e.plot.Style.Dark {
				e.color = color.Black
			} else {
				e.color = color.White
			}
		}
		draw.Draw(e.im, r, &image.Uniform{e.color}, image.ZP, draw.Src)
	}
	return e.im
}

func (e *emptyPlot) highlight([]HighlightID) *image.RGBA {
	return e.image()
}
func (e *emptyPlot) background() color.Color {
	e.image()
	return e.color
}
func (e *emptyPlot) zoom(x, y, dx, dy int) bool                  { return false }
func (e *emptyPlot) pan(x, y, dx, dy int) bool                   { return false }
func (e *emptyPlot) click(int, int, bool, bool) (Callback, bool) { return Callback{}, false }
func (e *emptyPlot) line(x0, y0, x1, y1 int) (complex128, bool)  { return complex(0, 0), false }
func (e *emptyPlot) limits() Limits                              { return Limits{} }
