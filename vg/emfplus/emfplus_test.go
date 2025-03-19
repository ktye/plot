package emfplus

import (
	"os"
	"testing"
)

const writeToDisk = true

func TestEmf(t *testing.T) {
	example(t)
}
func example(t *testing.T) {
	f := New(500, 500)
	f.push(Record{0x401f, 0x3, 12, []uint32{0}})
	f.push(Record{0x4030, 0x2, 16, []uint32{4, 1065353216}}) //page-scale: 1
	f.push(Record{0x4022, 0x4, 12, []uint32{0}})
	f.push(Record{0x401e, 0x9, 12, []uint32{0}})
	f.push(Record{0x4021, 0x7, 12, []uint32{0}})

	f.push(Record{0x402a, 0x0, 36, []uint32{24, 975962800, 0, 0, 975962800, 1159137963, 1151722155}})
	f.push(Record{0x400a, 0x8000, 36, []uint32{24, 4284193749, 1, 0, 0, 1230978560, 1230978560}})

	f.push(Record{0x4008, 0x200, 52, []uint32{40, 3686797314, 0, 144, 0, 1179021312, 1090519040, 0, 3686797314, 0, 4282479004}})
	f.push(Record{0x4008, 0x301, 60, []uint32{48, 3686797314, 4, 0, 0, 0, 1230978560, 0, 1230978560, 1230978560, 0, 1230978560, 2164326656}})
	f.push(Record{0x4015, 0x1, 16, []uint32{4, 0}})
	//f.Ellipse(10, 10, 100, 100)
	/*
		f.AntiAlias()
		f.CreatePen(Pen{Width: 5, Color: Red})
		f.CreateBrush(Brush{Color: Blue})
		f.Select(0)
		f.Rectangle(10, 10, 410, 210)
		//f.SetTextColor(0)
		f.CreateFont(Font{Height: -20, Weight: 400, Face: "Consolas"})
		f.Select(2)
		f.Text(100, 100, "Hello 123-xyz")
		f.MoveTo(0, 0)
		f.LineTo(500, 250)
		f.Polyline([]int16{100, 200, 300, 400}, []int16{500, 400, 500, 400})
		f.Select(1)
		f.Ellipse(10, 10, 100, 100)
	*/
	f.write(t, "example.emf")
}
func (f *File) write(t *testing.T, file string) {
	b := f.MarshallBinary()
	if writeToDisk {
		e := os.WriteFile(file, b, 0744)
		if e != nil {
			t.Fatal(e)
		}
	}
}
