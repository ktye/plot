// Package raster contains rasterization algorithms.
//
// Opposed to the image/draw conventions, the coordinates x and y
// are located at the pixel centers.
//
// Used Bresenham algorithm source:
// http://members.chello.at/~easyfilter/bresenham.html
package raster

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/ktye/plot/xmath"
)

// A Setter can set a single pixel at x and y.
type Setter interface {
	Set(int, int)
}

// Image implements a setter by wrapping the underlying draw.Image
// and color type.
type Image struct {
	draw.Image
	color.Color
}

// Set a pixel to the image with the current color.
func (im *Image) Set(x, y int) {
	im.Image.Set(x, y, im.Color)
}

// Line rasterizes a straight 1px line.
func Line(im Image, x0, y0, x1, y1 int) {
	var dx, dy, sx, sy, e, e2 int
	dx = abs(x1 - x0)
	dy = -abs(y1 - y0)
	if sx = -1; x0 < x1 {
		sx = 1
	}
	if sy = -1; y0 < y1 {
		sy = 1
	}
	e = dx + dy
	for {
		im.Set(x0, y0)
		if x0 == x1 && y0 == y1 {
			break
		}
		if e2 = 2 * e; e2 >= dy {
			e += dy
			x0 += sx
		} else if e2 <= dx {
			e += dx
			y0 += sy
		}
	}
}

// Circle rasterizes a circle at center xm and ym with radius r.
// The width and height in number of pixels of the circle is 2*r + 1.
func Circle(im Image, xm, ym, r int) {
	var x, y, e int
	x = -r
	e = 2 - 2*r
	for x < 0 {
		im.Set(xm-x, ym+y)
		im.Set(xm-y, ym-x)
		im.Set(xm+x, ym-y)
		im.Set(xm+y, ym+x)
		r = e
		if r <= y {
			y++
			e += 2*y + 1
		}
		if r > x || e > y {
			x++
			e += 2*x + 1
		}
	}
}

// FillCircle fills a circle with center at xm, ym and radius r.
func FillCircle(img Image, xm, ym, r int) {
	draw.DrawMask(img.Image, img.Image.Bounds(), &image.Uniform{img.Color}, image.ZP, &circle{image.Point{xm, ym}, r}, image.ZP, draw.Over)
}

type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r-1, c.p.Y-c.r-1, c.p.X+c.r+1, c.p.Y+c.r+1)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X), float64(y-c.p.Y), float64(c.r)+0.5
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

// Rectangle draws a rectangle.
func Rectangle(im Image, x0, y0, x1, y1 int) {
	Line(im, x0, y0, x1, y0)
	Line(im, x1, y0, x1, y1)
	Line(im, x1, y1, x0, y1)
	Line(im, x0, y1, x0, y0)
}

// FillRectangle fills a rectangle.
func FillRectangle(im Image, x0, y0, x1, y1 int) {
	draw.Draw(im.Image, image.Rect(x0, y0, x1+1, y1+1), &image.Uniform{im.Color}, image.ZP, draw.Src)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// CoordinateSystem defines the upper left point and the lower right point
// of the image in floating-point coordinates.
type CoordinateSystem struct {
	X0, Y0, X1, Y1 float64
}

// FloatLines draws connected lines given in floating point coordinates.
func FloatLines(im Image, x, y []float64, cs CoordinateSystem) { //todo ztranslate
	doStart := true
	var X0, Y0 int
	for i := range x {
		if math.IsNaN(x[i]) || math.IsNaN(y[i]) {
			doStart = true
			continue
		}
		X, Y := transform(x[i], y[i], cs, im.Image.Bounds())
		if doStart {
			X0, Y0 = X, Y
			doStart = false
		} else {
			Line(im, X0, Y0, X, Y)
			X0, Y0 = X, Y
		}
	}
}

func transform(x, y float64, cs CoordinateSystem, bounds image.Rectangle) (int, int) {
	x0, x1 := float64(bounds.Min.X), float64(bounds.Max.X)
	y0, y1 := float64(bounds.Min.Y), float64(bounds.Max.Y)
	return clip(xmath.Scale(x, cs.X0, cs.X1, x0, x1)), clip(xmath.Scale(y, cs.Y0, cs.Y1, y0, y1))
}

// We clip to 8192 pixels.
// No image should be larger than that.
// This ensures that not too many invisible pixels are drawn, when zoomed in
// for large values.
func clip(f float64) int {
	if f > 8192 {
		return 8192
	}
	if f < -8192 {
		return -8192
	}
	return int(f)
}
