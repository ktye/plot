// Package vg implements vector graphics operations on an image.
package vg

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/ktye/plot/xmath"

	"github.com/golang/freetype/raster"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Coordinate systems.
// In this file, there are 3 different coordinates:
//	image pixels as integers (top left to bottom right)
//	fixed.Int26_6 for the rasterizer
//	float64 with a coordinate transformation for all Float* Drawers.
//
// The fixed.Int26_6 is used only because the rasterizer expects it.
// For the interface the two versions pixels (int) and coordinates as float64 serve a
// different purpose.
// Int coordinates are used, whenever the position in the image is important.
// Float coordinates are used to draw data lines or points into an axis.
// This is why they bring their own CoordinateSystem.

// Drawer is an object with a draw method.
// It knows how to rasterize itself on an image.
// Some types which implement the interface are defined in this package, but
// others may be defined by the package clients.
// Drawers should only call painters methods: SetColor, Fill and Stroke.
type drawer interface {
	Draw(*Painter)
}

type Painter struct {
	im           *image.RGBA
	p            raster.Painter
	r            *raster.Rasterizer
	currentColor color.Color
	currentFace  font.Face
	fontDrawer   font.Drawer
	drawers      []drawer
	colors       []color.Color
	faces        []font.Face
	x0, y0       int
}

func (p *Painter) Add(d drawer) {
	p.drawers = append(p.drawers, d)
	p.colors = append(p.colors, p.currentColor)
	p.faces = append(p.faces, p.currentFace)
}
func (p *Painter) Paint() {
	for i, d := range p.drawers {
		switch p.p.(type) {
		case *raster.RGBAPainter:
			p.p.(*raster.RGBAPainter).SetColor(p.colors[i])
			p.currentColor = p.colors[i] // used by all drawers who dont rasterize.
			p.fontDrawer.Src = image.NewUniform(p.colors[i])
		}
		p.fontDrawer.Face = p.faces[i]
		d.Draw(p)
	}
}
func (p *Painter) SetColor(c color.Color) {
	p.currentColor = c
}
func (p *Painter) GetColor() color.Color {
	return p.currentColor
}
func (p *Painter) SetFont(f font.Face) {
	p.fontDrawer.Face = f // to calculate extent, before Paint.
	p.currentFace = f
}
func (p *Painter) Stroke(path raster.Path, lineWidth int) {
	p.r.UseNonZeroWinding = true
	p.r.AddStroke(path, fixed.I(lineWidth), raster.SquareCapper, raster.BevelJoiner)
	p.r.Rasterize(p.p)
	p.r.Clear()
}
func (p *Painter) Fill(path raster.Path) {
	p.r.AddPath(path)
	p.r.Rasterize(p.p)
	p.r.Clear()
}
func (p *Painter) Image() image.Image {
	return p.im
}
func NewPainter(im *image.RGBA) *Painter {
	fd := font.Drawer{
		Dst:  im,
		Src:  image.Black,
		Face: basicfont.Face7x13,
	}
	return &Painter{
		im:           im,
		p:            raster.NewRGBAPainter(im),
		r:            raster.NewRasterizer(im.Bounds().Max.X, im.Bounds().Max.Y),
		currentColor: color.Black,
		currentFace:  fd.Face,
		fontDrawer:   fd,
		x0:           im.Bounds().Min.X,
		y0:           im.Bounds().Min.Y,
	}
}

// Text alignment is an integer between 0 and 8:
// 6---5---4
// 7tex8tex3
// 0---1---2
type Text struct {
	X, Y     int
	S        string
	Align    int
	Vertical bool
}

