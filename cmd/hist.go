// +build ignore
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/ktye/plot"
)

var line int

func main() {
	min, max, bins := 0.0, 0.0, 60
	a := os.Args[1:]
	if len(a) == 1 && (a[0] == "-h" || a[0] == "--help") {
		fmt.Println("hist [min max] N < data  (ktye/plot/cmd/hist.go)")
		return
	}
	if len(a) > 1 {
		min, max = atof(a[0]), atof(a[1])
	}
	if len(a) > 2 {
		bins = int(atof(a[2]))
	}
	if len(a) == 1 {
		bins = int(atof(a[0]))
	}
	var data []float64
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line++
		t := strings.TrimSpace(s.Text())
		if len(t) == 0 {
			continue
		}
		data = append(data, atof(t))
	}
	if len(data) == 0 {
		return
	}
	if min == max {
		min, max = data[0], data[0]
		for _, u := range data {
			if u < min {
				min = u
			}
			if u > max {
				max = u
			}
		}
	}
	spacing := 0.0
	min, max, spacing = plot.NiceLimits(min, max)
	dx := (max - min) / float64(bins-1)
	m := make([]int, bins)
	hmax := 0
	for _, u := range data {
		i := int((u - min) / dx)
		if i < 0 {
			i = 0
		}
		if i > bins-1 {
			i = bins - 1
		}
		m[i]++
		if m[i] > hmax {
			hmax = m[i]
		}
	}
	const height = 16
	M := make([][]byte, height)
	for i := range M {
		M[i] = make([]byte, bins)
	}
	for i := range m {
		h := float64(height) * float64(m[i]) / float64(hmax)
		a := int(math.Floor(h))
		b := int(8 * (h - math.Floor(h)))
		if a+b == 0 && h != 0 {
			b = 1
		}
		for j := 0; j < a; j++ {
			M[height-1-j][i] = 8
		}
		if a < height {
			M[height-a-1][i] = byte(b)
		}
	}
	u := [][]byte{[]byte("▁"), []byte("▂"), []byte("▃"), []byte("▄"), []byte("▅"), []byte("▆"), []byte("▇"), []byte("█")}
	for i := range M {
		s := bytes.ReplaceAll(M[i], []byte{0}, []byte{32})
		for j := byte(1); j < 9; j++ {
			s = bytes.ReplaceAll(s, []byte{j}, u[j-1])
		}
		fmt.Println(string(s))
	}

	l := ""
	x := min
	j := 0.0
	for i := 0; i < bins; i++ {
		if x >= min+j*spacing {
			fmt.Print("┬")
			s := fmt.Sprintf("%v", min+j*spacing)
			if r := i - len(l) - len(s)/2; r > 0 {
				l += strings.Repeat(" ", r)
			} else if i > 0 {
				l += " "
			}
			l += s
			j += 1.0
		} else {
			fmt.Print("─")
		}
		x += dx
	}
	fmt.Println("\n" + l)
}
func atof(s string) float64 {
	s = strings.ReplaceAll(s, ",", ".")
	f, e := strconv.ParseFloat(s, 64)
	fatal(e)
	return f
}
func fatal(e error) {
	if e != nil {
		panic(fmt.Errorf("line %d: %v", line, e))
	}
}
