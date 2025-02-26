package wmf

// [MS-WMF]-220625.pdf
// https://wvware.sourceforge.net/caolan/ora-wmf.html
// https://github.com/g21589/wmf2canvas/blob/master/js/wmf.js

import (
	"bytes"
	"encoding/binary"
)

func New(w, h int) *File {
	f := File{
		Header: Header{
			Key:     2596720087,
			Right:   int16(w),
			Bottom:  int16(h),
			Type:    2, //1(in-mem) 2(on-disk)
			HdrSize: 9,
			Version: 0x0300,
		},
	}
	return &f
}

type File struct {
	Header
	Records []Record
}
type Header struct {
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

func (f *File) MoveTo(x, y int) {
	f.push(Record{5, 0x0214, []uint16{uint16(int16(y)), uint16(int16(x))}})
}
func (f *File) LineTo(x, y int) {
	f.push(Record{5, 0x0213, []uint16{uint16(int16(y)), uint16(int16(x))}})
}
func (f *File) push(x Record) { x.Size = uint32(binary.Size(x) / 2); f.Records = append(f.Records, x) }
func (f *File) setChecksum() { //of the leading 10words on the header
	var b bytes.Buffer
	fatal(binary.Write(&b, binary.LittleEndian, f.Header))
	var v [10]uint16
	fatal(binary.Read(bytes.NewReader(b.Bytes()), binary.LittleEndian, &v))
	var s uint16
	for _, x := range v {
		s ^= x
	}
	f.Checksum = s
}
func (f *File) MarshallBinary() ([]byte, error) {
	f.setChecksum()
	f.push(Record{3, 0, nil}) //eof

	m := 0
	for _, r := range f.Records {
		n := binary.Size(r)
		if n > m {
			m = n
		}
	}
	f.MaxRecord = uint32(m / 2)
	f.TotalSize = uint32((binary.Size(f) - 22) / 2)
	//todo number of objects

	var b bytes.Buffer
	e := binary.Write(&b, binary.LittleEndian, f)
	if e != nil {
		return nil, e
	}
	return b.Bytes(), nil
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
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
