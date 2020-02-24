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
const use = `|plot [-dts] [-csv] [-plt] [-a(p|l)(sleep)] [0 2 ..]
 -d   print (uniform) plot data [csv]
 -t   print table(caption) data [csv]
 -s   plot over screen (def. console)
 -al  animate(over lines) sleep 100ms
 -ap  animate(over plots/columns)
 -p   convert to plt format
  0.. plot number
`

var dst, dat, tab, csv, plt, ani, fps = CONSOLE, false, false, false, false, 0, 100
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
		} else if pre(s, "-p") {
			plt = true
		} else if pre(s, "-al") {
			ani, fps = 1, parseRate(s[3:])
		} else if pre(s, "-ap") {
			ani, fps = 2, parseRate(s[3:])
		} else {
			fmt.Println(use)
			return
		}
	}
	if dat {
		data(os.Stdin)
	} else {
		do(os.Stdin)
	}
}
func do(r io.Reader) {
	plts, e := plot.DecodeAny(r)
	fatal(e)
	plts = at(plts)
	pp(plot.AxisFromEnv(plts))
}
func pp(p plot.Plots) {
	if ani > 0 {
		animate(p)
	} else if plt {
		fatal(p.Encode(os.Stdout))
		return
	}
	w, h := screensize()
	ip, e := p.IPlots(w, h, 0)
	fatal(e)
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
func parseRate(s string) int {
	if n, e := strconv.Atoi(s); e != nil || n < 0 {
		return 100
	} else {
		return n
	}
}
func fatal(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
func init() {
	plot.SetFonts(Face10x20, Face10x20)
}