func (t Text) Draw(p *Painter) {
	bounds, _ := font.BoundString(p.fontDrawer.Face, t.S)
	metrics := p.fontDrawer.Face.Metrics()
	x, y, w, h := p.Extent(t)
	ascent := metrics.Ascent.Ceil()
	y += ascent
	x -= bounds.Min.X.Ceil()
	if t.Vertical {
		x, y := t.X, t.Y
		if a := t.Align; a == 1 || a == 5 { //x-correction nyi
			y -= w / 2
		} else if a == 2 || a == 3 || a == 4 {
			y -= w
		}
		vp := NewPainter(image.NewRGBA(image.Rect(0, 0, 1+w, h)))
		r, g, b, _ := p.currentColor.RGBA() //assume black or white
		bg := color.RGBA{255 - uint8(r>>8), 255 - uint8(g>>8), 255 - uint8(b>>8), 255}
		draw.Draw(vp.im, vp.im.Bounds(), image.NewUniform(bg), image.Point{0, 0}, draw.Src) //assume black on white
		vp.SetFont(p.currentFace)
		vp.SetColor(p.currentColor)
		t.Vertical = false
		t.X = 1
		t.Y = 0
		t.Align = 6
		vp.Add(t)
		vp.Paint()
		src := rotate(vp.im)
		draw.Draw(p.im, p.im.Bounds(), src, image.Point{-x, -y}, draw.Src)
	} else {
		p.fontDrawer.Dot = fixed.P(x+p.x0, y+p.y0)
		p.fontDrawer.DrawString(t.S)
	}
}
func rotate(src *image.RGBA) *image.RGBA {
	srcW := src.Bounds().Max.X
	srcH := src.Bounds().Max.Y
	dstW := srcH
	dstH := srcW
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))

	for dstY := 0; dstY < dstH; dstY++ {
		for dstX := 0; dstX < dstW; dstX++ {
			srcX := dstH - dstY - 1
			srcY := dstX

			srcOff := srcY*src.Stride + srcX*4
			dstOff := dstY*dst.Stride + dstX*4

			copy(dst.Pix[dstOff:dstOff+4], src.Pix[srcOff:srcOff+4])
		}
	}

	return dst
}

// Extent returns the Rectangle values which covers the text.
// x, y is the top left corner.
func (p *Painter) Extent(t Text) (x, y, width, height int) {
	x, y = t.X, t.Y
	bounds, _ := font.BoundString(p.fontDrawer.Face, t.S)
	d := bounds.Max.Sub(bounds.Min)
	width = d.X.Ceil()
	metrics := p.fontDrawer.Face.Metrics()
	height = (metrics.Ascent + metrics.Descent).Ceil()

	if t.Align == 1 || t.Align == 5 || t.Align == 8 {
		x -= width / 2
	}
	if t.Align == 2 || t.Align == 3 || t.Align == 4 {
		x -= width
	}
	if t.Align == 7 || t.Align == 8 || t.Align == 3 {
		y -= height / 2
	}
	if t.Align == 0 || t.Align == 1 || t.Align == 2 {
		y -= height
	}

	return x, y, width, height
}

// Pixel
type Pixel struct {
	X, Y int
}

func (x Pixel) Draw(p *Painter) {
	p.im.Set(x.X+p.x0, x.Y+p.y0, p.currentColor)
}

// Rectangle within these limits.
// The complete LineWidth will be painted to the inside.
// X, Y is the top left point (top left of the pixels)
// W, H is the number of pixels in horizontal and vertical direction.
// This is not like image.Rectangle, 1px is added to W and H.
type Rectangle struct {
	X, Y, W, H int
	LineWidth  int
	Fill       bool
}

func (r Rectangle) Draw(p *Painter) {
	rect := image.Rect(r.X+p.x0, r.Y+p.y0, r.X+p.x0+r.W+1, r.Y+p.y0+r.H+1)
	if r.Fill {
		draw.Draw(p.im, rect, image.NewUniform(p.currentColor), image.ZP, draw.Src)
	} else {
		// frame is a alpha mask image which is opaque at the border (linewidth)
		// and transparent in the center.
		frame := image.NewAlpha(image.Rect(0, 0, r.W, r.H))
		draw.Draw(frame, frame.Bounds(), image.NewUniform(color.Opaque), image.ZP, draw.Src)
		innerRect := image.Rect(r.LineWidth, r.LineWidth, r.W-r.LineWidth, r.H-r.LineWidth)
		draw.Draw(frame, innerRect, image.NewUniform(color.Transparent), image.ZP, draw.Src)
		//todo add p.x0 p.y0 ..........vv...vv?
		draw.DrawMask(p.im, image.Rect(r.X, r.Y, r.W+1, r.H+1), image.NewUniform(p.currentColor), image.ZP, frame, image.ZP, draw.Over)
	}
}

