package vg

import (
	"image"
	"image/color"
)

type Drawer interface {
	Reset()
	Size() (int, int)
	SubImage(image.Rectangle) Drawer
	Bounds() image.Rectangle
	Paint()
	Clear(color.Color)
	Color(color.Color)
	Line(Line)
	Circle(Circle)
	Rectangle(Rectangle)
	Triangle(Triangle)
	Ray(Ray)
	Text(Text)
	Font(bool)
	ArrowHead(ArrowHead)

	FloatTics(FloatTics)
	FloatText(FloatText)
	FloatTextExtent(FloatText) (int, int, int, int)
	FloatBars(FloatBars)
	FloatCircles(FloatCircles)
	FloatEnvelope(FloatEnvelope)
	FloatPath(FloatPath)
}
type PngEmbedder interface {
	Embed(int, int, []byte)
}

/*
func EncodeToPng(m image.Image) (string, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, m); err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
*/
