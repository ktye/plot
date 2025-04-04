package plot

import (
	"image/color"

	vgfont "github.com/ktye/plot/vg/font"
	"golang.org/x/image/font"
)

var font1, font2 font.Face

func init() {
	font1, font2 = vgfont.MakeFontSizes(vgfont.TTF(), 16, 12)
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
	return 5 //font1.Metrics().Height.Ceil() / 4
}

func (p *Plot) defaultForegroundColor() color.Color {
	if p.Style.Dark {
		return color.RGBA{200, 200, 200, 255}
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
	if p.Title == "" {
		return 2
	}
	return 20 //2 + font1.Metrics().Height.Ceil()
}

// defaultXLabelHeight is the vertical space for the x-axis label.
func (p *Plot) defaultXlabelHeight() int {
	if p.Xlabel+p.Xunit == "" {
		return 2
	}
	return 20 //2 + font1.Metrics().Height.Ceil()
}

// defaultYLabelWidth is the horizontal space for the rotated y-axis label.
func (p *Plot) defaultYlabelWidth() int {
	if p.Ylabel+p.Yunit == "" {
		return 2
	}
	return 16 //2 + font1.Metrics().Height.Ceil()
}

// defaultTicLabelHeight is the vertical space for the tic labels.
func (p *Plot) defaultTicLabelHeight() int {
	return 22 //font2.Metrics().Height.Ceil() * 2
}

// defaultPolarTicLabelWidth is the horizontal space for polar tic labels.
// which is basically the size around the string "270".
func (p *Plot) defaultPolarTicLabelWidth() int {
	return 42 //3 * font2.Metrics().Height.Ceil()
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
	return width
}

// defaultRightXYWidth is the horizontal space on the right side of an xy axis.
// This is needed to draw half the last x-axis tic.
func (p *Plot) defaultRightXYWidth(lastXLabel string) int {
	width := stringWidth(font2, lastXLabel)
	return (width+1)/2 + 7 // Half the label + some space.
}

// defaultTicLength computes the tic length.
func (p *Plot) defaultTicLength() int {
	return 6 //font2.Metrics().Height.Ceil() / 2
}

// defaultAxesGridLineWidth returns the grid line width.
func (p *Plot) defaultAxesGridLineWidth() int {
	return 1
}

// defaultAmpAngSpace is the vertical space between the amplitude
// and the phase axes for an ampang plot.
func (p *Plot) defaultAmpAngSpace() int {
	return 2 * p.defaultTicLength()
}
