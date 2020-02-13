package plot

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/ktye/plot/xmath"
)

func (plts Plots) Encode(w io.Writer) error {
	for _, p := range plts {
		if e := p.Encode(w); e != nil {
			return e
		}
	}
	return nil
}
func (p Plot) Encode(w io.Writer) error {
	if p.Type == Foto || p.Type == Raster {
		return fmt.Errorf("exporting a raster image is not supported")
	}
	fmt.Fprintf(w, "Plot\n")
	e := js(w, " Type", p.Type, nil)
	e = js(w, " Style", p.Style, e)
	e = js(w, " Limits", p.Limits, e)
	e = js(w, " Xlabel", p.Xlabel, e)
	e = js(w, " Ylabel", p.Ylabel, e)
	e = js(w, " Xunit", p.Xunit, e)
	e = js(w, " Yunit", p.Yunit, e)
	e = js(w, " Zunit", p.Zunit, e)
	if e != nil {
		return e
	}
	for _, l := range p.Lines {
		if e := l.Encode(w); e != nil {
			return e
		}
	}
	if c := p.Caption; c != nil {
		if e := c.Encode(w); e != nil {
			return e
		}
	}
	// p.Data is ignored
	return e
}
func (l Line) Encode(w io.Writer) error {
	fmt.Fprintf(w, " Line\n")
	e := js(w, "  X", l.X, nil)
	e = js(w, "  Y", l.Y, e)
	e = js(w, "  R", xmath.RealVector(l.C), e)
	e = js(w, "  I", xmath.ImagVector(l.C), e)
	e = js(w, "  V", l.V, e)
	e = js(w, "  Segments", l.Segments, e)
	e = js(w, "  Style", l.Style, e)
	e = js(w, "  Id", l.Id, e)
	return e
}
func (c Caption) Encode(w io.Writer) error {
	fmt.Fprintf(w, " Caption\n")
	e := js(w, "  Title", c.Title, nil)
	e = js(w, "  LeadText", c.LeadText, e)
	if e != nil {
		return e
	}
	if c.Columns == nil {
		return e
	}
	for _, col := range c.Columns {
		if e := col.Encode(w); e != nil {
			return e
		}
	}
	return nil
}
func (c CaptionColumn) Encode(w io.Writer) error {
	fmt.Fprintf(w, "  CaptionColumn\n")
	e := js(w, "   Class", c.Class, nil)
	e = js(w, "   Name", c.Name, e)
	e = js(w, "   Unit", c.Unit, e)
	e = js(w, "   Format", c.Format, e)
	if e != nil {
		return e
	}
	switch d := c.Data.(type) {
	case []string:
		e = js(w, "   StringData", d, e)
	case []int:
		e = js(w, "   IntData", d, e)
	case []float64:
		e = js(w, "   FloatData", d, e)
	case []complex128:
		e = js(w, "   ReData", xmath.RealVector(d), e)
		e = js(w, "   ReData", xmath.ImagVector(d), e)
	default:
		return fmt.Errorf("illegal caption column type: %T", c.Data)
	}
	return e
}

func js(w io.Writer, name string, v interface{}, e error) error {
	if e != nil {
		return e
	}
	var b []byte
	b, e = json.Marshal(v)
	if e != nil {
		return e
	}
	fmt.Fprintf(w, "%s %s\n", name, string(b))
	return nil
}
