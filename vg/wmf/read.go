package wmf

import (
	"encoding/binary"
	"fmt"
	"io"
)

func Read(r io.Reader) error {
	var h hdr
	e := binary.Read(r, binary.LittleEndian, &h)
	if e != nil {
		return e
	}
	fmt.Printf("%+v\n", h)
	for {
		var s uint32
		e := binary.Read(r, binary.LittleEndian, &s)
		if e == io.EOF {
			return nil
		}
		u := make([]uint16, s-2)
		e = binary.Read(r, binary.LittleEndian, u)
		if e != nil {
			return e
		}
		fmt.Printf("0x%04x #%d\n", u[0], s)
	}
}

type hdr struct {
	Key      uint32
	Hwmf     uint16
	Left     int16
	Top      int16
	Bottom   int16
	Right    int16
	Inch     int16
	Reserved uint32
	Checksum uint16

	Type       uint16
	HdrSize    uint16
	Version    uint16
	TotalSize  uint32
	NumObjects uint16
	MaxRecord  uint32
	Unused     uint16
}