type Triangle struct {
	X0, Y0, X1, Y1, X2, Y2 int
	LineWidth              int
	Fill                   bool
}

func (t Triangle) Draw(p *Painter) {
	var path raster.Path
	path = append(path, 0, fixed.I(t.X0+p.x0), fixed.I(t.Y0+p.y0), 0)
	path = append(path, 1, fixed.I(t.X1+p.x0), fixed.I(t.Y1+p.y0), 1)
	path = append(path, 1, fixed.I(t.X2+p.x0), fixed.I(t.Y2+p.y0), 1)
	path = append(path, 1, fixed.I(t.X0+p.x0), fixed.I(t.Y0+p.y0), 1)
	if t.LineWidth > 0 {
		p.Stroke(path, t.LineWidth)
	} else {
		p.Fill(path)
	}
}

// Circle draws a circle within the given image rectangle X, Y, X+D, Y+D.
// X and Y are the upper left coordinates.
// We do not give the center coordinates, because we will be able to draw a
// 1px center line ON a pixel and not in between.
// For this case, the diameter should be an odd number.
// The circle is completely inside, independend on the linewidth, similar to Rectangle
// with W = H = D.
type Circle struct {
	X, Y, D   int
	LineWidth int
	Fill      bool
}

func (c Circle) Draw(p *Painter) {
	// Move the origin to the pixel center.
	R := fixed.Int26_6((c.D * 64) / 2)
	x := fixed.Int26_6(c.X*64) + R
	y := fixed.Int26_6(c.Y*64) + R
	r := R
	if c.Fill == false {
		r -= fixed.Int26_6(c.LineWidth * 32) // substract half the linewidth.
	}
	c26_6 := circle{x + fixed.I(p.x0), y + fixed.I(p.y0), r}
	path := c26_6.getPath()
	if c.Fill {
		p.Fill(path)
	} else {
		p.Stroke(path, c.LineWidth)
	}
}

// circle is used internally by Circle and FloatCircles.
// x, y are center coordinates, r is the radius.
type circle struct {
	x, y, r fixed.Int26_6
}

// getPath approximates a circle by 8 quadrativ curve segments.
func (c circle) getPath() raster.Path {
	d := fixed.Point26_6{c.x, c.y}
	r := c.r
	s := fixed.Int26_6(float64(c.r) * math.Sqrt(2.0) / 2.0)
	t := fixed.Int26_6(float64(c.r) * math.Tan(math.Pi/8))
	P := func(x, y fixed.Int26_6) fixed.Point26_6 {
		return fixed.Point26_6{x, y}
	}
	var path raster.Path
	path.Start(d.Add(P(r, 0)))
	path.Add2(d.Add(P(r, t)), d.Add(P(s, s)))
	path.Add2(d.Add(P(t, r)), d.Add(P(0, r)))
	path.Add2(d.Add(P(-t, r)), d.Add(P(-s, s)))
	path.Add2(d.Add(P(-r, t)), d.Add(P(-r, 0)))
	path.Add2(d.Add(P(-r, -t)), d.Add(P(-s, -s)))
	path.Add2(d.Add(P(-t, -r)), d.Add(P(0, -r)))
	path.Add2(d.Add(P(t, -r)), d.Add(P(s, -s)))
	path.Add2(d.Add(P(r, -t)), d.Add(P(r, 0)))
	return path
}

type LineCoords struct {
	X, Y, DX, DY int
}

// Single Line.
type Line struct {
	LineCoords
	LineWidth int
	Floor     bool
}

func (l Line) Draw(p *Painter) {
	x0, y0, x1, y1 := fixed.I(l.X+p.x0), fixed.I(l.Y+p.y0), fixed.I(l.X+p.x0+l.DX), fixed.I(l.Y+p.y0+l.DY)
	hp := fixed.Int26_6(32) // half a pixel, this makes sure, a 1px line is not shared between 2 pixels.
	if l.Floor && l.LineWidth%2 == 0 {
		hp = 0
	}
	path := raster.Path{0, x0 + hp, y0 + hp, 0, 1, x1 + hp, y1 + hp, 1}
	p.Stroke(path, l.LineWidth)
}
func NewLine(x, y, dx, dy, lw int) Line {
	return Line{LineCoords{x, y, dx, dy}, lw, false}
}

