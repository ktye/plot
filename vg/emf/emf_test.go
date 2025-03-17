package emf

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
