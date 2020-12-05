// +build ignore

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const help = `x axis generator ktye/plot/cmd/x.go
x [N]    env: xmin xmax N
`

func main() {
	n := atoi(os.Getenv("N"))
	args := os.Args[1:]
	if len(args) > 0 {
		if strings.HasPrefix(args[0], "-") {
			fmt.Fprintln(os.Stderr, help)
			os.Exit(1)
		}
		n = atoi(args[0])
	}
	if n <= 1 {
		n = 10
	}
	xmin := atof(os.Getenv("xmin"))
	xmax := atof(os.Getenv("xmax"))
	if xmax == xmin {
		xmin, xmax = 0, 1
	}
	d := (xmax - xmin) / float64(n-1)
	for i := 0; i < n; i++ {
		fmt.Println(xmin + float64(i)*d)
	}
}
func atoi(s string) int {
	if n, e := strconv.Atoi(s); e != nil {
		return 0
	} else {
		return n
	}
}
func atof(s string) float64 {
	if n, e := strconv.ParseFloat(s, 64); e != nil {
		return 0.0
	} else {
		return n
	}
}
