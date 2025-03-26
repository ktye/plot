package emfplus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"unicode/utf16"
)

func New(w, h int) *File {
	wf, hf := 25.4*float64(w)/96, 25.4*float64(h)/96
	f := File{
		Header: Header{Type: 1, Size: 108, Bounds1: 0, Bounds2: 0, Bounds3: int32(w), Bounds4: int32(h),
			Frame1: 0, Frame2: 0, Frame3: int32(wf * 100), Frame4: int32(hf * 100), Signature: 1179469088, Version: 65536,
			Bytes: 1228, Records: 34, Handles: 2, Reserved: 0,
			NDesc: 0, OffDesc: 0, NPals: 0,
			DevWidth: 1920, DevHeight: 960, MilliX: 508, MilliY: 254,
			PixelFormat: 0, OffPixel: 0, Bopengl: 0,
			MicroX: 508000, MicroY: 254000},
		width:   w,
		height:  h,
		objects: 1,
	}
	return &f
}

type File struct {
	Header
	Records       []Record
	width, height int
	objects       uint8
	pens          map[pen]uint8
	fontmap       map[string]uint8
	alignments    map[string]uint8
}
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
type Record struct {
	Type  uint16
	Flags uint16
	Size  uint32
	Data  []uint32
}
type Emf struct {
	Type, Size uint32
	Data       []uint32
}

type pen struct {
	linewidth int16
	color     uint32
}

