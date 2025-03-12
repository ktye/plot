package wmf

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

const writeToDisk = true

func (f *File) write(t *testing.T, file string) {
	b := f.MarshallBinary()
	if writeToDisk {
		e := os.WriteFile(file, b, 0744)
		if e != nil {
			t.Fatal(e)
		}
	}
}

func line(t *testing.T) {
	f := New(400, 300)
	f.LineTo(400, 300)
	f.MoveTo(0, 400)
	f.LineTo(400, 0)
	f.write(t, "line.wmf")
}
func text(t *testing.T) {
	f := New(100, 20)
	f.Text(10, 10, "häöüÄÖÜß€°µe")
	f.write(t, "text.wmf")
}
func ax(t *testing.T) {
	f := New(200, 100)
	f.CreatePen(Pen{Width: 4, Color: Blue})
	//f.CreateBrush(Brush{Color: Color{R: 255}})
	f.Select(0)
	//f.SelectClip(1)
	f.Rectangle(0, 0, 150, 70)
	f.Text(10, 10, "abc A B C")
	f.write(t, "ax.wmf")
}
func circ(t *testing.T) {
	f := New(800, 600)
	f.Ellipse(0, 0, 600, 600)
	f.write(t, "circ.wmf")
}
func ell(t *testing.T) {
	f := New(200, 100)
	f.CreatePen(Pen{Width: 5, Color: Red})
	f.CreateBrush(Brush{Color: Blue})
	f.Select(0)
	f.Select(1)
	f.Ellipse(0, 0, 200, 100)
	f.write(t, "ell.wmf")
}
func poly(t *testing.T) {
	f := New(300, 200)
	f.CreatePen(Pen{Width: 5, Color: Green})
	f.Select(0)
	f.Polygon([]int16{0, 50, 100, 150, 200, 250, 300}, []int16{0, 200, 0, 200, 0, 200, 0})
	f.write(t, "poly.wmf")
}
func font(t *testing.T) {
	f := New(500, 100)
	//f.SetMapMode(1)
	//f.Rectangle(0, 0, 200, 50)
	f.CreateFont(Font{Height: 50, Face: "Consolas"})
	f.Select(0)
	f.Text(0, 0, "Hans Werner")
	f.write(t, "font.wmf")
}
func vert(t *testing.T) {
	f := New(300, 200)
	//f.SetMapMode(1)
	//f.Rectangle(0, 0, 200, 50)
	f.Text(0, 0, "horizontal text")
	f.CreateFont(Font{Height: 20, Face: "Arial", Escapement: 900, Orientation: 900})
	f.Select(0)
	f.Text(0, 150, "vertical text")
	f.write(t, "vert.wmf")
}
func align(t *testing.T) {
	f := New(200, 200)
	f.CreatePen(Pen{Width: 2, Color: Red})
	f.CreateBrush(Brush{Color: LightGrey})
	f.Select(0)
	f.Select(1)
	f.Rectangle(50, 50, 150, 150)
	f.Text(50, 50, "ABC")
	f.SetTextAlign(2)
	f.Text(150, 50, "abc")
	f.SetTextAlign(8)
	f.Text(50, 150, "DEF")
	f.SetTextAlign(2 + 8)
	f.Text(150, 150, "def")
	f.SetTextAlign(6)
	f.Text(100, 50, "123")
	f.SetTextAlign(6 + 0x18)
	f.Text(100, 150, "456")
	//f.Text(100, 100, "A")
	f.write(t, "align.wmf")
}
func clip(t *testing.T) {
	f := New(300, 200)
	f.CreatePen(Pen{Width: 5, Color: Green})
	f.Select(0)
	f.IntersectClipRect(50, 50, 250, 200)
	f.Polygon([]int16{0, 50, 100, 150, 200, 250, 300}, []int16{0, 200, 0, 200, 0, 200, 0})
	f.write(t, "clip.wmf")
}

func TestWmf(t *testing.T) {
	line(t)
	text(t)
	ax(t)
	circ(t)
	ell(t)
	poly(t)
	font(t)
	vert(t)
	align(t)
	clip(t)
}

func TestRead(t *testing.T) {
	t.Skip()
	files := []string{"ex.wmf", "ax.wmf", "ell.wmf", "text.wmf", "line.wmf"}
	for _, f := range files {
		b, e := os.ReadFile(f)
		if e != nil {
			t.Fatal(e)
		}
		fmt.Print(f + ": ")
		Read(bytes.NewReader(b))
	}
}
