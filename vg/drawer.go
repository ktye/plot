package vg

import (
	"image"
	"image/color"

	"golang.org/x/image/font"
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
	Font(font.Face)
	ArrowHead(ArrowHead)

	FloatTics(FloatTics)
	FloatText(FloatText)
	FloatTextExtent(FloatText) (int, int, int, int)
	FloatBars(FloatBars)
	FloatCircles(FloatCircles)
	FloatEnvelope(FloatEnvelope)
	FloatPath(FloatPath)
}
