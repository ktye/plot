package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ktye/plot/vg/emf"
)

func main() {
	b, e := os.ReadFile(os.Args[1])
	fatal(e)
	r := bytes.NewReader(b)

	var h emf.Header
	fatal(binary.Read(r, binary.LittleEndian, &h))
	fmt.Printf("EmfHeader%s\n", strings.ReplaceAll(fmt.Sprintf("%+v", h), " ", ","))
	if h.Size != 108 {
		fmt.Println("header size too large", h.Size)
		r.Seek(int64(h.Size), io.SeekStart)
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
			fmt.Printf("f(Emf{0x%x,%d,[]uint32{%d,%d}})\n", u, s, x[0], x[1])
			x = x[2:]
			for len(x) > 0 {
				t, fl := uint16(x[0]), uint16(x[0]>>16)
				s = x[1]
				if x[2] != s-12 {
					fmt.Println("datasize mismatch")
				}
				n := s / 4
				v := vec(x[2:n])
				fmt.Printf("p(Record{0x%x,0x%x,%d,[]uint32{%s}})\n", t, fl, s, v)
				x = x[n:]
			}
		} else {
			//fmt.Printf("%x %d %v\n", u, s, x)
			v := vec(x)
			fmt.Printf("f(Emf{0x%x,%d,[]uint32{%s}})\n", u, s, v)
		}
	}
}
func vec(x []uint32) string {
	r := strings.ReplaceAll(fmt.Sprintf("%v", x), " ", ",")
	return r[1 : len(r)-1]
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
