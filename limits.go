package plot

import (
	"fmt"
	"math"
	"math/cmplx"
	"strconv"
	"strings"
)

// Limits are axes limits for a single plot.
type Limits struct {
	Equal                              bool
	Xmin, Xmax, Ymin, Ymax, Zmin, Zmax float64
}

func (l Limits) String() string {
	if l.Equal {
		return "equal"
	}
	v := []float64{l.Xmin, l.Xmax, l.Ymin, l.Ymax, l.Zmin, l.Zmax}
	s := make([]string, len(v))
	for i, f := range v {
		s[i] = strconv.FormatFloat(f, 'g', -1, 64)
	}
	return strings.Join(s, ",")
}

// EqualLimits returns limits which include all given plots.
func (plts Plots) EqualLimits() (Limits, error) {
	setLimits := func(elMin, elMax, min, max float64) (float64, float64) {
		if elMin == elMax {
			return min, max
		}
		if min == max {
			return elMin, elMax
		}
		if min < elMin {
			elMin = min
		}
		if max > elMax {
			elMax = max
		}
		return elMin, elMax
	}

	var el, l Limits
	for _, p := range plts {
		switch p.Type {
		case Polar:
			l = p.getPolarLimits(false)
		case Ring:
			l = p.getPolarLimits(true)
		case AmpAng:
			l = p.getAmpAngLimits()
		case XY, Raster:
			l = p.getXYLimits()
		default:
			return Limits{}, fmt.Errorf("cannot get equal limits: plot type '%s' is not implemented", p.Type)
		}

		el.Xmin, el.Xmax = setLimits(el.Xmin, el.Xmax, l.Xmin, l.Xmax)
		el.Ymin, el.Ymax = setLimits(el.Ymin, el.Ymax, l.Ymin, l.Ymax)
		el.Zmin, el.Zmax = setLimits(el.Zmin, el.Zmax, l.Zmin, l.Zmax)
	}
	return el, nil
}

// getXYLimits returns the set limits for an XY Plot
// or calculates default limits, if no limits have been given by the user.
func (p *Plot) getXYLimits() Limits {
	limits := Limits{false, p.Xmin, p.Xmax, p.Ymin, p.Ymax, 0, 0}
	if p.Xmin == p.Xmax {
		a := autoscale{}
		for _, l := range p.Lines {
			a.add(l.X)
		}
		limits.Xmin, limits.Xmax, _ = a.niceLimits()
	}
	if p.Ymin == p.Ymax {
		a := autoscale{}
		for _, l := range p.Lines {
			if l.Y == nil {
				a.addEnvelope(l.C)
			} else {
				a.add(l.Y)
			}
		}
		limits.Ymin, limits.Ymax, _ = a.niceLimits()
	}
	if p.Type == Raster && len(p.Lines) > 0 {
		limits.Zmin = p.Lines[0].ImageMin
		limits.Zmax = p.Lines[0].ImageMax
	}
	return limits
}

// getPolarLimits returns the limits for polar or ring plots.
func (p *Plot) getPolarLimits(ring bool) Limits {
	limits := Limits{false, p.Xmin, p.Xmax, p.Ymin, p.Ymax, p.Zmin, p.Zmax}
	if ring {
		// user defines rmin for a ring plot as ymin, but it is store in Zmin internally.
		limits.Zmin = p.Ymin
	}
	limits.Ymin = 0
	rmin := 0.0
	if p.Ymax == 0 {
		a := autoscale{}
		if ring {
			for _, l := range p.Lines {
				a.add(l.X) // radial value for a ring plot
			}
			rmin, limits.Ymax, _ = a.niceLimits()
		} else {
			for _, l := range p.Lines {
				a.addComplex(l.C)
			}
			_, limits.Ymax, _ = a.niceLimits()
		}
	}
	r := limits.Ymax
	limits.Xmin = -r
	limits.Xmax = r
	limits.Ymin = -r
	limits.Zmin = rmin // we store rmin as zmin for a ring plot
	limits.Zmax = r
	return limits
}

func (l Limits) isPolarLimits() bool {
	if r := l.Xmax; l.Xmin == -r && l.Ymin == -r && l.Ymax == r {
		return true
	}
	return false
}

// getAmpAngLimits returns speed and amplitude limits.
func (p *Plot) getAmpAngLimits() Limits {
	limits := Limits{false, p.Xmin, p.Xmax, p.Ymin, p.Ymax, 0, 0}
	if p.Xmin == p.Xmax {
		a := autoscale{}
		for _, l := range p.Lines {
			a.add(l.X)
		}
		limits.Xmin, limits.Xmax, _ = a.niceLimits()
	}
	limits.Ymin = 0
	if p.Ymax == 0 {
		a := autoscale{}
		for _, l := range p.Lines {
			a.addComplex(l.C)
		}
		_, limits.Ymax, _ = a.niceLimits()
	}
	return limits
}

type autoscale struct {
	min, max float64
	isInit   bool
}

func (a *autoscale) add(x []float64) {
	for _, v := range x {
		if a.isInit == false {
			a.min = v
			a.max = v
			a.isInit = true
		}
		if v > a.max {
			a.max = v
		}
		if v < a.min {
			a.min = v
		}
	}
}

func (a *autoscale) addEnvelope(e []complex128) {
	for _, c := range e {
		if a.isInit == false {
			a.min = real(c)
			a.max = imag(c)
			a.isInit = true
		}
		if min := real(c); min < a.min {
			a.min = min
		}
		if max := imag(c); max > a.max {
			a.max = max
		}
	}
}

func (a *autoscale) addComplex(z []complex128) {
	for _, v := range z {
		r := cmplx.Abs(v)
		if math.IsNaN(r) == false {
			if a.isInit == false {
				a.min = 0
				a.max = r
				a.isInit = true
			}
			if r > a.max {
				a.max = r
			}
		}
	}
}

// uplimits returns a rounded upper limit for the range [0, a.max] used for polar or amplitude diagrams.
func (a *autoscale) uplimits() float64 {
	x := math.Abs(a.max)
	if x == 0 {
		return 1
	}
	p := math.Pow(10.0, math.Ceil(math.Log10(x)))
	if x < p/5 {
		return p / 5
	} else if x < p/2 {
		return p / 2
	} else if x < p/10 {
		return p / 10
	}
	return p
}

func (a *autoscale) niceLimits() (niceMin, niceMax, spacing float64) {
	return NiceLimits(a.min, a.max)
}

func (a *autoscale) getTics() Tics {
	return NiceTics(a.min, a.max)
}
