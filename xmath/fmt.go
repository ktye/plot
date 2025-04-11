package xmath

import (
	"fmt"
	"io"
	"math"
	"math/cmplx"
	"strings"
)

// Absang formats a complex128 number in abs@ang format.
func Absang(x complex128, format string) string {
	if format == "" {
		format = "%v@%v"
	}
	r, phi := cmplx.Polar(x)
	phi *= 180.0 / math.Pi
	if phi < 0 {
		phi += 360.0
	}
	if r == 0.0 {
		phi = 0.0 // We want predictable angles in this case.
	}
	if phi == -0.0 || phi == 360.0 {
		phi = 0.0
	}
	return fmt.Sprintf(format, r, phi)
}
func ParseComplex(s string) (z complex128, err error) {
	s = strings.ReplaceAll(s, ",", ".")
	var amp, ang float64
	var extra string
	n, err := fmt.Sscanf(s, "%f@%f%s", &amp, &ang, &extra)
	if n != 2 || (err != nil && err != io.EOF) {
		return z, fmt.Errorf("cannot parse complex number: %s (example: 1.2@90)", s)
	}
	switch ang {
	case 0:
		z = complex(amp, 0)
	case 90:
		z = complex(0, amp)
	case 180:
		z = complex(-amp, 0)
	case 270:
		z = complex(0, -amp)
	default:
		z = cmplx.Rect(amp, ang/180*math.Pi)
	}
	return z, nil
}
