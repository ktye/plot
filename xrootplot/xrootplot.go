// revival of an old work horse
package main

import (
	"fmt"
	"image"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ktye/plot"
)

type i = int32
type k = uint32
type c = byte
type f = float64

const (
	CONSOLE = 0
	SCREEN  = 1
)
const use = `|plot [-dts] [-csv] [0 2 ..]
 -d   print (uniform) plot data [csv]
 -t   print table(caption) data [csv]
 -s   plot over screen (def. console)
  0.. plot number
`

var dst, dat, tab, csv = CONSOLE, false, false, false
var idx []int

func main() {
	args := os.Args[1:]
	for _, s := range args {
		if n, e := strconv.Atoi(s); e == nil && n >= 0 {
			idx = append(idx, n)
		} else if pre(s, "-d") {
			dat = true
		} else if pre(s, "-t") {
			dat, tab = true, true
		} else if pre(s, "-c") {
			dat, csv = true, true
		} else {
			fmt.Println(use)
			return
		}
	}
	if dat {
		do(os.Stdin)
	} else {
		data(os.Stdin)
	}
}
func do(r io.Reader) {
	plts, e := plot.DecodeAny(r)
	fatal(e)
	plts = at(plts)
	pp(plot.AxisFromEnv(plts))
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
func at(p plot.Plots) plot.Plots {
	if idx == nil {
		return p
	} else {
		r := make(plot.Plots, len(idx))
		for i, n := range idx {
			if n >= 0 && n < len(p) {
				r[i] = p[n]
			} else {
				fatal(fmt.Errorf("index out of range %d [0,%d]", n, len(p)-1))
			}
		}
		return r
	}
}
func pre(s, p string) bool { return strings.HasPrefix(s, p) }
func suf(s, x string) bool { return strings.HasSuffix(s, x) }
func fatal(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
func init() {
	plot.SetFonts(Face10x20, Face10x20)
}