func (f *File) Pen(linewidth int16, color uint32) uint8 { //color: 0xaarrggbb
	p := pen{linewidth, color}
	if f.pens == nil {
		f.pens = make(map[pen]uint8)
	}
	if id, ok := f.pens[p]; ok {
		return id
	}

	id := f.objects
	f.pens[p] = id
	lw := f32(linewidth)
	f.push(Record{0x4008, 0x200 | uint16(id), 52 - 8, []uint32{40 - 8, 3686797314, 0, 0, 0, lw, 3686797314, 0, color /*4278219453*/}})
	f.objects++
	return id
}
func rect(x, y, w, h int16) (u1, u2 uint32) {
	return uint32(x) | uint32(y)<<16, uint32(w) | uint32(h)<<16
}
func (f *File) DrawEllipse(pen uint8, x, y, w, h int16) {
	u1, u2 := rect(x, y, w, h)
	f.push(Record{0x400f, 0x4000 | uint16(pen), 0x14, []uint32{0x8, u1, u2}})
}
func (f *File) FillEllipse(color uint32, x, y, w, h int16) {
	u1, u2 := rect(x, y, w, h)
	f.push(Record{0x400e, 0xc000, 24, []uint32{12, color, u1, u2}})
}
func (f *File) DrawRects(pen uint8, x, y, w, h []int16) {
	count := uint32(len(x))
	n := 4 + 8*count
	u := []uint32{n, count}
	for i := range x {
		u1, u2 := rect(x[i], y[i], w[i], h[i])
		u = append(u, u1, u2)
	}
	f.push(Record{0x400b, 0x4000 | uint16(pen), 12 + n, u})
}
func (f *File) FillRects(color uint32, x, y, w, h []int16) {
	count := uint32(len(x))
	n := 8 + 8*count
	u := []uint32{n, color, count}
	for i := range x {
		u1, u2 := rect(x[i], y[i], w[i], h[i])
		u = append(u, u1, u2)
	}
	f.push(Record{0x400a, 0xc000, 12 + n, u})
}
func (f *File) DrawPolyline(pen uint8, closepath bool, x, y []int16) {
	cp := uint16(0)
	if closepath {
		cp = uint16(0x6000)
	}
	count := uint32(len(x))
	u := make([]uint32, count)
	for i := range x {
		u[i] = uint32(x[i]) | uint32(y[i])<<16
	}
	f.push(Record{0x400d, 0xc000 | cp | uint16(pen), 16 + 4*count, append([]uint32{4 + 4*count, count}, u...)})
}
func (f *File) FillPolygon(color uint32, x, y []int16) {
	count := uint32(len(x))
	u := make([]uint32, count)
	for i := range x {
		u[i] = uint32(x[i]) | uint32(y[i])<<16
	}
	f.push(Record{0x400c, 0xc000, 0x14 + 4*count, append([]uint32{8 + 4*count, color, count}, u...)})
}
func (f *File) LineSegments(pen uint8, x0, x1, y0, y1 []int16) { //multiple individual lines using a path.
	pathid := f.objects
	//f.objects++ (dont store path)
	count := uint32(2 * len(x0))
	u := make([]uint32, 2*len(x0))
	t := make([]byte, 2*len(x0))
	for i := range x0 {
		u[2*i+0] = uint32(x0[i]) | uint32(y0[i])<<16
		u[2*i+1] = uint32(x1[i]) | uint32(y1[i])<<16
		t[2*i+0] = 0
		t[2*i+1] = 9
	}
	for len(t)%4 != 0 {
		t = append(t, 0xbf)
	}
	v := make([]uint32, len(t)/4)
	binary.Read(bytes.NewReader(t), binary.LittleEndian, v)

	s := uint32(0)
	data := []uint32{s, 3686797314, count, 24576}
	data = append(data, u...)
	data = append(data, v...)
	s = uint32(4*len(data) - 4)
	data[0] = s

	f.push(Record{0x4008, 0x300 | uint16(pathid), 12 + s, data})
	f.push(Record{0x4015, uint16(pathid), 16, []uint32{4, uint32(pen)}})
	// pathpointtype [flags|type] 4bits|4bits --> 0 9 0 9 0 9
	// flag: close subpath 0x08
	// type: 0x00(start) 0x01(line)
}
func (f *File) DrawImage(id uint8, x, y, w, h int16) {
	u1, u2 := rect(x, y, w, h)
	f.push(Record{0x4030, 2, 0x10, []uint32{0x4, math.Float32bits(1.0)}})
	f.push(Record{0x401a, 0x4000 | uint16(id), 0x2c, []uint32{0x20, 0, 0 /*2*/, f32(0), f32(0), f32(w), f32(h), u1, u2}})
	f.push(Record{0x4030, 2, 0x10, []uint32{0x4, math.Float32bits(1.0 / 8.0)}})
}
func (f *File) Clip(x, y, w, h int16) {
	f.push(Record{0x4032, 0, 0x1c, []uint32{0x10, f32(x), f32(y), f32(w), f32(h)}}) //replace current clip region
}
func (f *File) Png(w, h int, png []byte) uint8 {
	for len(png)%4 != 0 {
		png = append(png, 0)
	}
	v := make([]uint32, len(png)/4)
	binary.Read(bytes.NewReader(png), binary.LittleEndian, &v)
	u := append([]uint32{0, 3686797314, 1, uint32(w), uint32(h), 0, 0, 1}, v...)
	n := uint32(4*len(u) - 4)
	u[0] = n

	id := f.objects
	f.objects++
	f.push(Record{0x4008, 0x0500 | uint16(id), 12 + n, u})
	return id
}
func uni(s string) ([]uint32, uint32) {
	u := utf16.Encode([]rune(s))
	count := len(u)
	if len(u)%2 != 0 {
		u = append(u, 0)
	}
	r := make([]uint32, len(u)/2)
	for i := range r {
		r[i] = uint32(u[2*i]) | uint32(u[1+2*i])<<16
	}
	return r, uint32(count)
}
func (f *File) Font(size int16, name string) uint8 {
	id := fmt.Sprintf("name@%d", size)
	if f.fontmap == nil {
		f.fontmap = make(map[string]uint8)
	}
	if r, o := f.fontmap[id]; o {
		return r
	}
	fontid := f.objects
	f.objects++
	f.fontmap[id] = fontid

	family, count := uni(name)
	em := f32(size)
	u := append([]uint32{0, 3686797314, em, 2, 0, 0, count}, family...)
	n := uint32(4*len(u) - 4)
	u[0] = n
	f.push(Record{0x4008, 0x600 | uint16(fontid), 12 + n, u})
	return fontid
}
func (f *File) align(a int) uint8 {
	if f.alignments == nil {
		f.alignments = make(map[string]uint8)
	}
	aid := fmt.Sprintf("%d", a)
	if id, o := f.alignments[aid]; o {
		return id
	}

	id := f.objects
	f.alignments[aid] = id
	f.objects++
	s := []uint32{0, 1, 2, 2, 2, 1, 0, 0, 1}[a]
	l := []uint32{2, 2, 2, 1, 0, 0, 0, 1, 1}[a]
	u := []uint32{
		60, 3686797314,
		26628,    //flags: 0x6804 (dont clip)
		27459584, //lang(?)
		s, l,     //string,line align: 0(near) 1(center) 2(far)
		1,         //no digit substitution
		136052736, //digit language
		0, 0, 0, 0,
		1065353216, //tracking 1.0
		0, 0, 0}

	f.push(Record{0x4008, 0x700 | uint16(id), 72, u})
	return id
}
func (f *File) Text(x, y int16, t string, fn uint8, align int, vertical bool, color uint32) {
	al := uint32(f.align(align))
	if vertical {
		f.push(Record{0x402d, 0, 0x14, []uint32{0x8, f32(x), f32(y)}}) //translate
		f.push(Record{0x402f, 0, 0x10, []uint32{0x4, f32(-90)}})       //rotate
		x, y = 0, 0
	}
	s, count := uni(t)
	u := append([]uint32{0, color, al, count, f32(x), f32(y), 0, 0}, s...)
	n := uint32(4*len(u) - 4)
	u[0] = n
	f.push(Record{0x401c, 0x8000 | uint16(fn), 12 + n, u})
	if vertical {
		f.push(Record{0x402b, 0, 0xc, []uint32{0}}) //reset transform
	}
}
func f32(x int16) uint32 { return math.Float32bits(float32(x)) }

