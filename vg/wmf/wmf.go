package wmf

// [MS-WMF]-220625.pdf
// https://wvware.sourceforge.net/caolan/ora-wmf.html
// https://github.com/g21589/wmf2canvas/blob/master/js/wmf.js

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func New(w, h int) *File {
	f := File{
		width:  w,
		height: h,
		Header: Header{
			//Key:     2596720087,
			//Right:   int16(w),
			//Bottom:  int16(h),
			//Inch:    1440,
			Type:    1, //1(in-mem) 2(on-disk)
			HdrSize: 9,
			Version: 0x0300,
		},
	}
	return &f
}

type File struct {
	Header
	Records       []Record
	width, height int
}
type Header struct {
	/*
		//meta placeable record
		Key      uint32
		Hwmf     uint16
		Left     int16
		Top      int16
		Right    int16
		Bottom   int16
		Inch     int16
		Reserved uint32
		Checksum uint16
	*/

	//header record
	Type       uint16
	HdrSize    uint16
	Version    uint16
	TotalSize  uint32
	NumObjects uint16
	MaxRecord  uint32
	Unused     uint16
}
type Record struct {
	Size uint32
	Cmd  uint16
	Data []uint16
}

type Pen struct {
	Style   uint16 // 0(solid) 1(dash) 2(dot) 3(dashdot) 4(dashdotdot) .. 0x0100(endsquare) 0x0200(endflat) 0x1000(joinbevel) 0x2000(joinmiter) 2.1.1.23
	Width   uint16
	Ignored uint16
	Color   Color
}
type Brush struct {
	Style uint16 // 0(solid) 1(null/hollow) 2.1.1.4
	Color Color
	Hatch uint16 //2.1.1.12
}
type Color uint32

const (
	White     Color = 0xffffff
	Black     Color = 0x0
	Grey      Color = 0x808080
	LightGrey Color = 0xc0c0c0
	Red       Color = 0xff0000
	Green     Color = 0x00ff00
	Blue      Color = 0x0000ff
)

func (c Color) Values() (uint16, uint16) { return uint16(c>>16) | uint16(c)&0xff00, uint16(c & 0xff) }

type Font struct {
	Height                                               int16 //<0: transform into device units
	Width                                                int16 //0: keep aspect ratio
	Escapement, Orientation                              uint16
	Weight                                               uint16 //0 or 400: normal, 700 bold
	Italic, Underline, Strikeout                         uint8
	Charset                                              uint8 //0(ansi) 1(default) 2(symbol) ..
	OutPrecision, ClipPrecision, Quality, PitchAndFamily uint8
	Face                                                 string
}

func (f Font) Value() []uint16 {
	b := make([]byte, 32)
	copy(b, []byte(f.Face))
	name := make([]uint16, 16)
	for i := range name {
		name[i] = binary.LittleEndian.Uint16(b[2*i:])
	}
	return append([]uint16{uint16(f.Height), uint16(f.Width), f.Escapement, f.Orientation, f.Weight, uint16(f.Italic) | uint16(f.Underline)<<8,
		uint16(f.Strikeout) | uint16(f.Charset)<<8, uint16(f.OutPrecision) | uint16(f.ClipPrecision)<<8,
		uint16(f.Quality) | uint16(f.PitchAndFamily)<<8}, name...)
}

func (r Record) size() int { return 6 + 2*len(r.Data) }

