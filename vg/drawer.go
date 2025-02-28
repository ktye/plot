package vg

import (
	"image/color"

	"golang.org/x/image/font"
)

type Drawer interface {
	Size() (int, int)
	SubImage(int, int, int, int) Drawer
	Paint()
	Clear(color.Color)
	Color(color.Color)
	Line(Line)
	Circle(Circle)
	Triangle(Triangle)
	Ray(Ray)
	Text(Text)
	Font(font.Face)
	ArrowHead(ArrowHead)

	FloatTics(FloatTics)
	FloatText(FloatText)
	FloatBars(FloatBars)
	FloatCircles(FloatCircles)
	FloatEnvelope(FloatEnvelope)
	FloatPath(FloatPath)
}
