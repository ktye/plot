// Package plot plots data into images.
package plot

import "github.com/ktye/plot/color"

// Plotters is the plural of a Plotter.
type Plotters interface {
	Plots() (Plots, error)
}

// Plotter can make a Plot.
type Plotter interface {
	Plot() (Plot, error)
}

// Plots combines multiple Plots that are typically shown next to each other.
type Plots []Plot

type PlotType string

const (
	XY     PlotType = "xy"
	Foto   PlotType = "foto"
	Polar  PlotType = "polar"
	Ring   PlotType = "ring"
	AmpAng PlotType = "ampang"
	Raster PlotType = "raster"
)

// Plot is the plot definition structure.
type Plot struct {
	Type  PlotType
	Style Style
	Limits
	Xlabel, Ylabel, Title string
	Xunit, Yunit, Zunit   string
	Lines                 []Line      // data lines
	Foto                  string      // encoded Foto starting with "data:image/png;base64," or "data:image/jpeg;base64,"
	Caption               *Caption    // linked caption for this plot
	Data                  interface{} // custom userdata (unused)
}

// Style definitions for the plot. All values have defaults if unset.
type Style struct {
	Dark        bool          // Dark: green on black, otherwise black on white
	Transparent bool          // Transparent background.
	Map         color.Palette // Palette for spectrogram image.
	Order       color.Order   // Color order for lines and markers.
}

// WriteFile writes Plots to disk.
func (p *Plots) WriteFile(filename string) error {
	// TODO determine file type by extension.
	// TODO do it.
	return nil
}

// Line is a collection of data from a Plot.
type Line struct {
	X                  []float64    // x axis vector
	Y                  []float64    // y axis vector (for XY plots)
	C                  []complex128 // complex amplitude vector
	V                  []float64    // x values for vertical lines
	Image              [][]uint8    // image data: .Image[i][k] corresponds to axes values .X[i], .Y[i].
	ImageMin, ImageMax float64      // image color scale corresponding to [0..255]
	Segments           bool         // Lines contain segments
	Style              DataStyle    // line style
	Id                 int          // line id used for selection
}

// DataStyle is the line and marker style definition, which is set for
// every line object individually.
type DataStyle struct {
	Line   LineStyle
	Marker MarkerStyle
}

// FillLineStyle returns a DataStyle with default values for a line
// and ignored values for a marker.
func (d DataStyle) FillLineStyle() DataStyle {
	s := d
	s.Marker.Size = 0
	if s.Line.Width == 0 {
		s.Line.Width = 1
	}
	return s
}

// FillMarkerStyle returns a DataStyle with default values for a marker
// and ignored values for a line.
func (d DataStyle) FillMarkerStyle() DataStyle {
	s := d
	s.Line.Width = 0
	if s.Marker.Size == 0 {
		s.Marker.Size = 3
	}
	if s.Marker.Marker == NoMarker {
		s.Marker.Marker = PointMarker
	}
	return s
}

// Combined style includes the plot style and the data style.
// It is used by plot objects which do not style for every line manually.
type CombinedStyle struct {
	Plot Style
	Data DataStyle
}

// RemoveEmpty deletes empty plots from the plot slice.
func (p *Plots) RemoveEmpty() Plots {
	var y Plots
	for _, v := range *p {
		if v.Foto != "" || len(v.Lines) > 0 {
			y = append(y, v)
		}
	}
	return y
}

// NextNegativeLineId checks if the last line has a negative Id, and returns it's decrement.
// Otherwise it returns -1.
// It is used for interactively adding lines (or points) and giving them unique IDs.
func (p *Plot) nextNegativeLineId() int {
	if n := len(p.Lines); n > 0 {
		if id := p.Lines[n-1].Id; id < 0 {
			return id - 1
		}
	}
	return -1
}

// MergedCaption returns a single caption from plots.
func (p *Plots) MergedCaption() (c Caption, err error) {
	var caps []Caption
	for _, v := range *p {
		if v.Caption != nil {
			caps = append(caps, *v.Caption)
		}
	}
	return MergeCaptions(caps)
}
