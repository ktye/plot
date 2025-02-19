package xmath

import "math/cmplx"

func Iota(n int) []float64 {
	f := make([]float64, n)
	for i := 0; i < n; i++ {
		f[i] = float64(i)
	}
	return f
}

// RealVector returns the real part of complex slices.
func RealVector(v []complex128) (re []float64) {
	re = make([]float64, len(v))
	for i, z := range v {
		re[i] = real(z)
	}
	return
}

// ImagVector returns the imag part of complex slices.
func ImagVector(v []complex128) (im []float64) {
	im = make([]float64, len(v))
	for i, z := range v {
		im[i] = imag(z)
	}
	return
}

// ComplexVector returns a complex slice from real and imag parts.
func ComplexVector(re, im []float64) []complex128 {
	v := make([]complex128, len(re))
	n := len(re)
	if len(im) < n {
		for i, r := range re {
			v[i] = complex(r, 0.0)
		}
		return v
	}
	for i := 0; i < n; i++ {
		v[i] = complex(re[i], im[i])
	}
	return v
}

// AbsVector returns the absolute values of complex slices.
func AbsVector(vec []complex128) (abs []float64) {
	abs = make([]float64, len(vec))
	for i, z := range vec {
		abs[i] = cmplx.Abs(z)
	}
	return abs
}

// PhaseVector returns the phase of complex slices.
func PhaseVector(vec []complex128) (phase []float64) {
	phase = make([]float64, len(vec))
	for i, z := range vec {
		phase[i] = cmplx.Phase(z)
	}
	return phase
}

// MinMax returns the min and max value of a slice.
func MinMax(x []float64) (min float64, max float64) {
	if len(x) < 1 {
		return
	}
	min = x[0]
	max = x[0]
	for _, v := range x {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return
}
