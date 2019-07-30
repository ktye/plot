package plot

import (
	"math"

	"github.com/ktye/plot/xmath"
)

// An xyer should return x and y data from a line object.
// x and y coordinate are to be interpreted in the mathematical sense,
// not in image coordinates.
// The interface is used by axes.
type xyer interface {
	XY(Line) ([]float64, []float64, bool)
}

// xyPolar returns x and y data from a line object for a polar plot.
// Complex phase starts at the y axis and 90 degree is on the x axis.
// It is used for polar and ring plots. Polar plots ignore rmin and rmax.
type xyPolar struct {
	rmin, rmax float64
}

const innerRing = 0.5 // inner ring display ratio to outer ring

func (p xyPolar) XY(l Line) (x, y []float64, isEnvelope bool) {
	if p.rmin == 0 { // polar coordinates use l.C
		x = xmath.ImagVector(l.C)
		y = xmath.RealVector(l.C)
		return x, y, false
	} else { // ring coordinates r, phi are stored in l.X and l.Y
		// As l.X may be negative, this cannot be stored as a complex number.
		r := make([]float64, len(l.X))
		copy(r, l.X)
		θ := make([]float64, len(l.Y))
		copy(θ, l.Y)
		for i := range r {
			if r[i] < p.rmin {
				r[i] = p.rmin
			} else if r[i] > p.rmax {
				r[i] = p.rmax
			} else {
				r[i] = xmath.Scale(r[i], p.rmin, p.rmax, innerRing*p.rmax, p.rmax)
			}
		}
		x = make([]float64, len(r))
		y = make([]float64, len(r))
		for i := range r { // len(l.X) must be len(l.Y)
			x[i], y[i] = math.Sincos(θ[i])
			x[i] *= r[i]
			y[i] *= r[i]
		}
		return x, y, false
	}
}

// xyXY return x and y data for an XY plot.
type xyXY struct{}

func (p xyXY) XY(l Line) (x, y []float64, isEnvelope bool) {
	// Return envelope data.
	if l.Y == nil && len(l.C) > 0 {
		n := len(l.C)
		x = make([]float64, 2*n+1)
		y = make([]float64, 2*n+1)
		for i := 0; i < n; i++ {
			x[i] = l.X[i]
			y[i] = real(l.C[i])
		}
		k := n
		for i := n; i < 2*n; i++ {
			k--
			x[i] = l.X[k]
			y[i] = imag(l.C[k])
		}
		x[2*n] = l.X[0]
		y[2*n] = real(l.C[0])
		return x, y, true
	}
	return l.X, l.Y, false
}

// xyAmp returns x and y data for an amplitude plot.
type xyAmp struct{}

func (p xyAmp) XY(l Line) (x, y []float64, isEnvelope bool) {
	return l.X, xmath.AbsVector(l.C), false
}

// xyAng returns x and y data for a phase plot (-180,180).
// It inserts NaNs after phase jumps avoid vertical lines and to draw
// a new line instead.
type xyAng struct{}

func (p xyAng) XY(l Line) (x, y []float64, isEnvelope bool) {
	phase := xmath.PhaseVector(l.C)
	for i, v := range phase {
		phase[i] = v * 180.0 / math.Pi
	}
	var idx []int
	for i := 1; i < len(phase); i++ {
		if dp := math.Abs(phase[i] - phase[i-1]); dp > 300 {
			idx = append(idx, i)
		}
	}
	if len(idx) == 0 {
		return l.X, phase, false
	}
	x = make([]float64, len(l.X)+len(idx))
	y = make([]float64, len(phase)+len(idx))
	nextIdx := idx[0]
	iidx := 0
	off := 0
	for i := 0; i < len(phase); i++ {
		if i == nextIdx {
			x[i+off] = math.NaN()
			y[i+off] = math.NaN()
			off++
			if iidx < len(idx)-1 {
				iidx++
				nextIdx = idx[iidx]
			} else {
				nextIdx = len(phase) + 1
			}
		}
		x[i+off] = l.X[i]
		y[i+off] = phase[i]
	}
	return x, y, false
}
