package wmf

import (
	"encoding/binary"
	"fmt"
	"io"
)

func Read(r io.Reader) error {
	var h Header
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
		fmt.Printf("0x%04x #%d %v\n", u[0], s, u[1:])
	}
}
