package plot

import (
	"image"
	"image/draw"
)

// This file provides a basic widget implementation for for github.com/ktye/ui.

type UI struct {
	p    Plots
	hi   []HighlightID
	rect image.Rectangle
	save *image.RGBA
}

func NewUI(p Plots) *UI {
	var u UI
	u.SetPlots(p)
	return &u
}
func (u *UI) SetPlots(p Plots) {
	u.p = p
	u.hi = nil
	u.rect = image.Rectangle{}
	u.save = nil
}
func (u *UI) Highlight(hi []HighlightID) {
	u.hi = hi
	u.save = nil
}
func (u *UI) Draw(dst *image.RGBA, force bool) {
	if force && u.save != nil && u.save.Rect == dst.Rect {
		// Forced draw, but unchanged plot.
		draw.Draw(dst, dst.Rect, u.save, dst.Rect.Min, draw.Src)
		return
	}
	if u.rect != dst.Rect {
		u.rect = dst.Rect
		size := u.rect.Size()
		if h, err := u.p.IPlots(size.X, size.Y); err == nil {
			u.save = Image(h, u.hi, size.X, size.Y).(*image.RGBA)
			if u.save == nil {
				u.save = image.NewRGBA(dst.Rect)
			} else {
				u.save.Rect = u.save.Rect.Add(dst.Rect.Min)
			}
		}
		draw.Draw(dst, dst.Rect, u.save, dst.Rect.Min, draw.Src)
	}
}
func (u *UI) Mouse(pos image.Point, but int, dir int, mod uint32) int { return 0 }
func (u *UI) Key(r rune, code uint32, dir int, mod uint32) int        { return 0 }
