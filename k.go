package plot

import "fmt"

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
