// revival of an old work horse
package main

import (
	"image"
	"io"
	"os"

	"github.com/ktye/plot"
)

type i = int32
type k = uint32
type c = byte
type f = float64

func main() { do(os.Stdin) }
func do(r io.Reader) {
	plts, e := plot.TextDataPlot(r)
	if e != nil {
		panic(e)
	}
	pp(plts)
}
func pp(p plot.Plots) {
	w, h := screensize()
	ip, e := p.IPlots(w, h, 0)
	if e != nil {
		panic(e)
	}
	m := plot.Image(ip, nil, w, h, 0).(*image.RGBA)
	draw(w, h, m.Pix)
}

func init() {
	plot.SetFonts(Face10x20, Face10x20)
}
