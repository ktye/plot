package emfplus

import (
	"bytes"
	"encoding/hex"
	"image/png"
	"os"
	"testing"
)

func TestUni(t *testing.T) {
	s := "ab"
	u, c := uni(s)
	if c != 2 {
		t.Fatal()
	}
	if u[0] != uint32(0x00620061) {
		t.Fatalf("%s: %x", s, u)
	}
}

func TestPolygon(t *testing.T) {
	f := New(50, 50)
	f.FillPolygon(0xffff0000, []int16{100, 200, 300, 400, 400, 300, 200, 100}, []int16{100, 200, 100, 200, 400, 300, 400, 300})
	b := f.MarshallBinary()
	_ = b
	os.WriteFile("polygon.emf", b, 0744)
}

func TestEmfPlusSmall(t *testing.T) {
	f := New(20, 20)
	f.Pen(5*8, 0xffaa0000)
	f.DrawEllipse(1, 0, 0, 8*20, 8*20)
	b := f.MarshallBinary()
	_ = b
	//os.WriteFile("small.emf", b, 0744)
}

func TestEmfPlus(t *testing.T) {
	pngdata, e := hex.DecodeString("89504e470d0a1a0a0000000d4948445200000010000000100803000000282d0f530000000467414d410000b18f0bfc6105000000206348524d00007" +
		"a26000080840000fa00000080e8000075300000ea6000003a98000017709cba513c000000ae504c5445000000f1592af1592af1592af1592af1592af" +
		"1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af" +
		"1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af" +
		"1592af1592af1592af1592af1592af1592af1592af1592af1592af1592af1592affffff450bee170000003874524e5300006b6de48a95f16fc40820e" +
		"2f6f2fb6ce36035b49d0717ce6e09911a068c0a0ef80b0553ed8d3777fde5b80fbdc82d5f79f4cdd2f9c9df7a977e5500000001624b474439d700954" +
		"00000000774494d4507e3081f0d1815831020650000007b4944415418d36360606064b2800226460610c01460666165634716e0e0e4e2e6c114e0e5e" +
		"317401210e4131216e140161015139790840b48714bcbc8ca5920b4c88b293072282219aac4a8cca8a28a2420afa0a6cea8a1896c8b96b68e8e2edc50" +
		"3d717d167e03314366a8808091b189a9058b99394ecfa1080000520612531f70dc980000002574455874646174653a63726561746500323031392d303" +
		"82d33315431333a32343a32312b30323a30301f589ade0000002574455874646174653a6d6f6469667900323031392d30382d33315431333a32343a32" +
		"312b30323a30306e052262000000577a5458745261772070726f66696c65207479706520697074630000789ce3f20c0871562828ca4fcbcc49e552000" +
		"3230b2e630b1323134b9314031320448034c3640323b35420cbd8d4c8c4ccc41cc407cb8048a04a2e00ea171174f24235950000000049454e44ae426082")
	if e != nil {
		t.Fatal(e)
	}

	f := New(500, 400)
	f.Records = []Record{
		//Record{0x4030, 0x2, 16, []uint32{4, 1065353216}}, //page transform: f32(1.0)
		//Record{0x401e, 0x5, 12, []uint32{0}},                                                            //anti-alias
		//Record{0x4008, 0x400, 40, []uint32{28, 3686797314, 0, 268435456, 0, 0, 1141637120, 1137836032}}, //object
		//Record{0x4034, 0x0, 12, []uint32{0}},                                                            //clip-region
		//Record{0x401e, 0x0, 12, []uint32{0}},                                                            //anti-alias
		//Record{0x400a, 0xc000, 28, []uint32{16, 4294967295, 1, 0, 27525680}},                            //fill-rects
		//Record{0x4031, 0x0, 12, []uint32{0}},                                 //reset clip
		//Record{0x401e, 0x6, 12, []uint32{0}},                                 //anti-alias
		//Record{0x400a, 0xc000, 28, []uint32{16, 4294967295, 1, 0, 27525680}}, //fill rects

		//Record{0x401e, 0x9, 12, []uint32{0}}, //anti-alias

		//     object  pen                     version     0  flag unit w:0.6666  join/off   brush     solid aarrggbb:ff0072bd
		//Record{0x4008, 0x201, 52, []uint32{40, 3686797314, 0, 136, 0, 1059760811, 2, 0, 3686797314, 0, 4278219453}}, //object pen(2)
		//          c(type)(id)
		//              2   1

		//     object  path                    version     #  flags  point     point    align/pad
		//Record{0x4008, 0x302, 36, []uint32{24, 3686797314, 2, 24576, 24510537, 2032123, 256}}, //object path(3)
		//Record{0x4015, 0x2, 16, []uint32{4, 1}},                                               //draw-path
	}

	f.Pen(1, 0xffaa0000)
	f.DrawEllipse(1, 50, 50, 50, 50)
	f.FillEllipse(0xff0000aa, 50, 50, 20, 20)
	//f.DrawPolyline(1, true, []int16{0, 100, 200, 300}, []int16{0, 50, 0, 70})

	w, h := pngsize(pngdata)
	im := f.Png(w, h, pngdata)
	f.DrawImage(im, 300, 100, int16(w), int16(h))

	//f.LineSegments(1, []int16{0, 100, 200}, []int16{0, 100, 200}, []int16{0, 0, 0}, []int16{200, 200, 200})
	f.LineSegments(1, []int16{0, 100}, []int16{200, 100}, []int16{100, 0}, []int16{100, 200}) //todo

	fn := f.Font(16, "Consolas")
	f.Text(100, 100, "0.123", fn, 1, false, 0xff00aa00)

	f.FillRects(0xff0000ff, []int16{10, 20}, []int16{200, 210}, []int16{50, 60}, []int16{70, 80})
	f.DrawRects(1, []int16{200, 220}, []int16{200, 210}, []int16{50, 60}, []int16{70, 80})

	//f.FillPolygon(0xffff0000, []int16{100, 200, 300, 400, 400, 300, 200, 100}, []int16{100, 200, 100, 200, 400, 300, 400, 300})

	//f.Brush(0xff00ff00)
	//f.FillEllipse(0xff00ff00, 50, 50, 200, 200)

	b := f.MarshallBinary()
	_ = b
	//os.WriteFile("example.emf", b, 0744)
}
func pngsize(b []byte) (int, int) {
	m, e := png.Decode(bytes.NewReader(b))
	if e != nil {
		fatal(e)
	}
	return m.Bounds().Dx(), m.Bounds().Dy()
}
