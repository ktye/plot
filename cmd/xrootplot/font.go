package main

import (
	"encoding/base64"
	"image"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func dec(w, h int, s string) (r []byte) {
	n := w * h * 96
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	p, k := uint(0), 0
	c := b[0]
	r = make([]byte, n)
	for i := range r {
		if 1&(c>>p) != 0 {
			r[i] = 255
		}
		p++
		if p == 8 {
			p = 0
			k++
			if k < len(b) {
				c = b[k]
			}
		}
	}
	return r
}

type Range struct {
	Low, High rune
	Offset    int
}
type Face struct {
	Advance int
	Width   int
	Height  int
	Ascent  int
	Descent int
	Left    int
	Mask    image.Image
	Ranges  []Range
}

func (f *Face) Close() error                   { return nil }
func (f *Face) Kern(r0, r1 rune) fixed.Int26_6 { return 0 }
func (f *Face) Metrics() font.Metrics {
	return font.Metrics{
		Height:     fixed.I(f.Height),
		Ascent:     fixed.I(f.Ascent),
		Descent:    fixed.I(f.Descent),
		XHeight:    fixed.I(f.Ascent),
		CapHeight:  fixed.I(f.Ascent),
		CaretSlope: image.Point{X: 0, Y: 1},
	}
}
func (f *Face) Glyph(dot fixed.Point26_6, r rune) (
	dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
loop:
	for _, rr := range [2]rune{r, '\ufffd'} {
		for _, rng := range f.Ranges {
			if rr < rng.Low || rng.High <= rr {
				continue
			}
			maskp.Y = (int(rr-rng.Low) + rng.Offset) * (f.Ascent + f.Descent)
			ok = true
			break loop
		}
	}
	if !ok {
		return image.Rectangle{}, nil, image.Point{}, 0, false
	}
	x := int(dot.X+32)>>6 + f.Left
	y := int(dot.Y+32) >> 6
	dr = image.Rectangle{
		Min: image.Point{
			X: x,
			Y: y - f.Ascent,
		},
		Max: image.Point{
			X: x + f.Width,
			Y: y + f.Descent,
		},
	}
	return dr, f.Mask, maskp, fixed.I(f.Advance), true
}
func (f *Face) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	return fixed.R(0, -f.Ascent, f.Width, +f.Descent), fixed.I(f.Advance), true
}
func (f *Face) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	return fixed.I(f.Advance), true
}

var Face10x20 = &Face{
	Advance: 10,
	Width:   10,
	Height:  20,
	Ascent:  20,
	Descent: 0,
	Mask:    mask10x20,
	Ranges: []Range{
		{'\u0020', '\u007f', 0},
		{'\ufffd', '\ufffe', 95},
	},
}

var mask10x20 = &image.Alpha{
	Stride: 10,
	Rect:   image.Rectangle{Max: image.Point{10, 96 * 20}},
	Pix:    dec(10, 20, "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")}
