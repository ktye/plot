package plot

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/cmplx"
	"os"
	"strings"
	"testing"

	"github.com/ktye/plot/xmath"
)

func writeTest(b []byte, file string) {
	if true {
		os.WriteFile(file, b, 0744)
	}
}

func TestPlot(t *testing.T) {
	p := Plots{xy, polar, ampang, heatmap, envelope, bars}
	w, h := 1000, 600
	idx := []HighlightID{HighlightID{Line: 1, Point: -1}}

	b, e := p.Png(w, h, 3, idx)
	if e != nil {
		t.Fatal(e)
	}
	writeTest(b, "plot.png")

	b, e = p.Svg(w, h, 3, idx)
	if e != nil {
		t.Fatal(e)
	}
	writeTest(b, "plot.svg")

	b, e = p.Wmf(w, h, 3, idx)
	if e != nil {
		t.Fatal(e)
	}
	writeTest(b, "plot.wmf")

	b, e = p.Emf(w, h, 3, idx, "", 0, 0)
	if e != nil {
		t.Fatal(e)
	}
	writeTest(b, "plot.emf")

}
func TestEncDec(t *testing.T) {
	var buf bytes.Buffer
	p := Plots{xy, polar, ampang}
	if e := p.Encode(&buf); e != nil {
		t.Fatal(e)
	}
	s1 := string(buf.Bytes())

	pp, e := DecodePlots(strings.NewReader(s1))
	if e != nil {
		t.Fatal(e)
	}
	if len(pp) != len(p) {
		t.Fatalf("expected %d, got %d", len(p), len(pp))
	}

	buf.Reset()
	if e := pp.Encode(&buf); e != nil {
		t.Fatal(e)
	}
	s2 := string(buf.Bytes())
	if s1 != s2 {
		//t.Fatal("plots differ")
		t.Fatalf("first:\n%s\nnow:\n%s", s1, s2)
	}
	//fmt.Println(s2)
}

var x = linspace(0, 10, 100)

var xy Plot = Plot{
	Type:   XY,
	Title:  "XY",
	Xlabel: "x-axis",
	Ylabel: "y-axis",
	Xunit:  "μm",
	Yunit:  "kPa",
	Lines: []Line{
		Line{
			X: x,
			Y: apply(x, math.Sin),
		},
		Line{
			Id: 1,
			X:  x,
			Y:  apply(x, math.Cos),
		},
	},
}
var xy0 Plot = Plot{
	Type:   XY,
	Limits: Limits{Ymax: 1},
	Lines: []Line{
		Line{X: []float64{0, 1}, Y: []float64{1, 1}},
		Line{X: []float64{0, 1}, Y: []float64{.2, .2}, Id: 1},
		Line{X: []float64{0, 0}, Y: []float64{0, 1}, Id: 2},
		Line{X: []float64{.2, .2}, Y: []float64{0, 1}, Id: 2},
		Line{X: []float64{1, 1}, Y: []float64{0, 1}, Id: 2},
	},
}
var polar0 Plot = Plot{
	Type:   Polar,
	Limits: Limits{Ymax: 1},
	Lines: []Line{
		Line{C: []complex128{0}},
		Line{C: []complex128{0.2 + 0i, 0 + 0.2i, -0.2 + 0i, 0 - 0.2i}, Id: 1},
		Line{C: []complex128{1 + 0i, 0 + 1i, -1 + 0i, 0 - 1i}, Id: 2},
	},
}
var polar Plot = Plot{
	Type:   Polar,
	Title:  "Polar",
	Yunit:  "km/s",
	Limits: Limits{Ymax: 8},
	Lines: []Line{
		Line{
			X: x,
			C: spiral(x, 0, 1),
		},
		Line{
			Id: 1,
			X:  x,
			C:  spiral(x, math.Pi/2, 1),
		},
		Line{
			Id:    2,
			C:     []complex128{0, 4 + 2i},
			Style: DataStyle{Line: LineStyle{Width: 2, Arrow: 2}},
		},
		Line{
			Id:    3,
			C:     []complex128{0, 4 - 2i},
			Style: DataStyle{Line: LineStyle{Width: 2, Arrow: 2}},
		},
		Line{
			Id:    4,
			C:     []complex128{0, -4 + 2i},
			Style: DataStyle{Line: LineStyle{Width: 2, Arrow: 2}},
		},
		Line{
			Id:    5,
			C:     []complex128{0, -4 - 2i},
			Style: DataStyle{Line: LineStyle{Width: 2, Arrow: 2}},
		},
	},
}

