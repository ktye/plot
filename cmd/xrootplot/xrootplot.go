// revival of an old work horse
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/ktye/plot"
	"github.com/ktye/plot/vg"
)

type i = int32
type k = uint32
type c = byte
type f = float64

const (
	TERM    = 0
	CONSOLE = 1
	FILE    = 2
)
const use = `|plot [-dtcw:h:o:] [-plt] [0 2 ..]
ktye/plot/cmd/xrootplot/xrootplot.go
 -d   print (uniform) plot data [csv]
 -t   print table(caption) data [csv]
 -1   single axes
 FILE.png (to file as png instead of stdout)
 -c   console output (default iterm2 image)
 -p   convert to plt format
 -wWIDTH -hHEIGHT (also from env)
  0.. plot number`

var dst, dat, tab, sgl, plt, wid, hei, out = TERM, false, false, false, false, 0, 0, ""
var idx []int

func main() {
	wid, hei = atoi0(os.Getenv("WIDTH")), atoi0(os.Getenv("HEIGHT"))
	args := os.Args[1:]
	for _, s := range args {
		if n, e := strconv.Atoi(s); e == nil && n >= 0 {
			idx = append(idx, n)
		} else if pre(s, "-d") {
			dat = true
		} else if pre(s, "-t") {
			dat, tab = true, true
		} else if s == "-1" {
			sgl = true
		} else if pre(s, "-c") {
			dst = CONSOLE
		} else if pre(s, "-w") {
			wid = atoi(s[2:])
		} else if pre(s, "-h") {
			hei = atoi(s[2:])
		} else if suf(s, ".png") {
			dst = FILE
			out = s[1:]
		} else if pre(s, "-p") {
			plt = true
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
func atoi(s string) int {
	n, e := strconv.Atoi(s)
	fatal(e)
	return n
}
func atoi0(s string) int {
	if n, e := strconv.Atoi(s); e != nil {
		return 0
	} else {
		return n
	}
}
func do(r io.Reader) {
	plts, e := plot.DecodeAny(r)
	fatal(e)
	plts = dark(single(at(plts)))
	pp(plot.AxisFromEnv(plts))
}
func pp(p plot.Plots) {
	if plt {
		fatal(p.Encode(os.Stdout))
		return
	}
	w, h := screensize()
	d := vg.NewImage(w, h)
	ip, e := p.Iplots(d, 0)
	fatal(e)
	m := ip.Image(nil).(*image.RGBA)
	switch dst {
	case TERM:
		draw(pngData(m))
	case CONSOLE:
		drawConsole(w, h, m.Pix)
	case FILE:
		fatal(ioutil.WriteFile(out, pngData(m), 0644))
	default:
		fatal(fmt.Errorf("unknown dst: %d", dst))
	}
}
func screensize() (w, h int) {
	if dst == CONSOLE {
		w, h = consoleSize()
	}
	if wid != 0 {
		w = wid
	}
	if hei != 0 {
		h = hei
	}
	if w == 0 {
		w = 800
	}
	if h == 0 {
		h = 600
	}
	return w, h
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
func dark(p plot.Plots) plot.Plots { // console is always dark, file or plt not so.
	if dst != FILE { // && plt == false {
		for i := range p {
			p[i].Style.Dark = true
			p[i].Style.Order = "green"
		}
	}
	return p
}
func single(p plot.Plots) plot.Plots {
	if len(p) < 2 || sgl == false {
		return p
	}
	for _, x := range p[1:] {
		p[0].Lines = append(p[0].Lines, x.Lines...)
	}
	for i := range p[0].Lines {
		p[0].Lines[i].Id = i
	}
	return p[:1]
}
func pre(s, p string) bool { return strings.HasPrefix(s, p) }
func suf(s, x string) bool { return strings.HasSuffix(s, x) }
func data(r io.Reader) {
	plts, e := plot.DecodeAny(r)
	fatal(e)
	plts = at(plts)
	if len(plts) == 0 {
		return
	} else if tab {
		if c := plts[0].Caption; c != nil {
			c.WriteTable(os.Stdout, 0)
		}
		return
	}
	fatal(plts.WriteCsv(os.Stdout, false))
}
func pngData(m image.Image) []byte {
	var buf bytes.Buffer
	fatal(png.Encode(&buf, m))
	return buf.Bytes()
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
