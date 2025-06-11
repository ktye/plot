package plot

import (
	"encoding/json"
	"testing"
)

func TestJson(t *testing.T) {
	c0 := Caption{Columns: []CaptionColumn{
		CaptionColumn{Data: []float64{3.14, 2}},
		CaptionColumn{Data: []complex128{10 + 20i, 30 + 40i}},
	}}
	p0 := Plot{Lines: []Line{Line{X: []float64{1, 2}, C: []complex128{1 + 2i, 3 + 4i, 5 + 6i}}}}
	p0.Caption = &c0
	p := Plots{p0}
	b, e := json.Marshal(p)
	if e != nil {
		t.Fatal(e)
	}
	//fmt.Println(string(b))

	var r Plots
	json.Unmarshal(b, &r)
	if z := r[0].Lines[0].C[1]; z != 3+4i {
		t.Fatalf("expected 3+4i got %v\n", z)
	}
	if z := r[0].Caption.Columns[1].Data.([]complex128)[1]; z != 30+40i {
		t.Fatalf("expected 30+40i got %v\n", z)
	}
	//fmt.Println(r[0].Caption)
}