// Ray is a line with length L.
// It starts at X+R*cos(Phi), Y+R*sin(Phi).
type Ray struct {
	X, Y, R, L int
	Phi        float64
	LineWidth  int
}

func (r Ray) Draw(p *Painter) {
	x0 := float64(r.X+p.x0) + float64(r.R)*math.Cos(r.Phi) + 0.5
	y0 := float64(r.Y+p.y0) + float64(r.R)*math.Sin(r.Phi) + 0.5
	x1 := float64(r.X+p.x0) + float64(r.R+r.L)*math.Cos(r.Phi) + 0.5
	y1 := float64(r.Y+p.y0) + float64(r.R+r.L)*math.Sin(r.Phi) + 0.5
	path := raster.Path{
		0, fixed.Int26_6(int(x0 * 64.0)), fixed.Int26_6(int(y0 * 64.0)), 0,
		1, fixed.Int26_6(int(x1 * 64.0)), fixed.Int26_6(int(y1 * 64.0)), 1,
	}
	p.Stroke(path, r.LineWidth)
}

// Lines are unconnected lines with the same width.
type Lines struct {
	Coords    []LineCoords
	LineWidth int
}

func (l Lines) Draw(p *Painter) {
	var path raster.Path
	for _, c := range l.Coords {
		x0, y0, x1, y1 := fixed.I(c.X+p.x0), fixed.I(c.Y+p.y0), fixed.I(c.X+p.x0+c.DX), fixed.I(c.Y+p.y0+c.DY)
		path = append(path, 0, x0, y0, 0, 1, x1, y1, 1)
	}
	p.Stroke(path, l.LineWidth)
}

// CoordinateSystem defines the upper left point and the lower right point
// of the image in floating-point coordinates.
type CoordinateSystem struct {
	X0, Y0, X1, Y1 float64
}

// FloatTics are unconnected lines given in floating points coordinates.
// The lines are drawn perpendicular to the given direction.
type FloatTics struct {
	P          []float64 // Coordinates x for horizontal, y for vertical.
	Q          float64   // The other coordinate, which is the same for every point.
	Horizontal bool      // Direction of the axis, horizontal or vertical.
	LeftTop    bool      // Mark if it is the left or top axis.
	L          int       // line length in pixels.
	LineWidth  int
	CoordinateSystem
	Rect image.Rectangle // Rectangle of the coordinate system.
}

func (f FloatTics) Draw(p *Painter) {
	var path raster.Path
	for _, v := range f.P {
		X, Y := f.Q, v
		if f.Horizontal {
			X, Y = Y, X
		}
		rect := rect26_6(f.Rect)
		x, y := transform(X, Y, f.CoordinateSystem, rect)
		// round to pixel
		x /= 64
		x *= 64
		y /= 64
		y *= 64
		if f.LineWidth%2 == 1 {
			x += 32
			y += 32
		}
		off := fixed.Int26_6(f.L * 64)
		if f.LeftTop == false {
			off = -off
		}
		if f.Horizontal {
			y -= off
		} else {
			x -= off
		}
		path = append(path, 0, x, y, 0)
		if f.Horizontal {
			y += off
		} else {
			x += off
		}
		path = append(path, 1, x, y, 0)
	}
	p.Stroke(path, f.LineWidth)
}

// ArrowHead adds an arrow to a line. The coordinates of a FloatPath can be re-used.
// Only the last two values are used to calculate the direction.
type ArrowHead struct {
	X, Y []float64
	CoordinateSystem
	LineWidth int // scales arrow head size
	Arrow     int // currently ignored, may be used for two side arrows or styles
}

