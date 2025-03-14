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
		fmt.Printf("%x #%d %d\n", u, s, x)
	}
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
