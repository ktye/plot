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

func main() {
	b, e := os.ReadFile(os.Args[1])
	fatal(e)
	if len(b) > 22 && b[0] == 0xd7 && b[1] == 0xcd { //skip meta placeable record
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