func (f *File) MoveTo(x, y int) {
	f.push(Record{5, 0x0214, []uint16{uint16(int16(y)), uint16(int16(x))}})
}
func (f *File) LineTo(x, y int) {
	f.push(Record{5, 0x0213, []uint16{uint16(int16(y)), uint16(int16(x))}})
}
func (f *File) Rectangle(left, top, right, bottom int16) {
	f.push(Record{7, 0x041b, []uint16{uint16(bottom), uint16(right), uint16(top), uint16(left)}})
}
func (f *File) Ellipse(left, top, right, bottom int16) {
	f.push(Record{7, 0x0418, []uint16{uint16(bottom), uint16(right), uint16(top), uint16(left)}})
}
func (f *File) Polyline(x, y []int16) {
	u := make([]uint16, 1+2*len(x))
	u[0] = uint16(len(x))
	for i := range x {
		u[2*i+1] = uint16(x[i])
		u[2*i+2] = uint16(y[i])
	}
	f.push(Record{3 + uint32(len(u)), 0x0325, u})
}
func (f *File) Polygon(x, y []int16) {
	u := make([]uint16, 1+2*len(x))
	u[0] = uint16(len(x))
	for i := range x {
		u[2*i+1] = uint16(x[i])
		u[2*i+2] = uint16(y[i])
	}
	f.push(Record{3 + uint32(len(u)), 0x0324, u})
}
func (f *File) Text(x, y int16, s string) { //textout
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
	if len(b)%2 != 0 {
		b = append(b, 0)
	}
	u := make([]uint16, 3+len(b)/2)
	u[0] = uint16(n)
	j := 0
	for i := 1; i < len(u)-2; i++ {
		u[i] = uint16(b[j]) | uint16(b[1+j])<<8
		j += 2
	}
	u[len(u)-2], u[len(u)-1] = uint16(y), uint16(x)
	f.push(Record{3 + uint32(len(u)), 0x0521, u})
}
func (f *File) CreateFont(fn Font) { f.push(Record{28, 0x02fb, fn.Value()}); f.NumObjects++ }
func (f *File) SetTextAlign(mode uint16) { //mode: refpoint is 0(left top) 2(right) 6(hor-center) 8(bottom) 0x18(baseline) 2.1.2.3  2.3.5.24
	f.push(Record{5, 0x012e, []uint16{mode, 0}})
}
func (f *File) SetTextColor(b Color) { //0x0209
	a, c := b.Values()
	f.push(Record{5, 0x0209, []uint16{a, c}})
}
func (f *File) SetBkColor(b Color) {
	a, c := b.Values()
	f.push(Record{5, 0x0201, []uint16{a, c}})
}
func (f *File) SetBkMode(mode uint16) { //mode: transparent(0x0001) opaque(0x0002) 2.1.1.20
	f.push(Record{5, 0x0102, []uint16{mode, 0}})
}
func (f *File) CreatePen(p Pen) { //createpenindirect
	a, b := p.Color.Values()
	f.push(Record{8, 0x02fa, []uint16{p.Style, p.Width, 0, a, b}})
	f.NumObjects++
}
func (f *File) CreateBrush(b Brush) { //createbrushindirect
	a, c := b.Color.Values()
	f.push(Record{7, 0x02fc, []uint16{b.Style, a, c, b.Hatch}})
	f.NumObjects++
}
func (f *File) SetMapMode(mode uint16) { //mode: 1(text:pixels) 2(1u=.1mm) 3(1u=.01mm) 4(1u=.01in) 5(1u=.001in) 6(twip 1/20th of a point which is 1/1440 of an inch) 2.1.1.16
	f.push(Record{4, 0x0103, []uint16{mode}})
}

// func (f *File) SetWindowExt(w, h int16) { f.push(Record{5, 0x020c, []uint16{uint16(h), uint16(w)}}) }
func (f *File) Select(i int) { //selectobject
	if i < 0 {
		i = int(f.NumObjects - 1)
	}
	f.push(Record{4, 0x012d, []uint16{uint16(i)}})
}
func (f *File) SelectClip(region uint16) { //selectclipregion
	f.push(Record{4, 0x012c, []uint16{region}})
}
func (f *File) ExcludeClipRect(left, top, right, bottom int16) { //draw outside rect
	f.push(Record{7, 0x0415, []uint16{uint16(bottom), uint16(right), uint16(top), uint16(left)}})
}
func (f *File) IntersectClipRect(left, top, right, bottom int16) { //draw inside rect
	f.push(Record{7, 0x0416, []uint16{uint16(bottom), uint16(right), uint16(top), uint16(left)}})
}
func (f *File) push(x Record) {
	if s := uint32(3 + len(x.Data)); s != x.Size {
		panic(fmt.Sprintf("size mismatch: %d: %v", s, x))
	}
	f.Records = append(f.Records, x)
}

