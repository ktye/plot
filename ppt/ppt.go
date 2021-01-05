package ppt

import (
	"fmt"
	"image"
	"io"

	"github.com/ktye/plot"
	"github.com/ktye/pptx/pptxt"
)

// Plot defines a type that can use the text representation of a plot within ktye/pptx.
// It implements pptx.Raster without the need to import the pptx package.
// The program ktye/pptx/pptadd supports it.
type Plot struct {
	plot.Plots      // plot data, resolution independend
	image.Point     // rasterization size (not slide position)
	Columns     int // number of columns (for multi-row plot)
	Highlight   []plot.HighlightID
}

func (p Plot) Raster() (image.Image, error) {
	ip, err := p.Plots.IPlots(p.X, p.Y, p.Columns)
	if err != nil {
		return nil, err
	}
	return plot.Image(ip, p.Highlight, p.X, p.Y, p.Columns), nil
}
func (p Plot) Magic() string { return "Plot" }
func (p Plot) Encode(w io.Writer) error {
	_, e := fmt.Fprintf(w, "PlotSize [%d, %d]\nPlotColumns %d\n", p.Point.X, p.Point.Y, p.Columns)
	if e != nil {
		return e
	}
	return p.Plots.Encode(w)
}
func (p Plot) Decode(r pptxt.LineReader) (ra pptxt.Raster, e error) {
	var xy [2]int
	e = plot.Sj(r, "PlotSize", &xy, e)
	p.Point.X, p.Point.Y = xy[0], xy[1]
	e = plot.Sj(r, "PlotColumns", &p.Columns, e)
	if e != nil {
		return nil, e
	}
	plts, err := plot.DecodePlotsInline(r)
	p.Plots = plts
	return p, err
}
