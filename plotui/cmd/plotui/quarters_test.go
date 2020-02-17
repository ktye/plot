package main

import (
	"bytes"
	"io/ioutil"
	"math"
	"testing"

	"github.com/ktye/plot"
)

// Animation
func TestQuarters(t *testing.T) {
	ωdt := π / 30 // angular step
	var φ, ψ float64
	var z, zr, vdt complex128
	var ln []complex128
	for i := 0; i < 100; i++ {
		var l []complex128
		if i <= 30 {
			ψ = φ + π/4                                // cog angle
			zr = complex(ς*cos(ψ), ς*sin(ψ))           // cog position
			vdt = complex(ωdt, 0) * complex(0, 1) * zr // cog velocity * dt
			l = moveQuarter(zr, φ)
			z = zr
		} else {
			z += vdt
			l = moveQuarter(z, φ)
		}
		ln = append(ln, l...)
		ln = append(ln, complex(math.NaN(), math.NaN()))

		φ += ωdt
		if φ > 2.0*π {
			φ -= 2.0 * π
		}
	}

	p := plot.Plot{
		Type:    plot.Polar,
		Limits:  plot.Limits{Xmax: 3, Ymax: 3},
		Lines:   []plot.Line{plot.Line{C: ln, Segments: true}},
		Caption: &plot.Caption{LeadText: []string{"Animation of breaking rotor"}},
	}
	p.Lines[0].Style.Marker.Size = 3
	p.Lines[0].Style.Line.Width = 2
	plts := plot.Plots{p}

	var b bytes.Buffer
	if err := plts.Encode(&b); err != nil {
		t.Fatal(err)
	}

	if false {
		if err := ioutil.WriteFile("quarters.plt", b.Bytes(), 0744); err != nil {
			t.Fatal(err)
		}
	} else {
		if _, err := plot.DecodeAny(&b); err != nil {
			t.Fatal(err)
		} else if false {
			ioutil.WriteFile("quarters.plt", b.Bytes(), 0744)
		}
	}
}

func moveQuarter(z complex128, φ float64) []complex128 { // z: cog position
	l := newQuarter()
	e := complex(cos(φ), sin(φ))
	for i := range l {
		l[i] -= complex(ς, ς) / math.Sqrt2
		l[i] *= e
		l[i] += z
	}
	return l
}

func newQuarter() []complex128 { // unit quarter segment
	n := 12
	z := make([]complex128, n+2) // start and end point remains 0
	for i := 0; i < n; i++ {
		φ := π / 2 * float64(i) / float64(n-1)
		c, s := cosin(φ)
		z[i+1] = complex(c, s)
	}
	return z
}

func cos(φ float64) float64                  { return math.Cos(φ) }
func sin(φ float64) float64                  { return math.Sin(φ) }
func cosin(φ float64) (s float64, c float64) { s, c = math.Sincos(φ); return }

const π = math.Pi
const ς float64 = 4 * math.Sqrt2 / π / 3 // 0.600...
