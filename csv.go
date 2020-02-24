package plot

import (
	"bufio"
	"fmt"
	"io"
	"math"

	"github.com/ktye/plot/xmath"
)

func (plts Plots) WriteCsv(w io.Writer, comma bool) error {
	if len(plts) == 0 {
		return fmt.Errorf("csv: plot is empty")
	}
	t := PlotType(plts[0].Type)
	if t != XY && t != Polar && t != Ring && t != AmpAng {
		return fmt.Errorf("csv: unsupported plot type: %s", t)
	}
	n := len(plts[0].Lines)
	for _, p := range plts[1:] {
		if p.Type != t {
			return fmt.Errorf("different plot types")
		} else if len(p.Lines) != n {
			return fmt.Errorf("different number of lines")
		}
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprintln(w)
		}
		if e := csvLine(w, plts, i, comma); e != nil {
			return e
		}
	}
	return nil
}
func csvLine(w io.Writer, plts Plots, n int, comma bool) error {
	var tab [][]float64
	var cpx = false
	for i, p := range plts {
		l := p.Lines[n]
		if i == 0 && len(l.X) > 0 {
			tab = append(tab, l.X)
		}
		switch p.Type {
		case XY, Ring:
			tab = append(tab, l.Y)
		case Polar, AmpAng:
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
			return fmt.Errorf("unknown plot type: %s", p.Type)
		}
	}
	if len(tab) == 0 {
		return fmt.Errorf("csv: empty plot")
	}
	n = len(tab[0])
	f := bufio.NewWriter(w)
	for i := 0; i < n; i++ {
		var s string
		for k, c := range tab {
			if len(c) != n {
				return fmt.Errorf("data is not uniform")
			}
			if k == 0 {
				s = ""
			} else if comma {
				s = ","
			} else if cpx && k%2 == n%2 {
				s = "a"
			} else {
				s = " "
			}
			fmt.Fprintf(f, s+"%v", c[i])
			if k == len(tab)-1 {
				fmt.Fprintf(f, "\n")
			}
		}
	}
	f.Flush()
	return nil
}