func (f *File) push(x Record) {
	if s := uint32(8 + 4*len(x.Data)); s != x.Size {
		panic(fmt.Sprintf("size mismatch: %d: %v", s, x))
	}
	f.Records = append(f.Records, x)
}

func (f *File) MarshallBinary() []byte {
	r := make([]Record, 0, 3+len(f.Records))
	r = append(r, Record{0x4001, 0x1, 28, []uint32{16, 3686797314, 1, 96, 96}}) //emfplus header: Emf+dual, 96dpi
	r = append(r, Record{0x401e, 0x9, 12, []uint32{0}})                         //anti-alias
	r = append(r, Record{0x401f, 0x3, 12, []uint32{0}})                         //text rendering hint: antialias-grid-fit(3)
	//r = append(r, Record{0x4022, 0x2, 12, []uint32{0}})                                //pixel offset mode: high quality(2)
	r = append(r, Record{0x4030, 2, 0x10, []uint32{0x4, math.Float32bits(1.0 / 8.0)}}) //scale by 8: max coordinate values are now 4015 (int16)
	r = append(r, f.Records...)
	r = append(r, Record{0x4002, 0x0, 12, []uint32{0}}) //emfplus-eof

	n := 0
	for _, x := range r {
		n += 2 + len(x.Data)
	}
	s := uint32(4 * (1 + n))

	var b bytes.Buffer
	we := func(x Emf) {
		binary.Write(&b, binary.LittleEndian, x.Type)
		binary.Write(&b, binary.LittleEndian, x.Size)
		for _, u := range x.Data {
			binary.Write(&b, binary.LittleEndian, u)
		}
	}
	we(Emf{0x46, 12 + s, []uint32{s, 726027589}})

	for _, x := range r {
		binary.Write(&b, binary.LittleEndian, x.Type)
		binary.Write(&b, binary.LittleEndian, x.Flags)
		binary.Write(&b, binary.LittleEndian, x.Size)
		for _, u := range x.Data {
			binary.Write(&b, binary.LittleEndian, u)
		}
	}
	we(Emf{0xe, 20, []uint32{0, 16, 20}}) //emf-eof

	h := f.Header
	h.Bytes = uint32(len(b.Bytes()) + binary.Size(h))
	h.Records = uint32(3)

	var o bytes.Buffer
	binary.Write(&o, binary.LittleEndian, h)
	o.Write(b.Bytes())
	return o.Bytes()
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
