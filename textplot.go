package plot

import (
	"image"
	"image/color"
	"strings"

	"github.com/ktye/plot/vg"
	"golang.org/x/image/draw"
)

// textPlot shows text as a plot image.
type textPlot struct {
	plot *Plot
	im   *image.RGBA
}

// ErrorPlot returns Plots that render an error message.
func ErrorPlot(err error) Plots {
	s := ""
	if err != nil {
		s = err.Error()
	}
	return Plots{{Type: Text, Foto: s}}
}

// Create a new text plot in the subimage.
func (plt *Plot) NewTextPlot(width, height int) (p textPlot, err error) {
	p.plot = plt

	p.im = image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(p.im, p.im.Bounds(), image.NewUniform(p.plot.defaultBackgroundColor()), image.ZP, draw.Src)

	lines := strings.Split(plt.Foto, "\n")
	paint := vg.NewPainter(p.im)
	//paint.SetFont(font1)
	paint.SetColor(p.plot.defaultForegroundColor())
	x, y := 10, 10
	dy := font1.Metrics().Height.Ceil()
	for _, s := range lines {
		paint.Add(vg.Text{X: x, Y: y, S: s, Align: 6})
		y += dy
	}
	paint.Paint()

	return p, nil
}

func (p textPlot) background() color.Color {
	return p.plot.defaultBackgroundColor()
}

func (p textPlot) image() *image.RGBA {
	return p.im
}

func (p textPlot) zoom(x, y, dx, dy int) bool {
	return false
}

func (p textPlot) pan(x, y, dx, dy int) bool {
	return false
}

func (p textPlot) limits() Limits {
	return Limits{}
}

func (p textPlot) line(x0, y0, x1, y1 int) (complex128, bool) {
	return complex(0, 0), false
}

func (p textPlot) click(x, y int, snapToPoint bool) (Callback, bool) {
	return Callback{}, false
}

func (p textPlot) highlight(id []HighlightID) *image.RGBA {
	return p.im
}
