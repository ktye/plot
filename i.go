package plot

// The following methods integrate the plot package into ktye/i/⍳.
// The Stringers for Plot and Plots return an url-encoded image.
// ConvertTo is used by dyadic $ (cst) to create a plot from a vector or list of vectors:
//	p: plot$?-100
// plot is a variable of type Plot, present in the k-tree by default.

// StringSize is the image size for the default stringers.
var StringSize struct {
	Width, Height int
}

// String prevents the implementation of a Stringer from the embedded Limits.
// This is for i to show a plot as a dict and not use a stringer.
func (p Plot) String() {}

/*
// String url-encodes a plot as a png.
// It uses StringSize as size.
func (p Plot) String() string {
	return Plots{p}.String()
}

func (p Plots) String() string {
	w := StringSize.Width
	if w == 0 {
		w = 800
	}
	h := StringSize.Height
	if h == 0 {
		h = 600
	}
	ip, err := p.IPlots(w, h)
	if err != nil {
		return err.Error()
	}
	m := Image(ip, nil, w, h)
	s, err := EncodeToPng(m)
	if err != nil {
		return err.Error()
	}
	return s
}
*/

// ConvertTo returns a single line plot from []complex128.
// or multiple lines from []interface{[]complex128{}}.
// If all imag parts are zero, the values are converted to []float64.
// The plot type is set to "xy" if all values are real, otherwise "polar".
func (p Plot) ConvertTo(in interface{}) interface{} {
	e := func(s string) Plot {
		return Plot{Title: s}
	}
	realvec := func(z []complex128) []float64 {
		r := make([]float64, len(z))
		for i := range z {
			r[i] = real(z[i])
		}
		return r
	}
	xvec := func(n int) []float64 {
		r := make([]float64, n)
		for i := range r {
			r[i] = float64(i)
		}
		return r
	}
	if z, ok := in.([]complex128); ok {
		in = []interface{}{z}
	}
	l, ok := in.([]interface{})
	if !ok {
		return e("input error")
	}
	isz := false
	for i := range l {
		zv, ok := l[i].([]complex128)
		if ok == false {
			return e("input error")
		}
		for _, u := range zv {
			if imag(u) != 0 {
				isz = true
				goto done
			}
		}
	}
done:
	p.Type = XY
	if isz {
		p.Type = Polar
	}
	p.Lines = make([]Line, len(l))
	for i := range l {
		z := l[i].([]complex128)
		p.Lines[i].X = xvec(len(z))
		if isz {
			p.Lines[i].C = z
		} else {
			p.Lines[i].Y = realvec(z)
		}
	}
	return p
}
