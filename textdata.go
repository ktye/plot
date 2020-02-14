package plot

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"strconv"
)

const TextDataFormat = `
y                 1 real column: x over 0..n-1
x y1 y2 y3..      n real columns: n-1 plots side-by-side
x c1 c2 c3..      1 real + n complex columns: n plots amplitude/phase over x
c1 c2 c3..        n complex columns: n polar plots
blank line:       new dataset (new line)

complex number: 1.2a30 (polar, degree)`

// TextDataPlot returns Plots from textual numeric data as described by TextDataFormat.
// e.g. xrootplot/
func TextDataPlot(r io.Reader) (Plots, error) {
	var plts Plots
	b, _ := ioutil.ReadAll(r)
	b = bytes.Trim(b, "\r\n\t ")
	lines := bytes.Split(b, []byte{'\n'})

	var data []col
	var e error
	for _, l := range lines {
		// lino++ ?
		cols := bytes.Fields(l)
		if n := len(cols); data == nil {
			data = make([]col, n)
		} else if n == 0 {
			plts, e = newline(plts, data)
			if e != nil {
				return nil, e
			}
			data = make([]col, len(data))
		} else if n != len(data) {
			return nil, fmt.Errorf("data is not uniform")
		}
		for k, c := range cols {
			data[k].push(c)
		}
	}
	return newline(plts, data)
}
func newline(plts Plots, data []col) (Plots, error) {
	var t PlotType
	var x []float64
	if len(data) < 1 {
		return nil, fmt.Errorf("data has no columns")
	}
	if len(data) == 1 && data[0].cmplx == false { // y
		data = []col{col{r: til(len(data[0].r))}, data[0]}
	}
	if data[0].cmplx {
		t = Polar
	} else if len(data) > 1 && data[1].cmplx {
		t = AmpAng
		x = data[0].r
		data = data[1:]
	} else {
		t = XY
		x = data[0].r
		data = data[1:]
	}
	if len(plts) == 0 {
		plts = make(Plots, len(data))
	} else if plts[0].Type != t {
		return nil, fmt.Errorf("plot types differ (not supported)")
	} else if len(plts) != len(data) {
		return nil, fmt.Errorf("data is not uniform")
	}

	for i := range plts {
		plts[i].Type = t
		var l Line
		id := len(plts[i].Lines)
		if t == XY {
			l = Line{X: xcp(x), Y: data[i].r, Id: id}
		} else {
			l = Line{X: xcp(x), C: data[i].c, Id: id}
		}
		plts[i].Lines = append(plts[i].Lines, l)
	}
	return plts, nil
}

/*
func pp(p plot.Plots) {
	w, h := screensize()
	ip, e := p.IPlots(w, h, 0)
	if e != nil {
		xx(e.Error())
	}
	m := plot.Image(ip, nil, w, h, 0).(*image.RGBA)
	draw(w, h, m.Pix)
}
*/

type col struct {
	r     []float64
	c     []complex128
	cmplx bool
}

func til(n int) []float64 {
	x := make([]float64, n)
	for i := range x {
		x[i] = float64(i)
	}
	return x
}
func xcp(x []float64) []float64 {
	r := make([]float64, len(x))
	copy(r, x)
	return r
}
func (c *col) push(s []byte) (e error) {
	if i := bytes.IndexByte(s, 'a'); i == -1 {
		var f float64
		f, e = parse(s, e)
		c.r = append(c.r, f)
	} else {
		var r, p float64
		r, e = parse(s[:i], e)
		p, e = parse(s[i+1:], e)
		im, re := math.Sincos(p * 180.0 / math.Pi)
		c.c = append(c.c, complex(r*re, r*im))
		c.cmplx = true
	}
	return e
}
func parse(s []byte, e error) (float64, error) {
	if e != nil {
		return 0, e
	}
	n, e := strconv.ParseFloat(string(s), 64)
	if e != nil {
		return 0, fmt.Errorf("cannot parse number: %q", string(s))
	}
	return n, nil
}
