package plot

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"strconv"

	"github.com/ktye/plot/xmath"
)

// TODO: use encoding.TextMarshaler/TextUnmarshaler instead? does than help?

func (plts Plots) Encode(w io.Writer) error {
	for _, p := range plts {
		if e := p.encode(w); e != nil {
			return e
		}
	}
	return nil
}
func DecodePlots(r io.Reader) (Plots, error) {
	b, e := ioutil.ReadAll(r)
	if e != nil {
		return nil, e
	}
	lr := &lineReader{Buffer: bytes.NewBuffer(b)}
	var plts Plots
	for {
		p, e := decodePlot(lr)
		if e == nil {
			plts = append(plts, p)
		} else if e == io.EOF {
			return append(plts, p), nil
		} else {
			return nil, e
		}
	}
}
func DecodeAny(r io.Reader) (Plots, error) {
	br := bufio.NewReader(r)
	c, err := br.Peek(1)
	if err != nil {
		return nil, err
	}
	if c[0] == 'P' {
		return DecodePlots(br)
	} else {
		return TextDataPlot(br)
	}
}
func (p Plot) encode(w io.Writer) error {
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
		if e := l.encode(w); e != nil {
			return e
		}
	}
	if c := p.Caption; c != nil {
		if e := c.encode(w); e != nil {
			return e
		}
	}
	// p.Data is ignored
	return e
}
func decodePlot(r LineReader) (p Plot, e error) {
	e = expect(r, "Plot", e)
	e = sj(r, " Type", &p.Type, e)
	e = sj(r, " Style", &p.Style, e)
	e = sj(r, " Limits", &p.Limits, e)
	e = sj(r, " Xlabel", &p.Xlabel, e)
	e = sj(r, " Ylabel", &p.Ylabel, e)
	e = sj(r, " Xunit", &p.Xunit, e)
	e = sj(r, " Yunit", &p.Yunit, e)
	e = sj(r, " Zunit", &p.Zunit, e)
	if e != nil {
		return p, e
	}
	for {
		if l, e := decodeLine(r); e == nil {
			p.Lines = append(p.Lines, l)
		} else if e.Error() == errNextCaption {
			break
		} else if e.Error() == errNextPlot {
			return p, nil
		} else if e != nil {
			return p, e // maybe io.EOF
		}
	}
	var c Caption
	c, e = decodeCaption(r)
	if e == nil || e == io.EOF {
		p.Caption = &c
	}
	return p, e
}
func (l Line) encode(w io.Writer) error {
	fmt.Fprintf(w, " Line\n")
	e := fs(w, "  X", l.X, nil)
	e = fs(w, "  Y", l.Y, e)
	e = fs(w, "  R", xmath.RealVector(l.C), e)
	e = fs(w, "  I", xmath.ImagVector(l.C), e)
	e = fs(w, "  V", l.V, e)
	e = js(w, "  Segments", l.Segments, e)
	e = js(w, "  Style", l.Style, e)
	e = js(w, "  Id", l.Id, e)
	return e
}
func decodeLine(r LineReader) (l Line, e error) {
	if b, e := r.Peek(); e != nil {
		return l, e
	} else if bytes.HasPrefix(b, []byte(" Caption")) {
		return l, errors.New(errNextCaption)
	} else if bytes.HasPrefix(b, []byte("Plot")) {
		return l, errors.New(errNextPlot)
	} else if bytes.HasPrefix(b, []byte(" Line")) == false {
		return l, fmt.Errorf("line %d: decode Line", r.LineNumber())
	}
	var re, im []float64
	e = expect(r, " Line", e)
	e = sf(r, "  X", &l.X, e)
	e = sf(r, "  Y", &l.Y, e)
	e = sf(r, "  R", &re, e)
	e = sf(r, "  I", &im, e)
	e = sf(r, "  V", &l.V, e)
	e = sj(r, "  Segments", &l.Segments, e)
	e = sj(r, "  Style", &l.Style, e)
	e = sj(r, "  Id", &l.Id, e)
	if e == nil {
		l.C = xmath.ComplexVector(re, im)
	}
	return l, e
}
func (c Caption) encode(w io.Writer) error {
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
		if e := col.encode(w); e != nil {
			return e
		}
	}
	return nil
}
func decodeCaption(r LineReader) (c Caption, e error) {
	e = expect(r, " Caption", e)
	e = sj(r, "  Title", &c.Title, e)
	e = sj(r, "  LeadText", &c.LeadText, e)
	if e != nil {
		return c, e
	}
	for {
		if col, e := decodeCaptionColumn(r); e == nil {
			c.Columns = append(c.Columns, col)
		} else if e.Error() == errNextPlot {
			return c, nil
		} else if e != nil {
			return c, e // maybe io.EOF
		}
	}
}
func (c CaptionColumn) encode(w io.Writer) error {
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
		e = fs(w, "   FloatData", d, e)
	case []complex128:
		e = fs(w, "   ReData", xmath.RealVector(d), e)
		e = fs(w, "   ImData", xmath.ImagVector(d), e)
	default:
		return fmt.Errorf("illegal caption column type: %T", c.Data)
	}
	return e
}
func decodeCaptionColumn(r LineReader) (c CaptionColumn, e error) {
	if b, e := r.Peek(); e != nil {
		return c, e
	} else if bytes.HasPrefix(b, []byte("Plot")) {
		return c, errors.New(errNextPlot)
	} else if bytes.HasPrefix(b, []byte(" CaptionColumn")) == false {
		return c, fmt.Errorf("line %d: decode CaptionColumn", r.LineNumber())
	}
	e = expect(r, "  CaptionColumn", e)
	e = sj(r, "  Class", &c.Class, e)
	e = sj(r, "  Name", &c.Name, e)
	e = sj(r, "  Unit", &c.Unit, e)
	e = sj(r, "  Format", &c.Format, e)
	if e != nil {
		return c, e
	}
	if b, e := r.Peek(); e != nil {
		return c, e
	} else if bytes.HasPrefix(b, []byte("   StringData")) {
		var d []string
		sj(r, "   StringData", &d, e)
		c.Data = d
	} else if bytes.HasPrefix(b, []byte("   IntData")) {
		var d []int
		sj(r, "   StringData", &d, e)
		c.Data = d
	} else if bytes.HasPrefix(b, []byte("   FloatData")) {
		var d []float64
		sf(r, "   FloatData", &d, e)
		c.Data = d
	} else if bytes.HasPrefix(b, []byte("   ReData")) {
		var re, im []float64
		sf(r, "   ReData", &re, e)
		sf(r, "   ImData", &im, e)
		c.Data = xmath.ComplexVector(re, im)
	} else {
		return c, fmt.Errorf("%d: caption column without data", r.LineNumber())
	}
	return c, e
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
func fs(w io.Writer, name string, v []float64, e error) error { // js cannot handle nan
	if e != nil {
		return e
	}
	w.Write([]byte(name))
	for _, f := range v {
		fmt.Fprintf(w, " %v", f)
	}
	w.Write([]byte{'\n'})
	return nil
}
func sj(r LineReader, name string, v interface{}, e error) error {
	if e != nil {
		return e
	}
	b, e := r.ReadLine()
	if e != nil {
		return e
	}
	if bytes.HasPrefix(b, []byte(name)) == false {
		return fmt.Errorf("line %d: expected %q", r.LineNumber(), name)
	} else {
		b = b[len(name):]
	}
	e = json.Unmarshal(b, v)
	if e != nil {
		return fmt.Errorf("line %d: %s", r.LineNumber(), e)
	}
	return nil
}
func sf(r LineReader, name string, v *[]float64, e error) error {
	if e != nil {
		return e
	}
	b, e := r.ReadLine()
	if e != nil {
		return e
	}
	if bytes.HasPrefix(b, []byte(name)) == false {
		return fmt.Errorf("line %d: expected %q", r.LineNumber(), name)
	} else {
		b = b[len(name):]
	}
	f := bytes.Fields(b)
	z := make([]float64, len(f))
	for i, b := range f {
		if len(b) == 2 && b[0] == '0' && b[1] == 'n' {
			z[i] = math.NaN()
		} else {
			z[i], e = strconv.ParseFloat(string(b), 64)
			if e != nil {
				return fmt.Errorf("line %d: %s", r.LineNumber(), e)
			}
		}
	}
	*v = z
	return nil
}
func expect(r LineReader, s string, e error) error {
	if e != nil {
		return e
	}
	if b, e := r.ReadLine(); e != nil {
		return e
	} else if string(b) != s {
		return fmt.Errorf("line %d: expected: %q got %q", r.LineNumber(), s, string(b))
	}
	return nil
}

type LineReader interface {
	ReadLine() ([]byte, error)
	Peek() ([]byte, error)
	LineNumber() int
}
type lineReader struct {
	*bytes.Buffer
	peek []byte
	lino int
}

func (l *lineReader) ReadLine() (c []byte, e error) {
	if l.peek != nil {
		c = l.peek
		l.peek = nil
		return c, nil
	}
	l.lino++
	c, e = l.Buffer.ReadBytes('\n')
	if e == io.EOF {
		return c, io.EOF
	} else if e != nil {
		return c, fmt.Errorf("line %d: %s", l.lino, e)
	}
	return bytes.Trim(c, "\r\n"), e
}
func (l *lineReader) Peek() (c []byte, e error) {
	c, e = l.ReadLine()
	l.peek = c
	return c, e
}
func (l *lineReader) LineNumber() int {
	return l.lino
}

const errNextCaption = "caption follows"
const errNextPlot = "plot follows"
