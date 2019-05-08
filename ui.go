package plot

import (
	"image"
	"image/draw"
)

// This file provides a basic widget implementation for for github.com/ktye/ui.

type UI struct {
	Plots
	Hi   []HighlightID
	rect image.Rectangle
}

func NewUI(p Plots) *UI {
	var u UI
	u.SetPlots(p)
	return &u
}
func (u *UI) SetPlots(p Plots) {
	u.Plots = p
	u.rect = image.Rectangle{}
	u.Hi = nil
}
func (u *UI) Draw(dst *image.RGBA, force bool) {
	if u.rect != dst.Rect || force {
		u.rect = dst.Rect
		size := u.rect.Size()
		if h, err := u.IPlots(size.X, size.Y); err == nil {
			if im := Image(h, u.Hi, size.X, size.Y); im != nil {
				draw.Draw(dst, u.rect, im, image.ZP, draw.Src)
			}
		}
	}
}
func (u *UI) Mouse(pos image.Point, but int, dir int, mod uint32) int { return 0 }
func (u *UI) Key(r rune, code uint32, dir int, mod uint32) int        { return 0 }
