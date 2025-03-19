package emfplus

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ktye/plot/vg/emf"
)

func New(w, h int) *File {
	f := File{
		Header: Header{
			Type:         0x4001,
			Flags:        1,
			Size:         0x1c,
			DataSize:     0x10,
			Version:      0xdbc01002,
			EmfPlusFlags: 0,   //1 << 31,
			DpiX:         600, //96,
			DpiY:         600, //96,
		},
		width:  w,
		height: h,
	}
	return &f
}

type File struct {
	Header
	Records       []Record
	objects       uint32
	width, height int
}
type Header struct {
	Type, Flags              uint16
	Size, DataSize, Version  uint32
	EmfPlusFlags, DpiX, DpiY uint32
}
type Record struct {
	Type  uint16
	Flags uint16
	Size  uint32
	Data  []uint32
}

func rect(left, top, right, bottom int16) (u1, u2 uint32) {
	return uint32(left) | uint32(top)<<16, uint32(right) | uint32(bottom)<<16
}

func (f *File) Ellipse(left, top, right, bottom int16) {
	u1, u2 := rect(left, top, right, bottom)
	f.push(Record{0x400f, 0x2, 0x14, []uint32{0x8, u1, u2}})
}

func (f *File) push(x Record) {
	if s := uint32(8 + 4*len(x.Data)); s != x.Size {
		panic(fmt.Sprintf("size mismatch: %d: %v", s, x))
	}
	f.Records = append(f.Records, x)
}
func (f *File) MarshallBinary() []byte {
	var r []Record
	for _, x := range f.Records {
		r = append(r, x)
	}
	r = append(r, Record{0x4002, 0, 0xc, []uint32{0}}) //eof

	m := emf.New(f.width, f.height)

	u := make([]uint32, binary.Size(f.Header)/4)
	var b bytes.Buffer
	fatal(binary.Write(&b, binary.LittleEndian, f.Header))
	fatal(binary.Read(&b, binary.LittleEndian, u))
	m.EmfplusComment(u)

	u = nil
	for _, x := range r {
		u = append(u, uint32(x.Type)|uint32(x.Flags)<<16, x.Size)
		u = append(u, x.Data...)
	}

	m.EmfplusComment(u)
	return m.MarshallBinary()
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
