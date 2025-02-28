//go:build skip
// +build skip

// go run animation.go lines | xrootplot -al
// go run animation.go plots | xrootplot -ap
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "plots" {
		makePlots()
	} else {
		makeLines()
	}
}
func makeLines() { // two points circulating in polar diagrams
	p0, p1 := cmplx.Rect(1.2, math.Pi*30.0/180.0), cmplx.Rect(0.8, math.Pi*300.0/180.0)
	dr := 0.2
	n := 360
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Println()
		}
		s, c := math.Sincos(2.0 * math.Pi * float64(i) / float64(n-1))
		dp := complex(dr*c, dr*s)
		fmt.Printf("%s %s\n", absang(p0+dp), absang(p1-dp))
	}
}
func makePlots() { // traveling wave
	n := 30
	for l := 0; l < 2; l++ {
		if l > 0 {
			fmt.Println()
		}
		for i := 0; i < 100; i++ {
			x := 2.0 * math.Pi * float64(i) / 99.0
			fmt.Printf("%v", x)
			for f := 0; f < n; f++ {
				t := 2.0 * math.Pi * float64(f-10*l) / float64(n-1)
				fmt.Printf(" %v", math.Cos(x-t))
			}
			fmt.Println()
		}
	}
}
func absang(z complex128) string {
	p := cmplx.Phase(z) / math.Pi * 180.0
	if p < 0 {
		p += 360.0
	}
	return fmt.Sprintf("%va%v", cmplx.Abs(z), p)
}