var ampang Plot = Plot{
	Type:   AmpAng,
	Title:  "Amplitude/Angle",
	Xlabel: "Inclination",
	Ylabel: "Velocity",
	Xunit:  "°",
	Yunit:  "m/s",
	Lines: []Line{
		Line{
			X: x,
			C: spiral(x, 0, 1),
		},
		Line{
			Id: 1,
			X:  x,
			C:  spiral(x, math.Pi/2, 1.5),
		},
	},
}

var foto Plot = Plot{
	Type:  Foto,
	Title: "Setup",
	Foto:  painting(),
}

var heatmap Plot = Plot{
	Type:   Raster,
	Title:  "Heatmap",
	Xlabel: "Force",
	Xunit:  "N",
	Ylabel: "Pressure",
	Yunit:  "kPa",
	Lines: []Line{
		Line{
			X:     x,
			Y:     x,
			Image: x_plus_y(x),
		},
	},
}

var envelope Plot = Plot{
	Type:   XY,
	Title:  "Envelope plot äüöß⍳↓→λφα",
	Xlabel: "time",
	Xunit:  "s",
	Ylabel: "quantity",
	Limits: Limits{Xmin: 1.1, Xmax: 6.5, Ymin: -2, Ymax: 2},
	Lines: []Line{
		Line{
			X: []float64{1, 2, 3, 4, 5, 6, 7},
			Y: []float64{1, -1, 2, -2, 3, -3, 3},
		},
		Line{
			X: []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			C: []complex128{0.6 + 0.4i, 0.7 + 0.3i, 0.6 + 0.4i, 0.7 + 0.3i, 0.6 + 0.4i, 0.7 + 0.3i, 0.6 + 0.4i, 0.7 + 0.3i, 0.6 + 0.4i, 0.8 + 0.2i},
		},
	},
}

var bars Plot = Plot{
	Type:   XY,
	Xlabel: "steps",
	Ylabel: "count",
	Lines: []Line{
		Line{
			Id:    1,
			X:     []float64{0.1, 1, 1.1, 2, 2.1, 3, 3.1, 4, 4.1, 5},
			Y:     []float64{0, 1, 0, 2, 0, 1, 0, 3, 0, 1},
			Style: DataStyle{Marker: MarkerStyle{Marker: Bar, Size: 1}},
		},
		Line{
			Id:    2,
			X:     []float64{0.1, 1, 1.1, 2, 2.1, 3, 3.1, 4, 4.1, 5},
			Y:     []float64{1.1, 3, 2.1, 4, 1.1, 3, 3.1, 4, 1.1, 3},
			Style: DataStyle{Marker: MarkerStyle{Marker: Bar, Size: 1}},
		},
	},
}

func linspace(min, max float64, n int) []float64 {
	y := make([]float64, n)
	for i := range y {
		y[i] = xmath.Scale(float64(i), 0, float64(len(y)-1), min, max)
	}
	return y
}

func apply(x []float64, f func(float64) float64) []float64 {
	y := make([]float64, len(x))
	for i, v := range x {
		y[i] = f(v)
	}
	return y
}

func spiral(r []float64, phi0, scale float64) []complex128 {
	z := make([]complex128, len(r))
	for i := range z {
		z[i] = cmplx.Rect(scale*r[i], r[i]+phi0)
	}
	return z
}

func painting() string {
	var r image.Rectangle
	r.Max = image.Point{30, 30}
	m := image.NewRGBA(r)
	draw.Draw(m, m.Bounds(), &image.Uniform{color.RGBA{0, 0, 0xFF, 0xFF}}, image.ZP, draw.Src)
	s, _ := EncodeToPng(m)
	return s
}

func x_plus_y(x []float64) [][]uint8 {
	m := make([][]uint8, len(x))
	for i := range m {
		m[i] = make([]uint8, len(x))
		for k := range x {
			m[i][k] = uint8(255 * (x[i] + x[k]) / 20)
		}
	}
	return m
}
