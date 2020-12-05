// +build ignore
package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const help = `random numbers ktye/plot/cmd/r.go
r       N(env) float [-1,1]
r 3     N=3
r n     normal distribution(float)
r 3n    N=3 normal
r 100z  complex binormal 1.2a34
r 3n-   append as column to stdin
`

func main() {
	rand.Seed(time.Now().UnixNano())
	n, pipe, u, z := atoi(os.Getenv("N")), false, true, false
	args := os.Args[1:]
	if len(args) > 0 {
		a := args[0]
		if strings.HasSuffix(a, "-") {
			pipe = true
			a = a[:len(a)-1]
		}
		if strings.HasPrefix(a, "-") {
			fmt.Fprintf(os.Stderr, help)
			os.Exit(1)
		}
		if strings.HasSuffix(a, "n") {
			u = false
			a = a[:len(a)-1]
		}
		if strings.HasSuffix(a, "a") || strings.HasSuffix(a, "z") {
			u, z = false, true
			a = a[:len(a)-1]
		}
		if len(a) > 0 {
			n = atoi(a)
		}
	}
	if n <= 0 && pipe == false {
		os.Exit(1)
	}
	if pipe {
		scn := bufio.NewScanner(os.Stdin)
		for scn.Scan() {
			os.Stdout.WriteString(scn.Text() + " ")
			out(u, z)
		}
	} else {
		for i := 0; i < n; i++ {
			out(u, z)
		}
	}
}
func atoi(s string) int {
	if n, e := strconv.Atoi(s); e != nil {
		return 0
	} else {
		return n
	}
}
func out(u, z bool) {
	var x, y float64
	if u {
		x = -1.0 + 2.0*rand.Float64()
	} else if z {
		x = rand.NormFloat64()
		y = rand.NormFloat64()
		r := math.Hypot(x, y)
		p := math.Atan2(y, x) * 180.0 / math.Pi
		if p < 0 {
			p += 360.0
		}
		fmt.Printf("%va%v\n", r, p)
		return
	} else {
		x = rand.NormFloat64()
	}
	fmt.Println(x)
}