func (a ArrowHead) Draw(p *Painter) {
	if len(a.X) < 2 || len(a.Y) < 2 {
		return
	}
	x, y, xa, ya, xb, yb, xc, yc := a.points(rect26_6(p.im.Bounds()))
	f := func(p float64) fixed.Int26_6 { return fixed.Int26_6(int(64.0 * p)) }
	path := raster.Path{0, x, y, 0, 1, f(xb), f(yb), 1, 1, f(xa), f(ya), 1, 1, f(xc), f(yc), 1, 1, x, y, 1}
	p.Fill(path)
}
func (a ArrowHead) points(r fixed.Rectangle26_6) (x, y fixed.Int26_6, xa, ya, xb, yb, xc, yc float64) {
	x, y = transform(a.X[len(a.X)-1], a.Y[len(a.Y)-1], a.CoordinateSystem, r)
	dx, dy := a.X[len(a.X)-1]-a.X[len(a.X)-2], a.Y[len(a.Y)-1]-a.Y[len(a.Y)-2]
	dx, dy = transformDirection(dx, dy, a.CoordinateSystem, r) // unit vector in pixels (float64)
	A := float64(a.LineWidth * 12)                             // arrow head length (pixels)
	B := A / 5.0                                               // short cathetus length
	//x += fixed.I(p.x0)
	//y += fixed.I(p.y0)
	xa, ya = float64(x)/64.0-A*dx, float64(y)/64.0-A*dy // arrow base point on the line
	xb, yb = xa-B*dy, ya+B*dx                           // corner points
	xc, yc = xa+B*dy, ya-B*dx                           //
	return x, y, xa, ya, xb, yb, xc, yc
}

// FloatPath is a connected line with many points.
// It is given in floating point coordinates with a transformation system.
// For every NaN in the path, a new line is started.
type FloatPath struct {
	X, Y []float64 // point coordinates
	Z    int
	CoordinateSystem
	LineWidth int
}

func (f FloatPath) Draw(p *Painter) {
	var path raster.Path
	doStart := true
	for i := range f.X {
		if math.IsNaN(f.X[i]) || math.IsNaN(f.Y[i]) {
			doStart = true
			continue
		}
		x, y := transform(f.X[i], f.Y[i], f.CoordinateSystem, rect26_6(p.im.Bounds()))
		z := fixed.I(f.Z)
		x += z
		y -= z
		if doStart {
			path = append(path, 0, x, y, 0)
			doStart = false
		} else {
			path = append(path, 1, x, y, 1)
		}
	}
	p.Stroke(path, f.LineWidth)
}

// FloatEnvelope is a filled area.
// X and Y must be a closed path.
type FloatEnvelope struct {
	X, Y []float64
	Z    int
	CoordinateSystem
	LineWidth int
}

func (f FloatEnvelope) Draw(p *Painter) {
	var path raster.Path
	bounds := rect26_6(p.im.Bounds())
	z := fixed.I(f.Z)
	for i := range f.X {
		x, y := transform(f.X[i], f.Y[i], f.CoordinateSystem, bounds)
		x += z
		y -= z
		if i == 0 {
			path = append(path, 0, x, y, 0)
		} else {
			path = append(path, 1, x, y, 1)
		}
	}
	p.Fill(path)
	p.Stroke(path, f.LineWidth)
}

// FloatText is text with floating point coordinates.
type FloatText struct {
	X, Y       float64
	Z          int
	S          string
	Xoff, Yoff int // Additional offset in pixel coordinates.
	Align      int
	CoordinateSystem
	Rect image.Rectangle // Rectangle of the coordinate system.
}

func (f FloatText) toText(x0, y0 int) Text {
	x, y := transform(f.X, f.Y, f.CoordinateSystem, rect26_6(f.Rect))
	x += fixed.Int26_6(f.Xoff*64) + fixed.I(f.Z) - fixed.I(x0) //x0,y0 is already applied by transform, but again in Text.Draw
	y += fixed.Int26_6(f.Yoff*64) - fixed.I(f.Z) - fixed.I(y0)
	t := Text{X: int(x / 64), Y: int(y / 64), S: f.S, Align: f.Align}
	return t
}
func (f FloatText) Draw(p *Painter) { f.toText(p.x0, p.y0).Draw(p) }
func (p *Painter) FloatTextExtent(f FloatText) (int, int, int, int) {
	return p.Extent(f.toText(p.x0, p.y0))
}

// FloatCircles are many circles given in float point coordinates.
type FloatCircles struct {
	X, Y []float64 // center coordinates
	Z    int
	CoordinateSystem
	Radius    int // circle radius in pixels
	LineWidth int
	Fill      bool
}

