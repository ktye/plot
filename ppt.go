package plot

import (
	"fmt"
	"image"
	"io"

	"github.com/ktye/pptx/pptxt"
)

// PptPlot defines a type that can use the text representation of a plot within ktye/pptx.
// It implements pptx.Raster without the need to import the pptx package.
// The program ktye/pptx/pptadd supports it.
type PptPlot struct {
	Plots           // plot data, resolution independend
	image.Point     // rasterization size (not slide position)
	Columns     int // number of columns (for multi-row plot)
	Highlight   []HighlightID
}

func (p PptPlot) Raster() (image.Image, error) {
	ip, err := p.Plots.IPlots(p.X, p.Y, p.Columns)
	if err != nil {
		return nil, err
	}
	return Image(ip, p.Highlight, p.X, p.Y, p.Columns), nil
}
func (p PptPlot) Magic() string { return "Plot" }
func (p PptPlot) Encode(w io.Writer) error {
	_, e := fmt.Fprintf(w, "PlotSize [%d, %d]\nPlotColumns %d\n", p.Point.X, p.Point.Y, p.Columns)
	if e != nil {
		return e
	}
	return p.Plots.Encode(w)
}
func (p PptPlot) Decode(r pptxt.LineReader) (ra pptxt.Raster, e error) {
	var xy [2]int
	e = sj(r, "PlotSize", &xy, e)
	p.Point.X, p.Point.Y = xy[0], xy[1]
	e = sj(r, "PlotColumns", &p.Columns, e)
	if e != nil {
		return nil, e
	}
	plts, err := DecodePlotsInline(r)
	p.Plots = plts
	return p, err
}
