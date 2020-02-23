package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/ktye/plot"
	"github.com/ktye/plot/xmath"
)

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
	t := plot.PlotType(plts[0].Type)
	if t != plot.XY && t != plot.Polar && t != plot.Ring && t != plot.AmpAng {
		fatal(fmt.Errorf("unsupported plot type: %s", t))
	}
	n := len(plts[0].Lines)
	for _, p := range plts[1:] {
		if p.Type != t {
			fatal(fmt.Errorf("different plot types"))
		} else if len(p.Lines) != n {
			fatal(fmt.Errorf("differnt number of lines"))
		}
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Println()
		}
		line(plts, i)
	}
}
func line(plts plot.Plots, n int) {
	var tab [][]float64
	var cpx = false
	for i, p := range plts {
		l := p.Lines[n]
		if i == 0 && len(l.X) > 0 {
			tab = append(tab, l.X)
		}
		switch p.Type {
		case plot.XY, plot.Ring:
			tab = append(tab, l.Y)
		case plot.Polar, plot.AmpAng:
			tab = append(tab, xmath.AbsVector(l.C))
			z := xmath.PhaseVector(l.C)
			for k := range z {
				z[k] *= 180.0 / math.Pi
				if z[k] < 0 {
					z[k] += 360.0
				}
			}
			tab = append(tab, z)
			cpx = true
		default:
			fatal(fmt.Errorf("unknown plot type: %s", p.Type))
		}
	}
	if len(tab) == 0 {
		return
	}
	n = len(tab[0])
	f := bufio.NewWriter(os.Stdout)
	for i := 0; i < n; i++ {
		var s string
		for k, c := range tab {
			if len(c) != n {
				fatal(fmt.Errorf("data is not uniform"))
			}
			if k == 0 {
				s = ""
			} else if csv {
				s = ","
			} else if cpx && k%2 == n%2 {
				s = "a"
			} else {
				s = " "
			}
			fmt.Fprintf(f, s+"%v", c[i])
			if i == len(tab)-1 {
				fmt.Fprintf(f, "\n")
			}
		}
	}
	f.Flush()
}