func (f *File) MarshallBinary() []byte {

	r := make([]Record, 3+len(f.Records))
	r[0] = Record{5, 0x020c, []uint16{uint16(f.height), uint16(f.width)}} //SetWindowExt
	r[1] = Record{4, 0x0103, []uint16{1}}                                 //setmapmode 1(pixel) 2(1unit=0.1mm) 4(1u=0.01in) 7(isotropic)
	for i, x := range f.Records {
		r[2+i] = x
	}
	r[len(r)-1] = Record{3, 0, nil} //eof

	m, total := 0, 0
	for _, r := range r {
		n := r.size()
		total += n
		if n > m {
			m = n
		}
	}
	f.MaxRecord = uint32(m / 2)
	f.TotalSize = uint32(9 + total/2) //meta-placeable header does not count

	var b bytes.Buffer
	//b.Write(f.placeableHeader())
	fatal(binary.Write(&b, binary.LittleEndian, f.Header))
	for _, x := range r {
		fatal(binary.Write(&b, binary.LittleEndian, x.Size))
		fatal(binary.Write(&b, binary.LittleEndian, x.Cmd))
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
func (f *File) placeableHeader() []byte {
	checksum := func(b []byte) uint16 {
		var v [10]uint16
		binary.Read(bytes.NewReader(b), binary.LittleEndian, &v)
		var s uint16
		for _, x := range v {
			s ^= x
		}
		return s
	}

	var b bytes.Buffer
	bw := func(x interface{}) { binary.Write(&b, binary.LittleEndian, x) }
	b.Write([]byte{0xd7, 0xcd, 0xc6, 0x9a}) //key
	b.Write([]byte{0, 0})                   //hwmf (on-disk)
	bw(uint16(0))                           //bbox left
	bw(uint16(0))                           //bbox right
	bw(uint16(f.width))                     //bbox right
	bw(uint16(f.height))                    //bbox bottom
	bw(uint16(96))                          //units per inch
	b.Write([]byte{0, 0, 0, 0})             //reserved
	bw(checksum(b.Bytes()))                 //checksum (xor of 10 previous 16bit words)
	return b.Bytes()
}

/*
func (w *wf) Bytes() []byte {
	var b bytes.Buffer
	bw := func(x interface{}) { binary.Write(&b, binary.LittleEndian, x) }

	// meta placeable record (2.3.2.3) 22-bytes
	b.Write([]byte{0xd7, 0xcd, 0xc6, 0x9a}) //key
	b.Write([]byte{0, 0})                   //hwmf (on-disk)
	bw(uint16(0))                           //bbox left
	bw(uint16(0))                           //bbox right
	bw(w.width)                             //bbox right
	bw(w.height)                            //bbox bottom
	bw(w.inch)                              //units per inch
	b.Write([]byte{0, 0, 0, 0})             //reserved
	cs := checksum(b.Bytes())               //
	bw(cs)                                  //checksum (xor of 10 previous 16bit words)

	total, max := 0, 0
	for _, o := range w.obj {
		n := len(o)
		total += n
		if n > max {
			max = n
		}
	}

	// header record (2.3.2.2)
	bw(uint16(2))       //type 1(inmemory) 2(ondisk)
	bw(uint16(9))       //header size in words
	bw(uint16(0x0300))  //version
	bw(uint32(0))       //lo/hi size in 16bit words: (total size in bytes - 22)/2 TODO
	bw(w.nobj)          //number of graphics objects (brushes, fonts, palettes, pens, regions) (defined prior to the records)
	bw(uint32(max / 2)) //maxrecord size of largest record in 16bit words         TODO
	bw(uint16(0))       //unused

	for _, o := range w.obj {
		b.Write(o)
	}
	bw(uint32(3)) // eof record size
	bw(uint16(0)) // eof
	return b.Bytes()
}

//  4 0x0301 setmapmode
//  5 0x0b02 setwindoworg
//  5 0x0c02 seteindowext
//  4 0x0601 setpolyfillmode
//  7 0xfc02 createbrushindirect
//  4 0x2d01 selectobj
//  8 0xfa02 createpenindirect
// 99 0x3805 polypolygon

type wf struct {
	width, height int16    //logical units
	inch          int16    //logical units per inch (convention 1440)
	nobj          uint16   // number of graphics objects brushes, etc..)
	obj           [][]byte //records
}
*/
