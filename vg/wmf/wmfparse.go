//go:build ignore

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type Header struct {
	Type       uint16
	HdrSize    uint16
	Version    uint16
	TotalSize  uint32
	NumObjects uint16
	MaxRecord  uint32
	Unused     uint16
}
type Phead struct {
	Key                      uint32
	Ondisk                   uint16 //0(ondisk) 1(mem)
	Left, Top, Right, Bottom int16
	Inch                     uint16
	Reserved                 uint32
	Sum                      uint16
}

type Emfhead struct {
	Type, Size                         uint32
	Bounds1, Bounds2, Bounds3, Bounds4 int32
	Frame1, Frame2, Frame3, Frame4     int32
	Signature, Version, Bytes, Records uint32
	Handles, Reserved                  uint16
	NDesc, OffDesc, NPals              uint32
	DevWidth, DevHeight                uint32
	MilliX, MilliY                     uint32
}

func pemf(b []byte) { //print emf header
	r := bytes.NewReader(b)
	var h Emfhead
	fatal(binary.Read(r, binary.LittleEndian, &h))
	fmt.Printf("emf: %+v\n", h)
}
func phdr(b []byte) {
	var p Phead
	fatal(binary.Read(bytes.NewReader(b), binary.LittleEndian, &p))
	fmt.Printf("phead: %+v\n", p)
}
func main() {
	b, e := os.ReadFile(os.Args[1])
	fatal(e)
	if len(b) > 4 && b[0] == 1 && b[1] == 0 && b[2] == 0 && b[3] == 0 {
		pemf(b)
		return
	}
	if len(b) > 22 && b[0] == 0xd7 && b[1] == 0xcd { //skip meta placeable record
		phdr(b[:22])
		b = b[22:]
	}
	r := bytes.NewReader(b)

	var h Header
	fatal(binary.Read(r, binary.LittleEndian, &h))
	fmt.Printf("%+v\n", h)
	for {
		var s uint32
		e := binary.Read(r, binary.LittleEndian, &s)
		if e == io.EOF {
			return
		}
		u := make([]uint16, s-2)
		fatal(binary.Read(r, binary.LittleEndian, u))
		fmt.Printf("0x%04x %v\n", u[0], u[1:])
		if u[0] == 0 {
			return
		}
	}
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
