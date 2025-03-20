package emfplus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

func New(w, h int) *File {
	f := File{
		Header: Header{Type: 1, Size: 108, Bounds1: 0, Bounds2: 0, Bounds3: 560, Bounds4: 420,
			Frame1: 0, Frame2: 0, Frame3: 15343, Frame4: 11484, Signature: 1179469088, Version: 65536,
			Bytes: 1228, Records: 34, Handles: 2, Reserved: 0,
			NDesc: 0, OffDesc: 0, NPals: 0,
			DevWidth: 1920, DevHeight: 1080, MilliX: 527, MilliY: 296,
			PixelFormat: 0, OffPixel: 0, Bopengl: 0,
			MicroX: 527000, MicroY: 296000},
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

func (f *File) Pen(linewidth int, color uint32) uint8 { //color: 0xaarrggbb
	id := f.objects
	lw := math.Float32bits(float32(linewidth))
	f.push(Record{0x4008, 0x200 | uint16(id), 52 - 8, []uint32{40 - 8, 3686797314, 0, 0, 0, lw, 3686797314, 0, color /*4278219453*/}})
	f.objects++
	return id
}
func (f *File) Brush(color uint32) uint8 {
	id := f.objects
	f.push(Record{0x4008, 0x100 | uint16(id), 24, []uint32{12, 3686797314, 0, color}})
	f.objects++
	return id
}
func rect(left, top, right, bottom int16) (u1, u2 uint32) {
	return uint32(left) | uint32(top)<<16, uint32(right) | uint32(bottom)<<16
}

func (f *File) DrawEllipse(pen uint8, left, top, right, bottom int16) {
	u1, u2 := rect(left, top, right, bottom)
	f.push(Record{0x400f, 0x4000 | uint16(pen), 0x14, []uint32{0x8, u1, u2}})
}
func (f *File) FillEllipse(color uint32, left, top, right, bottom int16) {
	u1, u2 := rect(left, top, right, bottom)
	f.push(Record{0x400e, 0xc000, 24, []uint32{12, color, u1, u2}})
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
	f.objects++
	count := uint32(len(x0))
	u := make([]uint32, 2*len(x0))
	t := make([]byte, 2*len(x0))
	for i := range x0 {
		u[2*i+0] = uint32(x0[i]) | uint32(y0[i]<<16)
		u[2*i+1] = uint32(x1[i]) | uint32(y1[i]<<16)
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