func (f FloatCircles) Draw(p *Painter) {
	for i := range f.X {
		x, y := transform(f.X[i], f.Y[i], f.CoordinateSystem, rect26_6(p.im.Bounds()))
		x += fixed.I(f.Z)
		y -= fixed.I(f.Z)
		path := circle{x, y, fixed.Int26_6(f.Radius * 64)}.getPath()
		if f.Fill {
			p.Fill(path)
		} else {
			p.Stroke(path, f.LineWidth)
		}
	}
}

// FloatBars draws rectangles for each 2 float points.
type FloatBars struct {
	X, Y []float64
	CoordinateSystem
	LineWidth int
	Fill      bool
}

func (f FloatBars) Draw(p *Painter) {
	for i := 0; i < len(f.X); i += 2 {
		x0, y0 := transform(f.X[i], f.Y[i], f.CoordinateSystem, rect26_6(p.im.Bounds()))
		x1, y1 := transform(f.X[i+1], f.Y[i+1], f.CoordinateSystem, rect26_6(p.im.Bounds()))
		path := raster.Path{0, x0, y0, 0, 1, x0, y1, 1, 1, x1, y1, 1, 1, x1, y0, 1, 1, x0, y0, 1}
		if f.Fill {
			p.Fill(path)
		}
		if f.LineWidth > 0 {
			p.Stroke(path, f.LineWidth)
		}
	}
}

// Pixel calculates the pixel position for a given FloatPoint.
func (cs CoordinateSystem) Pixel(x, y float64, r image.Rectangle) (X, Y int) {
	a, b := transform(x, y, cs, rect26_6(r))
	return int(a / 64), int(b / 64)
}

// Point calculates FloatPoint coordinates from a pixel position.
func (cs CoordinateSystem) Point(X, Y int, r image.Rectangle) (x, y float64) {
	x = xmath.Scale(float64(X), float64(r.Min.X)-0.5, float64(r.Max.X)-0.5, cs.X0, cs.X1)
	y = xmath.Scale(float64(Y), float64(r.Min.Y)-0.5, float64(r.Max.Y)-0.5, cs.Y0, cs.Y1)
	return
}

func transform(x, y float64, cs CoordinateSystem, bounds fixed.Rectangle26_6) (X, Y fixed.Int26_6) {
	fixedFloat := func(f float64) fixed.Int26_6 {
		return fixed.Int26_6(int(64.0 * f))
	}
	hp := fixed.Int26_6(32)
	x0, x1 := float64(bounds.Min.X+hp)/64.0, float64(bounds.Max.X+hp)/64.0
	y0, y1 := float64(bounds.Min.Y+hp)/64.0, float64(bounds.Max.Y+hp)/64.0
	return fixedFloat(clip(xmath.Scale(x, cs.X0, cs.X1, x0, x1))), fixedFloat(clip(xmath.Scale(y, cs.Y0, cs.Y1, y0, y1)))
}
func transformDirection(dx, dy float64, cs CoordinateSystem, bounds fixed.Rectangle26_6) (float64, float64) {
	hp := fixed.Int26_6(32)
	x0, x1 := float64(bounds.Min.X+hp)/64.0, float64(bounds.Max.X-hp)/64.0
	y0, y1 := float64(bounds.Min.Y+hp)/64.0, float64(bounds.Max.Y-hp)/64.0
	dx *= (x1 - x0) / (cs.X1 - cs.X0)
	dy *= (y1 - y0) / (cs.Y1 - cs.Y0)
	l := math.Hypot(dx, dy)
	return dx / l, dy / l // pixel units as float64
}

// rect26_6 transforms an image.Rectangle to a fixed.Rectangle26_6
func rect26_6(r image.Rectangle) fixed.Rectangle26_6 {
	return fixed.R(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)
}

// clip sets big values to a max value.
// The max value will be the max value when converted to a Int26_6.
// This should give the impression, as if the point is at infinity and lines
// drawn with these points should have the right angles and be clipped.
func clip(f float64) float64 {
	minInt := -33554432
	maxInt := 33554431
	if f*64.0 > float64(maxInt) {
		return float64(maxInt) / 64.0
	}
	if f*64.0 < float64(minInt) {
		return float64(minInt) / 64.0
	}
	return f
}
