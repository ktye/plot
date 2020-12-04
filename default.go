package plot

import (
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

var font1 font.Face = basicfont.Face7x13
var font2 font.Face = basicfont.Face7x13

// SetFonts sets two fonts used plots.
// A larger for title and axis labels, and usually a smaller one for tic labels.
func SetFonts(labels, tics font.Face) {
	font1, font2 = labels, tics
}

func Fonts() (font.Face, font.Face) {
	return font1, font2
}

// StringWidth return the space needed to draw a string in the font f.
func stringWidth(f font.Face, s string) int {
	d := font.Drawer{Face: f}
	return d.MeasureString(s).Ceil()
}

// defaultBorder is the outer border around a plot. It is cut from the requested size.
func (p *Plot) defaultBorder() int {
	return font1.Metrics().Height.Ceil() / 4
}

func (p *Plot) defaultForegroundColor() color.Color {
	if p.Style.Dark {
		return color.RGBA{32, 230, 32, 255}
	}
	return color.Black
}

func (p *Plot) defaultBackgroundColor() color.Color {
	if p.Style.Transparent {
		return color.Transparent
	}
	if p.Style.Dark {
		return color.Black
	}
	return color.White
}

// defaultTitleHeight is the vertical space for the plot title.
func (p *Plot) defaultTitleHeight() int {
	return font1.Metrics().Height.Ceil() * 2
}

// defaultTicLabelHeight is the vertical space for the tic labels.
func (p *Plot) defaultTicLabelHeight() int {
	return font2.Metrics().Height.Ceil() * 2
}

// defaultPolarTicLabelWidth is the horizontal space for polar tic labels.
// which is basically the size around the string "270".
func (p *Plot) defaultPolarTicLabelWidth() int {
	return 3 * font2.Metrics().Height.Ceil()
}

// defaultTicLabelWidth is the horizontal space for tic labels on the y axis.
func (p *Plot) defaultTicLabelWidth(ylabels []string) int {
	// Calculate the space needed for the longest tic label.
	width := 0
	for _, l := range ylabels {
		if w := stringWidth(font2, l); w > width {
			width = w
		}
	}
	return width + 3 // Add some space to the longest string.
}

// defaultRightXYWidth is the horizontal space on the right side of an xy axis.
// This is needed to draw half the last x-axis tic.
func (p *Plot) defaultRightXYWidth(lastXLabel string) int {
	width := stringWidth(font2, lastXLabel)
	return (width+1)/2 + 7 // Half the label + some space.
}

// defaultTicLength computes the tic length.
func (p *Plot) defaultTicLength() int {
	return font2.Metrics().Height.Ceil() / 2
}

// defaultAxesGridLineWidth returns the grid line width.
func (p *Plot) defaultAxesGridLineWidth() int {
	n := font2.Metrics().Height.Ceil() / 10
	if n == 0 {
		n = 1
	}
	return n
}

// defaultAmpAngSpace is the vertical space between the amplitude
// and the phase axes for an ampang plot.
func (p *Plot) defaultAmpAngSpace() int {
	return font1.Metrics().Height.Ceil() / 2
}
