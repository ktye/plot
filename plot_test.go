package plot

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"math/cmplx"
	"strings"
	"testing"

	"github.com/ktye/plot/xmath"
)

func TestPlot(t *testing.T) {
	p := Plots{xy, polar, ampang, heatmap}
	w, h := 800, 400
	ip, err := p.IPlots(w, h, 0)
	if err != nil {
		t.Fatal(err)
	}
	im := Image(ip, nil, w, h, 0)

	var buf bytes.Buffer
	if err := png.Encode(&buf, im); err != nil {
		t.Fatal(err)
	}

	/*
		if err := ioutil.WriteFile("out.png", buf.Bytes(), 0644); err != nil {
			t.Fatal(err)
		}
	*/
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

var polar Plot = Plot{
	Type:  Polar,
	Title: "Polar",
	Yunit: "km/s",
	Lines: []Line{
		Line{
			X: x,
			C: spiral(x, 0),
		},
		Line{
			Id: 1,
			X:  x,
			C:  spiral(x, math.Pi/2),
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
			C: spiral(x, 0),
		},
		Line{
			Id: 1,
			X:  x,
			C:  spiral(x, math.Pi/2),
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

func spiral(r []float64, phi0 float64) []complex128 {
	z := make([]complex128, len(r))
	for i := range z {
		z[i] = cmplx.Rect(r[i], r[i]+phi0)
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
