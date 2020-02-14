package main

import (
	"encoding/gob"
	"math"
	"os"
	"testing"

	"github.com/ktye/plot"
)

// Animation
func TestQuarters(t *testing.T) {
	ωdt := math.Pi / 30.0 // angular step
	var φ, ψ float64
	var z, zr, vdt complex128
	var ln []complex128
	for i := 0; i < 100; i++ {
		var l []complex128
		if i <= 30 {
			ψ = φ + math.Pi/4.0                        // cog angle
			zr = complex(ς*math.Cos(ψ), ς*math.Sin(ψ)) // cog position
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
		if φ > 2.0*math.Pi {
			φ -= 2.0 * math.Pi
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

	if false {
		writeplot(t, p)
	}
}

func writeplot(t *testing.T, p plot.Plot) {
	f, err := os.Create("quarters.plt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := gob.NewEncoder(f).Encode(plot.Plots{p}); err != nil {
		t.Fatal(err)
	}
}

func moveQuarter(z complex128, φ float64) []complex128 { // z: cog position
	l := newQuarter()
	e := complex(math.Cos(φ), math.Sin(φ))
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
		φ := math.Pi / 2.0 * float64(i) / float64(n-1)
		s, c := math.Sincos(φ)
		z[i+1] = complex(c, s)
	}
	return z
}

const ς float64 = 4.0 * math.Sqrt2 / math.Pi / 3.0 // 0.600...
