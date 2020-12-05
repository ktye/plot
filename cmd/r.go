package main // +build ignore
import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const help = `random numbers /c/k/plot/cmd/r.go
r       N(env) float [-1,1]
r 3     N=3
r n            float normal distribution
r 3n
r 100z         complex binormal
r -            (last) cat col to stdin
`

func main() {
	n, pipe, u, z := atoi(os.Getenv("N")), false, true, false
	args := os.Args[1:]
	if len(args) > 0 && args[len(args)-1] == "-" {
		pipe = true
		args = args[:len(args)-1]
	}
	if len(args) > 0 {
		a := args[0]
		if strings.HasPrefix(a, "-") {
			fmt.Fprintf(os.Stderr, help)
			os.Exit(1)
		}
		if strings.HasSuffix(a, "n") {
			u = false
			a = a[:len(a)-1]
		}
		if strings.HasSuffix(a, "a") || strings.HasSuffix(a, "z") {
			z = true
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
