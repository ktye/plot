package plot

import (
	//"fmt"
	"bytes"
	"strconv"
	"strings"
	//"math"
	"encoding/binary"
)

func K(b []byte, x uint64, plts Plots, p int) (Plots, int, string) {
	br := func(x uint64, data interface{}) {
		binary.Read(bytes.NewReader(b[int32(x):]), binary.LittleEndian, data)
	}
	ux := func(p int32) uint64 { return binary.LittleEndian.Uint64(b[p:]) }
	ix := func(p int32) int32 { return int32(binary.LittleEndian.Uint32(b[p:])) }
	//fx := func(p int32) float64 { return math.Float64frombits(binary.LittleEndian.Uint64(b[p:])) }
	//zx := func(p int32) complex128 { return complex(fx(p), fx(p+8)) }
	nn := func(x uint64) int32 { return ix(int32(x) - 12) }
	CK := func(x uint64) []byte { return b[int32(x) : int32(x)+nn(x)] }
	IK := func(x uint64) []int32 { r := make([]int32, nn(x)); br(x, r); return r }
	UK := func(x uint64) []uint64 { r := make([]uint64, nn(x)); br(x, r); return r }
	FK := func(x uint64) []float64 { r := make([]float64, nn(x)); br(x, r); return r }
	ZK := func(x uint64) []complex128 { r := make([]complex128, nn(x)); br(x, r); return r }
	tp := func(x uint64) int { return int(x >> 59) }
	s0 := ix(0)
	sK := func(x uint64) string { return string(CK(ux(s0 + int32(x)))) }
	SK := func(x uint64) []string {
		v := IK(x)
		r := make([]string, len(v))
		for i, j := range v {
			r[i] = string(CK(ux(s0 + j)))
		}
		return r
	}
	floats := func(x uint64) []float64 {
		if tp(x) == 19 {
			v := IK(x)
			r := make([]float64, len(v))
			for i := range v {
				r[i] = float64(v[i])
			}
			return r
		} else {
			return FK(x)
		}
	}
	axis := func(s []string) { plts[p].Limits, _ = AxisFromString(strings.Join(s, ",")) }
	data := func(x []uint64) {
		if len(x) == 0 {
			return
		}
		var l Line
		var X []float64
		t0 := tp(x[0])
		if t0 == 19 || t0 == 21 {
			X = floats(x[0])
		} else if t0 == 22 {
			l.C = ZK(x[0])
			plts[p].Lines = append(plts[p].Lines, l)
			l = Line{}
		}
		push := func(n int) {
			if X != nil {
				if len(X) != n {
					return
				}
				x := make([]float64, len(X))
				copy(x, X)
				l.X = x
			}
			plts[p].Lines = append(plts[p].Lines, l)
			l = Line{}
		}
		for _, xi := range x[1:] {
			if t0 == 19 || t0 == 21 {
				l.Y = floats(xi)
				push(len(l.Y))
			} else if t0 == 22 {
				l.C = ZK(xi)
				push(len(l.C))
			}
		}
	}

	if len(plts) == 0 {
		plts = append(plts, Plot{})
	}

	tx, xi := tp(x), int(int32(x))
	switch tx {
	case 0: //reset
		return nil, 0, ""
	case 2:
		return plts, p, "todo: styles"
	case 3:
		if xi < 0 {
			return nil, 0, ""
		}
		if xi == len(plts) {
			p = xi
			plts = append(plts, Plot{})
		} else if xi < len(plts) {
			p = xi
		} else {
			return plts, p, "plot id does not exist"
		}
	case 4:
		s := sK(x) //todo equal
		plts[p].Type = PlotType(s)
	case 5:
		return plts, p, "todo rand float"
	case 6:
		return plts, p, "todo rand polar"
	case 18:
		return plts, p, "todo styles"
	case 19:
		v := IK(x)
		s := make([]string, len(v))
		for i := range v {
			s[i] = strconv.Itoa(int(v[i]))
		}
		axis(s)
	case 20:
		v := FK(x)
		s := make([]string, len(v))
		for i := range v {
			s[i] = strconv.FormatFloat(v[i], 'g', -1, 64)
		}
		axis(s)
	case 21:
		v := SK(x)
		if len(v) == 1 {
			plts[p].Title = v[0]
		} else if len(v) == 2 {
			plts[p].Xlabel, plts[p].Ylabel = v[0], v[1]
		} else if len(v) > 2 {
			plts[p].Xlabel, plts[p].Ylabel, plts[p].Title = v[0], v[1], v[2]
		}
	case 22:
		data([]uint64{21})
	case 23:
		data(UK(x))
	case 24:
		return plts, p, "todo key table plot"
	case 25:
		return plts, p, "todo table plot"
	default:
		return plts, p, "plot unknown k type"
	}
	return plts, p, ""
}

