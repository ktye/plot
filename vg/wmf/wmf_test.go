package wmf

import (
	"bytes"
	"os"
	"testing"
)

func TestWmf(t *testing.T) {
	w := New(1000, 1000)
	w.LineTo(500, 500)
	w.MoveTo(0, 1000)
	w.LineTo(400, 400)
	b := w.Bytes()
	e := os.WriteFile("out.wmf", b, 0744)
	if e != nil {
		t.Fatal(e)
	}
}

func TestRead(t *testing.T) {
	b, e := os.ReadFile("1.wmf")
	if e != nil {
		t.Fatal(e)
	}
	Read(bytes.NewReader(b))
}
