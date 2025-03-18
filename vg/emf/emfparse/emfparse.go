package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/ktye/plot/vg/emf"
)

func main() {
	b, e := os.ReadFile(os.Args[1])
	fatal(e)
	r := bytes.NewReader(b)

	var h emf.Header
	fatal(binary.Read(r, binary.LittleEndian, &h))
	fmt.Printf("%+v\n", h)
	if h.Size != 108 {
		fmt.Printf("unsupported header size: #%d\n", h.Size)
	}
	for {
		var u, s uint32
		e := binary.Read(r, binary.LittleEndian, &u)
		if e == io.EOF {
			break
		}
		fatal(binary.Read(r, binary.LittleEndian, &s))
		x := make([]uint32, (s/4)-2)
		fatal(binary.Read(r, binary.LittleEndian, x))
		if u == 0x46 && len(x) > 1 && x[1] == 0x2b464d45 { //"EMF+"
			fmt.Printf("%x #%d EMF+\n", u, s)
			x = x[2:]
			for len(x) > 0 {
				t := uint16(x[0])
				s = x[1]
				n := s / 4
				fmt.Printf("+%x #%d %v\n", t, s, x[2:n])
				x = x[n:]
			}
		} else {
			fmt.Printf("%x #%d %d\n", u, s, x)
		}
	}
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
