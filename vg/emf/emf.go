package emf

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// examples (TestRoundRect.emf)
// github.com/LibreOffice/core/tree/master/emfio/qa/cppunit/emf/data

func New(w, h int) *File {
	f := File{
		width:  w,
		height: h,
		Header: Header{
			Type: 1, Size: 108,
			Bounds3: int32(w), Bounds4: int32(h),
			Frame3: int32(20 * w), Frame4: int32(20 * h),
			Signature: 1179469088, Version: 65536,
			Handles:  1, //1st is predefined
			DevWidth: 1920, DevHeight: 1080,
			MilliX: 2 * 192, MilliY: 2 * 108,
			MicroX: 500, MicroY: 500,
		},
	}
	return &f
}

type Color uint32

func (c Color) Value() uint32 { return uint32(c&0xff0000>>16 | c&0x00ff00 | (c&0xff)<<16) }

type Header struct {
	Type, Size                         uint32
	Bounds1, Bounds2, Bounds3, Bounds4 int32
	Frame1, Frame2, Frame3, Frame4     int32
	Signature, Version, Bytes, Records uint32
	Handles, Reserved                  uint16
	NDesc, OffDesc, NPals              uint32
	DevWidth, DevHeight                uint32
	MilliX, MilliY                     uint32
	PixelFormat, OffPixel, Bopengl     uint32 //header extension1
	MicroX, MicroY                     uint32 //header extension2
}
type File struct {
	Header
	Records       []Record
	objects       uint32
	width, height int
}
type Record struct {
	Type uint32
	Size uint32
	Data []uint32
}

const (
	White     Color = 0xffffff
	Black     Color = 0x0
	Grey      Color = 0x808080
	LightGrey Color = 0xc0c0c0
	Red       Color = 0xff0000
	Green     Color = 0x00ff00
	Blue      Color = 0x0000ff
)

type Pen struct {
	Style, Width, BrushStyle uint16
	Color                    Color
}
type Brush struct {
	Style uint32 // 0(solid) 1(null/hollow) wmf:2.1.1.4
	Color Color
	Hatch uint32
}
type Font struct {
	Height       int16
	Escapement   uint16
	Orientation  uint16
	Weight       uint16
	OutPrecision uint16 //ignored
	Quality      uint16 //ignored
	Face         string
}

func (f *File) CreatePen(p Pen) {
	f.push(Record{0x26, 0x1c, []uint32{uint32(f.Handles), uint32(p.Style), uint32(p.Width), uint32(p.BrushStyle), p.Color.Value()}})
	//f.push(Record{0x5f, 56, []uint32{uint32(f.Handles), 56, 0, 56, 0, 65536, p.Width, 0, p.Color.Value(), 0, 0, 0}})
	f.Handles++
}
func (f *File) CreateBrush(b Brush) {
	f.push(Record{0x27, 0x18, []uint32{uint32(f.Handles), b.Style, b.Color.Value(), b.Hatch}})
	f.Handles++
}
func (f *File) Select(i int) { f.push(Record{0x25, 12, []uint32{uint32(1 + i)}}) }

func (f *File) push(x Record) {
	if s := uint32(8 + 4*len(x.Data)); s != x.Size {
		panic(fmt.Sprintf("size mismatch: %d: %v", s, x))
	}
	f.Records = append(f.Records, x)
}
func (f *File) Rectangle(left, top, right, bottom int16) {
	f.push(Record{0x2b, 24, []uint32{uint32(int32(left)), uint32(int32(top)), uint32(int32(right)), uint32(int32(bottom))}})
}
func (f *File) Ellipse(left, top, right, bottom int16) {
	f.push(Record{0x2a, 24, []uint32{uint32(int32(left)), uint32(int32(top)), uint32(int32(right)), uint32(int32(bottom))}})
}
func (f *File) MoveTo(x, y int16) {
	f.push(Record{0x1b, 16, []uint32{uint32(int32(x)), uint32(int32(y))}})
}
func (f *File) LineTo(x, y int16) {
	f.push(Record{0x36, 16, []uint32{uint32(int32(x)), uint32(int32(y))}})
}
func polyvals(x, y []int16) []uint32 {
	var x0, x1, y0, y1 int16
	if len(x) > 0 {
		x0, x1, y0, y1 = x[0], x[0], y[0], y[1]
	}
	for i := range x {
		xi, yi := x[i], y[i]
		if xi < x0 {
			x0 = xi
		}
		if xi > x1 {
			x1 = xi
		}
		if yi < y0 {
			y0 = yi
		}
		if yi > y1 {
			y1 = yi
		}
	}
	r := make([]uint32, 0, 5+len(x))
	r = append(r, []uint32{uint32(int32(x0)), uint32(int32(y0)), uint32(int32(x1)), uint32(int32(y1))}...)
	r = append(r, uint32(len(x)))
	for i := range x {
		r = append(r, uint32(x[i])|uint32(y[i])<<16)
	}
	return r
}
func (f *File) Polygon(x, y []int16)  { f.push(Record{0x56, uint32(4 * (7 + len(x))), polyvals(x, y)}) }
func (f *File) Polyline(x, y []int16) { f.push(Record{0x57, uint32(4 * (7 + len(x))), polyvals(x, y)}) }

