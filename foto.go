package plot

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/ktye/plot/vg"
	"golang.org/x/image/draw"
)

// fotoPlot is an implementation of the HiPlotter interface.
type fotoPlot struct {
	plot   *Plot
	im     *image.RGBA     // complete available sub image
	rect   image.Rectangle // destination rectangle for scaled foto keeping ratio
	drawer vg.Drawer
}

const PngPrefix = "data:image/png;base64,"
const JpegPrefix = "data:image/jpeg;base64,"

func EncodeToPng(m image.Image) (string, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, m); err != nil {
		return "", err
	}
	return PngPrefix + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// Create a new foto plot in the subimage.
func (plt *Plot) NewFoto(d vg.Drawer) (p fotoPlot, err error) {
	width, height := d.Size()
	p.drawer = d
	p.plot = plt
	border := plt.defaultBorder()

	var foto image.Image
	if strings.HasPrefix(plt.Foto, PngPrefix) {
		var b []byte
		b, err = base64.StdEncoding.DecodeString(plt.Foto[len(PngPrefix):])
		if err == nil {
			r := bytes.NewReader(b)
			foto, err = png.Decode(r)
		}
	} else if strings.HasPrefix(plt.Foto, JpegPrefix) {
		var b []byte
		b, err = base64.StdEncoding.DecodeString(plt.Foto[len(JpegPrefix):])
		if err == nil {
			r := bytes.NewReader(b)
			foto, err = jpeg.Decode(r)
		}
	} else {
		return p, fmt.Errorf("foto plot has wrong suffix")
	}
	if err != nil {
		return p, fmt.Errorf("cannot decode foto: %s", err)
	}

	// Available space for plotArea.
	widthAvail := width - 2*border
	heightAvail := height - 2*border

	scale := float64(widthAvail) / float64(foto.Bounds().Dx())
	if h := int(scale * float64(foto.Bounds().Dy())); h <= heightAvail {
		x0 := border
		x1 := x0 + widthAvail
		y0 := border + (heightAvail-h)/2
		y1 := y0 + h
		p.rect = image.Rect(x0, y0, x1, y1)
	} else {
		scale = float64(heightAvail) / float64(foto.Bounds().Dy())
		w := int(scale * float64(foto.Bounds().Dx()))
		x0 := border + (widthAvail-w)/2
		x1 := x0 + w
		y0 := border
		y1 := y0 + heightAvail
		p.rect = image.Rect(x0, y0, x1, y1)
	}
	p.rect = p.rect.Add(d.Bounds().Min)
	p.draw(foto)
	return p, nil
}
func (p fotoPlot) draw(foto image.Image) {
	m, o := p.drawer.(*vg.Image)
	if o == false {
		return
	}
	draw.Draw(m.RGBA, m.RGBA.Bounds(), image.NewUniform(p.plot.defaultBackgroundColor()), m.RGBA.Bounds().Min, draw.Src)
	draw.ApproxBiLinear.Scale(m.RGBA, p.rect, foto, foto.Bounds(), draw.Src, nil)
}
func (p fotoPlot) background() color.Color                        { return p.plot.defaultBackgroundColor() }
func (p fotoPlot) image() *image.RGBA                             { return p.drawer.(*vg.Image).RGBA }
func (p fotoPlot) zoom(x, y, dx, dy int) bool                     { return false }
func (p fotoPlot) pan(x, y, dx, dy int) bool                      { return false }
func (p fotoPlot) limits() Limits                                 { return Limits{} }
func (p fotoPlot) measure(x0, y0, x1, y1 int) (MeasureInfo, bool) { return MeasureInfo{}, false }
func (p fotoPlot) line(x0, y0, x1, y1 int) (complex128, bool)     { return complex(0, 0), false }
func (p fotoPlot) click(x, y int, snapToPoint, deleteLine, dodraw bool) (Callback, bool) {
	return Callback{}, false
}
func (p fotoPlot) highlight(id []HighlightID) *image.RGBA { return p.drawer.Rgba() }
