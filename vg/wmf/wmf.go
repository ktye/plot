package wmf

// [MS-WMF]-220625.pdf
// https://wvware.sourceforge.net/caolan/ora-wmf.html
// https://github.com/g21589/wmf2canvas/blob/master/js/wmf.js

import (
	"bytes"
	"encoding/binary"
)

func New(w, h int) *wf {
	r := wf{
		width:  int16(w),
		height: int16(h),
		inch:   1440,
	}
	return &r
}
func (w *wf) MoveTo(x, y int) {
	b := []byte{5, 0, 0, 0, 0x14, 0x02, 0, 0, 0, 0} //2.3.5.4
	binary.LittleEndian.PutUint16(b[6:], uint16(int16(y)))
	binary.LittleEndian.PutUint16(b[8:], uint16(int16(x)))
	w.push(b)
}
func (w *wf) LineTo(x, y int) {
	b := []byte{5, 0, 0, 0, 0x13, 0x02, 0, 0, 0, 0} //2.3.3.10
	binary.LittleEndian.PutUint16(b[6:], uint16(int16(y)))
	binary.LittleEndian.PutUint16(b[8:], uint16(int16(x)))
	w.push(b)
}
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

func (w *wf) push(b []byte) { w.obj = append(w.obj, b) }

func checksum(b []byte) (r uint16) {
	var v [10]uint16
	e := binary.Read(bytes.NewReader(b), binary.LittleEndian, &v)
	if e != nil {
		panic(e)
	}
	for _, x := range v {
		r ^= x
	}
	return r
}
