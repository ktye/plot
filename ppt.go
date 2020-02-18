package plot

import (
	"image"

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
func (p PptPlot) Decode(r pptxt.LineReader) (pptxt.Raster, error) {
	plts, err := DecodePlotsInline(r)
	p.Plots = plts
	p.Point = image.Point{600, 300} // todo
	// p.Columns
	return p, err
}

// Encode is already satisfied as both pptx and plot use the same interface definition.