/*
// KTablePlot creates plots from k tables.
// It does not modify k memory.
//
// The mapping from data to plots/lines/points is the same as TextDataFormat.
// The header of x is the x-label and the column headers are plot titles.
//
// If the input is a list of tables, they must have the same structure and create multiline plots.
func KTablePlot(x uint32, C []byte, I []uint32, F []float64) (Plots, error) {
	var plts Plots
	if t := tp(x, I); t == 6 {
		n := nn(x, I)
		for i := uint32(0); i < n; i++ {
			data, err := kdata(I[2+i+x>>2], C, I, F)
			if err != nil {
				return nil, err
			}
			plts, err = newline(plts, data)
			if err != nil {
				return nil, err
			}
		}
		if plts == nil {
			return nil, fmt.Errorf("k-table data is empty")
		}
	} else if t == 7 {
		data, err := kdata(x, C, I, F)
		if err != nil {
			return nil, err
		}
		plts, err = newline(plts, data)
	} else {
		return nil, fmt.Errorf("input is not a dict-table or list thereof (%d)", t)
	}
	return plts, nil
}
func kdata(x uint32, C []byte, I []uint32, F []float64) ([]col, error) {
	if tp(x, I) != 7 {
		return nil, fmt.Errorf("data is not a k-table")
	}
	key := I[2+x>>2]
	val := I[3+x>>2]
	n := nn(val, I)
	if n == 0 {
		return nil, fmt.Errorf("table is empty")
	}
	m := nn(I[2+val>>2], I)
	var r []col
	for i := uint32(0); i < n; i++ {
		s := I[2+i+key>>2]
		y := I[2+i+val>>2]
		yn := nn(y, I)
		if yn != m {
			return nil, fmt.Errorf("dict is not uniform")
		}
		c, e := kcol(y, C, I, F, ksym(s, C, I))
		if e != nil {
			return nil, e
		}
		r = append(r, c)
	}
	return r, nil
}
func ksym(x uint32, C []byte, I []uint32) string {
	r := I[(I[132>>2]+x)>>2]
	n := nn(r, I)
	return string(C[8+r : 8+r+n])
}
func kcol(x uint32, C []byte, I []uint32, F []float64, s string) (c col, e error) {
	t := tp(x, I)
	n := nn(x, I)
	switch t {
	case 2:
		c.r = make([]float64, n)
		for i := uint32(0); i < n; i++ {
			c.r[i] = float64(I[2+i+x>>2])
		}
	case 3:
		c.r = make([]float64, n)
		for i := uint32(0); i < n; i++ {
			c.r[i] = F[1+i+x>>3]
		}
	case 4:
		c.c = make([]complex128, n)
		for i := uint32(0); i < n; i++ {
			c.c[i] = complex(F[1+2*i+x>>3], F[2+2*i+x>>3])
		}
		c.cmplx = true
	default:
		return c, fmt.Errorf("column data is not numeric")
	}
	c.s = s
	return c, nil
}
func tp(x uint32, I []uint32) uint32 {
	if x < 256 {
		return 0
	}
	return I[x>>2] >> 29
}
func nn(x uint32, I []uint32) (xn uint32) {
	if x < 256 {
		return 1
	}
	return I[x>>2] & 536870911
}
*/
