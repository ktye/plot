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
 -d   print (uniform) plot data [csv]
 -t   print table(caption) data [csv]
 -c   console output (default iterm2 image)
 -wWIDTH -hHEIGHT
 -o   file.png
 -p   convert to plt format
  0.. plot number
`

var dst, dat, tab, plt, wid, hei, out = TERM, false, false, false, 0, 0, ""
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
			dst = CONSOLE
		} else if pre(s, "-w") {
			wid = atoi(s[2:])
		} else if pre(s, "-h") {
			hei = atoi(s[2:])
		} else if pre(s, "-o") {
			dst = FILE
			out = parseFile(s[2:])
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
func parseFile(s string) string {
	if strings.HasSuffix(s, ".png") == false {
		fatal(fmt.Errorf("-oFILE.png only png is supported"))
	}
	return s
}

func do(r io.Reader) {
	plts, e := plot.DecodeAny(r)
	fatal(e)
	plts = dark(at(plts))
	pp(plot.AxisFromEnv(plts))
}
func pp(p plot.Plots) {
	if plt {
		fatal(p.Encode(os.Stdout))
		return
	}
	w, h := screensize()
	ip, e := p.IPlots(w, h, 0)
	fatal(e)
	m := plot.Image(ip, nil, w, h, 0).(*image.RGBA)
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
func dark(p plot.Plots) plot.Plots {
	if dst != FILE {
		for i := range p {
			p[i].Style.Dark = true
			p[i].Style.Order = "#00ff00"
		}
	}
	return p
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
