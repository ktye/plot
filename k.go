package plot

import (
	"bytes"
	"fmt"
	"math"
	"math/rand/v2"
	"strconv"
	"strings"

	//"math"
	"encoding/binary"

	"github.com/ktye/plot/xmath"
)

// k-interface (readonly, does not unref)
// type
// c style(todo)
// i plot id, e.g. 1 to add a 2nd plot
// s plot type `polar `xy .. (todo `equal`bar`stacked)
// f
// z
// C
// S
// I
// F
// Z
// L
// D
// T
func K(b []byte, x uint64, plts Plots, p int) (Plots, int, string) {
	br := func(x uint64, data interface{}) {
		binary.Read(bytes.NewReader(b[int32(x):]), binary.LittleEndian, data)
	}
	ux := func(p int32) uint64 { return binary.LittleEndian.Uint64(b[p:]) }
	ix := func(p int32) int32 { return int32(binary.LittleEndian.Uint32(b[p:])) }
	fK := func(x uint64) float64 { return math.Float64frombits(binary.LittleEndian.Uint64(b[int32(x):])) }
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
			r[i] = string(CK(ux(s0 + int32(j))))
		}
		return r
	}
	floats := func(x uint64) []float64 {
		if t := tp(x); t == 19 {
			v := IK(x)
			r := make([]float64, len(v))
			for i := range v {
				r[i] = float64(v[i])
			}
			return r
		} else if t == 21 {
			return FK(x)
		} else {
			return nil
		}
	}
	axis := func(s []string) { plts[p].Limits, _ = AxisFromString(strings.Join(s, ",")) }
	pushline := func(l Line) {
		l.Id = len(plts[p].Lines)
		plts[p].Lines = append(plts[p].Lines, l)
	}
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
			pushline(l)
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
			pushline(l)
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
	fzdata := func(x []float64, z []complex128) {
		l := Line{}
		if x == nil {
			l.C = z
			plts[p].Type = Polar
		} else {
			l.X = xmath.Iota(len(x))
			l.Y = x
			plts[p].Type = XY
		}
		pushline(l)
	}
	if len(plts) == 0 {
		plts = append(plts, Plot{})
	}
	group := func(x []int32) (r []int) { //end-index when values change (single pass)
		if len(x) == 0 {
			return nil
		}
		r = append(r, 0)
		y := x[0]
		for i := range x {
			if x[i] != y {
				r = append(r, i)
				y = x[i]
			}
		}
		return append(r, len(x))
	}
	table := func(keys []string, vals []uint64, ids []int32) (Plots, int, string) {
		if len(keys) == 0 {
			return plts, p, ""
		}
		x := floats(vals[0])
		xlabel := ""
		if x != nil {
			xlabel = keys[0]
			keys = keys[1:]
			vals = vals[1:]
		}
		if len(vals) == 0 {
			return plts, p, ""
		}
		xx := func(a, b int) (r []float64) {
			if x == nil {
				return x
			}
			r = make([]float64, b-a)
			copy(r, x[a:b])
			return r
		}
		g := group(ids)
		for i, v := range vals {
			p = i
			if p == len(plts) {
				plts = append(plts, Plot{})
			}
			plts[p].Title = keys[i]
			plts[p].Xlabel = xlabel
			y := floats(v)
			if y == nil {
				if tp(v) != 22 {
					return plts, p, "plot table column types must be numeric"
				}
				z := ZK(v)
				for j := 0; j < len(g)-1; j++ {
					pushline(Line{X: xx(g[j], g[1+j]), C: z[g[j]:g[1+j]]})
				}
			} else {
				for j := 0; j < len(g)-1; j++ {
					pushline(Line{X: xx(g[j], g[1+j]), Y: y[g[j]:g[1+j]]})
				}
			}
		}
		return plts, p, ""
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
		fmt.Println("s?", len(s), s)
		plts[p].Type = PlotType(s)
	case 5:
		n := int(fK(x))
		var f []float64
		if n < 0 {
			f = make([]float64, -n)
			for i := range f {
				f[i] = rand.NormFloat64()
			}
		} else {
			f = make([]float64, n)
			for i := range f {
				f[i] = rand.Float64()
			}
		}
		fzdata(f, nil)
	case 6:
		n := int(fK(x))
		if n > 0 {
			z := make([]complex128, n)
			for i := range z {
				z[i] = complex(rand.NormFloat64(), rand.NormFloat64())
			}
			fzdata(nil, z)
		}
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
		v := SK(x)
		fmt.Println("SK", v)
		if len(v) == 1 {
			plts[p].Title = v[0]
		} else if len(v) == 2 {
			plts[p].Xlabel, plts[p].Ylabel = v[0], v[1]
		} else if len(v) > 2 {
			plts[p].Xlabel, plts[p].Ylabel, plts[p].Title = v[0], v[1], v[2]
		}
	case 21:
		v := FK(x)
		s := make([]string, len(v))
		for i := range v {
			s[i] = strconv.FormatFloat(v[i], 'g', -1, 64)
		}
		axis(s)
	case 22:
		data([]uint64{x})
	case 23:
		data(UK(x))
	case 24:
		keys := ux(int32(x))
		vals := ux(int32(x) + 8)
		if tp(vals) != 25 {
			return plts, p, "plot: dict must be a key table"
		}
		var g []int32
		if tp(keys) == 19 || tp(keys) == 20 {
			g = IK(keys)
		} else if tp(keys) == 25 {
			x := ux(int32(keys) + 8)
			if tp(x) != 23 {
				return plts, p, "plot: key table expected list"
			}
			if nn(x) != 1 {
				return plts, p, "plot: key table should have 1 key column"
			}
			x = ux(int32(x))
			if tp(x) != 19 && tp(x) != 20 {
				return plts, p, "plot: key table column type should be I or S"
			}
			g = IK(keys)
		}
		return table(SK(ux(int32(vals))), UK(ux(int32(vals)+8)), g)
	case 25:
		keys := ux(int32(x))
		vals := ux(int32(x) + 8)
		return table(SK(keys), UK(vals), make([]int32, nn(x))) //single line
	default:
		return plts, p, "plot unknown k type"
	}
	return plts, p, ""
}
func Kcaption(b []byte, x uint64, plts Plots) (Plots, string) { // t(table) S(units) I(decimals) i(decimals) Ii(neg:angle) C(leadtext)
	if len(plts) == 0 {
		plts = append(plts, Plot{})
	}
	tp := func(x uint64) int { return int(x >> 59) }
	br := func(x uint64, data interface{}) {
		binary.Read(bytes.NewReader(b[int32(x):]), binary.LittleEndian, data)
	}
	ix := func(p int32) int32 { return int32(binary.LittleEndian.Uint32(b[p:])) }
	nn := func(x uint64) int32 { return ix(int32(x) - 12) }
	CK := func(x uint64) []byte { return b[int32(x) : int32(x)+nn(x)] }
	IK := func(x uint64) []int32 { r := make([]int32, nn(x)); br(x, r); return r }
	ux := func(p int32) uint64 { return binary.LittleEndian.Uint64(b[p:]) }
	FK := func(x uint64) []float64 { r := make([]float64, nn(x)); br(x, r); return r }
	ZK := func(x uint64) []complex128 { r := make([]complex128, nn(x)); br(x, r); return r }
	s0 := ix(0)
	SK := func(x uint64) []string {
		v := IK(x)
		r := make([]string, len(v))
		for i, j := range v {
			r[i] = string(CK(ux(s0 + int32(j))))
		}
		return r
	}
	if tp(x) == 25 {
		c := Caption{}
		keys := ux(int32(x))
		vals := ux(int32(x) + 8)
		ks := SK(keys)
		vp := int32(vals)
		for _, s := range ks {
			vi := ux(int32(vp))
			vp += 8
			fm := ""
			var v interface{}
			switch tt := tp(vi); tt {
			case 19: // It
				u := IK(vi)
				w := make([]int, len(u))
				for i := range u {
					w[i] = int(u[i])
				}
				v = w
			case 20: // St
				v = SK(vi)
			case 21: // Ft
				v = FK(vi)
				fm = "%.2f"
			case 22: // Zt
				v = ZK(vi)
				fm = "%.2f@%3.0f"
			default:
				return plts, "caption type is not allowed"
			}
			c.Columns = append(c.Columns, CaptionColumn{
				Class:  s,
				Name:   s,
				Data:   v,
				Format: fm,
			})
		}
		plts[0].Caption = &c
	} else {
		if plts[0].Caption == nil {
			return plts, "plot has no caption"
		}
		c := plts[0].Caption
		colfmt := func(i int, d int) {
			if i >= 0 && i < len(c.Columns) {
				if _, o := c.Columns[i].Data.([]float64); o && d >= 0 {
					c.Columns[i].Format = "%." + strconv.Itoa(d) + "f"
				}
				if _, o := c.Columns[i].Data.([]complex128); o {
					a, b, _ := strings.Cut(c.Columns[i].Format, "@")
					if d >= 0 {
						a = "%." + strconv.Itoa(d) + "f"
					} else {
						b = "%" + strconv.Itoa(4+(-d)) + "." + strconv.Itoa(-d) + "f"
					}
					if a == "" {
						a = "%.2f"
					}
					if b == "" {
						b = "%3.0f"
					}
					c.Columns[i].Format = a + "@" + b
				}
			}
		}
		if tp(x) == 3 {
			for i := range c.Columns {
				colfmt(i, int(int32(x)))
			}
		} else if tp(x) == 18 {
			c.LeadText = strings.Split(string(CK(x)), "\n")
		} else if tp(x) == 19 {
			for i, j := range IK(x) {
				colfmt(i, int(j))
			}
		} else if tp(x) == 20 {
			for i, s := range SK(x) {
				if i < len(c.Columns) {
					c.Columns[i].Unit = s
				}
			}
		} else {
			return plts, "caption: unknown type (should be tiIS)"
		}
	}
	return plts, ""
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
