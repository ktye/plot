package xmath

import (
	"fmt"
	"math"
	"math/cmplx"
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