func (f *File) Text(x, y int16, s string) { //smalltextout

	ansi := func(s string) (b []byte) { // https://go.dev/play/p/bLORbSxRU63
		b = []byte(s)
		b = bytes.ReplaceAll(b, []byte("ä"), []byte{228})
		b = bytes.ReplaceAll(b, []byte("ö"), []byte{246})
		b = bytes.ReplaceAll(b, []byte("ü"), []byte{252})
		b = bytes.ReplaceAll(b, []byte("Ä"), []byte{196})
		b = bytes.ReplaceAll(b, []byte("Ö"), []byte{214})
		b = bytes.ReplaceAll(b, []byte("Ü"), []byte{220})
		b = bytes.ReplaceAll(b, []byte("ß"), []byte{223})
		b = bytes.ReplaceAll(b, []byte("€"), []byte{128})
		b = bytes.ReplaceAll(b, []byte("°"), []byte{176})
		b = bytes.ReplaceAll(b, []byte("µ"), []byte{181})
		return b
	}
	b := ansi(s)
	n := len(b)
	if m := len(b) % 4; m > 0 {
		for i := 0; i < 4-m; i++ {
			b = append(b, 0)
		}
	}
	t := make([]uint32, len(b)/4)
	binary.Read(bytes.NewReader(b), binary.LittleEndian, t)

	u := []uint32{uint32(x), uint32(y), uint32(n),
		0x200, //options
		2,     //graphics mode
		0, 0, 0, 0, 0xffffffff, 0xffffffff,
	}
	u = append(u, t...)
	f.push(Record{0x6c, 8 + 4*uint32(len(u)), u})
}

/*
func (f *File) Text(x, y int16, s string) { //textouta

		ansi := func(s string) (b []byte) { // https://go.dev/play/p/bLORbSxRU63
			b = []byte(s)
			b = bytes.ReplaceAll(b, []byte("ä"), []byte{228})
			b = bytes.ReplaceAll(b, []byte("ö"), []byte{246})
			b = bytes.ReplaceAll(b, []byte("ü"), []byte{252})
			b = bytes.ReplaceAll(b, []byte("Ä"), []byte{196})
			b = bytes.ReplaceAll(b, []byte("Ö"), []byte{214})
			b = bytes.ReplaceAll(b, []byte("Ü"), []byte{220})
			b = bytes.ReplaceAll(b, []byte("ß"), []byte{223})
			b = bytes.ReplaceAll(b, []byte("€"), []byte{128})
			b = bytes.ReplaceAll(b, []byte("°"), []byte{176})
			b = bytes.ReplaceAll(b, []byte("µ"), []byte{181})
			return b
		}
		b := ansi(s)
		n := len(b)
		if m := len(b) % 4; m > 0 {
			for i := 0; i < 4-m; i++ {
				b = append(b, 0)
			}
		}
		t := make([]uint32, len(b)/4)
		binary.Read(bytes.NewReader(b), binary.LittleEndian, t)
		u := []uint32{0, 0, 0xffffffff, 0xffffffff, 2, 0, 0,
			uint32(x), uint32(y), uint32(n),
			76, 0,
			0, 0, 0xffffffff, 0xffffffff,
			76 + 4*uint32(len(t))}
		u = append(u, t...)
		for i := 0; i < n; i++ {
			u = append(u, 10) //what is correct spacing?
		}
		f.push(Record{0x53, 8 + 4*uint32(len(u)), u})
	}
*/
func (f *File) SetBkMode(mode uint16) { f.push(Record{0x12, 12, []uint32{uint32(mode)}}) }
func (f *File) SetTextColor(c Color)  { f.push(Record{0x18, 12, []uint32{c.Value()}}) }
func (f *File) SetTextAlign(mode uint16) { //mode: refpoint is 0(left top) 2(right) 6(hor-center) 8(bottom) 0x18(baseline) 2.1.2.3  2.3.5.24
	f.push(Record{0x16, 12, []uint32{uint32(mode)}})
}
func (f *File) CreateFont(fn Font) {
	u := []uint32{uint32(f.Handles), uint32(int32(fn.Height)), 0, uint32(fn.Escapement), uint32(fn.Orientation), uint32(fn.Weight), 0, 33554436}
	u = append(u, name64(fn.Face)...)
	u = append(u, make([]uint32, 64)...) //fullname,style,script
	u = append(u, 134248036, 0)          // design vector?
	f.push(Record{0x52, 368, u})
	f.Handles++
}
func name64(s string) []uint32 {
	b := make([]byte, 64)
	for i, c := range s {
		b[2*i] = byte(c)
	}
	u := make([]uint32, 16)
	binary.Read(bytes.NewReader(b), binary.LittleEndian, u)
	return u
}
func (f *File) AntiAlias() { //emf+ header+antialias
	return /* todo..
	f.push(Record{0x46, 0x2c, []uint32{0x20, 726027589,
		0x4001 | 0x1<<16, 0x1c, 0x10, 0xDBC01002, 0x1,
		// 0x66, 0x6c
		0x32, 0x32,
	}})
	f.push(Record{0x46, 28, []uint32{12,
		726027589, //"EMF+"
		0x401e | 0xb<<16, 0xc, 0}})
	*/
}
func (f *File) MarshallBinary() []byte {
	var r []Record
	for _, x := range f.Records {
		r = append(r, x)
	}
	r = append(r, Record{0x0e, 20, []uint32{0, 16, 20}}) // eof
	f.Header.Records = uint32(len(r))

	total := uint32(binary.Size(f.Header))
	for _, r := range r {
		total += r.Size
	}
	f.Bytes = total

	var b bytes.Buffer
	fatal(binary.Write(&b, binary.LittleEndian, f.Header))
	for _, x := range r {
		fatal(binary.Write(&b, binary.LittleEndian, x.Type))
		fatal(binary.Write(&b, binary.LittleEndian, x.Size))
		for _, d := range x.Data {
			fatal(binary.Write(&b, binary.LittleEndian, d))
		}
	}
	return b.Bytes()
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
