package plot

import (
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/ktye/plot/raster"
)

// This file provides a basic widget implementation for github.com/ktye/ui.
// See github.com/ktye/ui/examples/plot.

type UI struct {
	UnitDialog func(int) int
	AxisDialog func(int, Limits) int
	LogMeasure func(complex128)
	p          Plots
	hi         []HighlightID
	iplots     []IPlotter
	rect       image.Rectangle
	save       *image.RGBA
	mouse      mouseState
	invalid    bool
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
	// possible states:
	//	- another widget requested a forced redraw, but UI known nothing has changed
	//		→ redraw from saved image
	//	- during zoom
	//		→ draw the saved image and a red rectange on top
	//	- zoom/pan/click event is finished (invalid==true)
	//		→ redraw stored iplots
	//	- dst has been resized or initial call or plot data has changed (u.rect does not match)
	//		→ plot must be recalculated completely

	zooming := u.mouse.rect != image.ZR
	if (force || zooming) && u.save != nil && u.save.Rect == dst.Rect {
		// Forced draw, but unchanged plot.
		draw.Draw(dst, dst.Rect, u.save, dst.Rect.Min, draw.Src)
		if r := u.mouse.rect; r != image.ZR { // during zoom
			r = r.Add(dst.Rect.Min)
			raster.Rectangle(raster.Image{dst, color.RGBA{R: 0xFF, A: 0xFF}}, r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)
		}
		return
	}
	if u.rect != dst.Rect || u.invalid {
		size := dst.Rect.Size()
		if u.rect != dst.Rect {
			u.rect = dst.Rect
			h, err := u.p.IPlots(size.X, size.Y)
			if err != nil {
				return
			}
			u.iplots = h
		}
		u.save = Image(u.iplots, u.hi, size.X, size.Y).(*image.RGBA)
		if u.save == nil {
			u.save = image.NewRGBA(dst.Rect)
		} else {
			u.save.Rect = u.save.Rect.Add(dst.Rect.Min)
		}
		draw.Draw(dst, dst.Rect, u.save, dst.Rect.Min, draw.Src)
	}
}
func (u *UI) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	switch {
	case but == -1 || but == -2:
		return u.wheel(2*but + 3) // -1→1 (up), -2→-1 (down)
	case dir > 0:
		return u.mouseDown(pos.Sub(u.rect.Min), but, mod)
	case dir < 0:
		return u.mouseUp(pos.Sub(u.rect.Min), but, mod)
	default:
		return u.mouseMove(pos.Sub(u.rect.Min), mod)
	}
}
func (u *UI) Key(r rune, code uint32, dir int, mod uint32) int { return 0 }

// all coordinates in mouse state are relative to the top left corner (not the global destination image).
type mouseState struct {
	pos  image.Point // mouse down position
	cur  image.Point // current position
	but  int         // mouse down state
	mod  uint32
	rect image.Rectangle // red rectangle during zooming
	time time.Time
}

func (u *UI) toHorOrVer(pos image.Point) image.Point {
	d := pos.Sub(u.mouse.pos)
	if x2, y2 := d.X*d.X, d.Y*d.Y; x2 > y2 {
		return image.Point{pos.X, u.mouse.pos.Y}
	}
	return image.Point{u.mouse.pos.X, pos.Y}
}
func (u *UI) mouseDown(pos image.Point, but int, mod uint32) int {
	u.mouse.pos = pos
	u.mouse.but = but
	u.mouse.mod = mod
	return 0
}
func (u *UI) mouseMove(pos image.Point, mod uint32) int {
	if u.mouse.but != 1 {
		return 0
	}
	if u.mouse.mod == 4 { // alt
		pos = u.toHorOrVer(pos)
	}
	u.mouse.cur = pos
	u.mouse.rect = image.Rectangle{u.mouse.pos, pos}
	return 1
}
func (u *UI) mouseUp(pos image.Point, but int, mod uint32) int {
	u.mouse.rect = image.ZR
	u.mouse.but = 0
	elapsed := time.Since(u.mouse.time)
	u.mouse.time = time.Now()
	dx := pos.X - u.mouse.pos.X
	dy := pos.Y - u.mouse.pos.Y
	bounds := u.rect
	if dx*dx+dy*dy > 100 {
		// Click and Move (zoom, pan or draw line)
		if u.mouse.mod == 1 || u.mouse.mod == 4 { // shift|alt
			if u.mouse.mod == 4 { // alt → draw line (aligned)
				pos = u.toHorOrVer(pos)
			}
			if vec, ok := LineIPlotters(u.iplots, u.mouse.pos.X, u.mouse.pos.Y, pos.X, pos.Y, bounds.Dx(), bounds.Dy()); ok {
				if u.LogMeasure != nil {
					u.LogMeasure(vec)
				}
				return u.invalidate()
			}
			u.mouse.rect = image.ZR
		} else {
			// Zoom or pan.
			X := u.mouse.pos.X
			Y := u.mouse.pos.Y
			if but == 1 {
				// Left Click and move: zoom
				if dx < 0 {
					X = pos.X
					dx = -dx
				}
				if dy < 0 {
					Y = pos.Y
					dy = -dy
				}
				if ok, n := ZoomIPlotters(u.iplots, X, Y, dx, dy, bounds.Dx(), bounds.Dy()); ok {
					_ = n
					/* TODO envelope
					if ui.ReplotEnvelope != nil {
						ui.ReplotEnvelope(n, ui.plots)
					}
					*/
					return u.invalidate()
				}
			} else if but == 3 {
				// Right Click and move: pan
				if ok, n := PanIPlotters(u.iplots, X, Y, dx, dy, bounds.Dx(), bounds.Dy()); ok {
					_ = n
					/* TODO envelope
					if ui.ReplotEnvelope != nil {
					ui.ReplotEnvelope(n, ui.plots)
					}
					*/
					return u.invalidate()
				}
			}
		}
	}
	if elapsed < 300*time.Millisecond {
		// Double click (without moving)
		snapToPoint := true
		if mod == 1 {
			snapToPoint = false
		}
		if cb, ok := ClickIPlotters(u.iplots, pos.X, pos.Y, bounds.Dx(), bounds.Dy(), snapToPoint); ok {
			if cb.Type == PointInfoCallback {
				pointInfo := cb.PointInfo
				hi := []HighlightID{HighlightID{
					Line:   pointInfo.LineID,
					Point:  pointInfo.PointNumber,
					XImage: pointInfo.X,
					YImage: pointInfo.Y,
				}}
				println("LineID", hi[0].Line, "Point", hi[0].Point)
				u.hi = hi
				// TODO u.highlightCaption([]int{pointInfo.LineID})
				// TODO u.slidePoints = pointInfo.NumPoints
				// TODO u.SetSlider(pointInfo.PointNumber + 1)
				if pointInfo.IsImage {
					println("pointInfo.String()") // TODO
				}
				return u.invalidate()
			} else if cb.Type == MeasurePoint { // shift-double click: add a point, double-click it again to see values
				return u.invalidate()
			} else if cb.Type == UnitCallback {
				if u.UnitDialog != nil {
					u.UnitDialog(cb.PlotIndex)
				}
			} else if cb.Type == AxisCallback {
				if u.AxisDialog != nil {
					u.AxisDialog(cb.PlotIndex, cb.Limits)
				}
			}
		}
	}
	return 0
}
func (u *UI) wheel(add int) int { // slide through points with the wheel, after double clicking a point in a line.
	if len(u.hi) == 1 {
		u.hi[0].Point += add
		if u.hi[0].Point < 0 {
			u.hi[0].Point = 0
		} // cannot clip max point.
		return u.invalidate()
	}
	return 0
}
func (u *UI) invalidate() int { // plot needs update from iplots due to interactive action
	u.invalid = true
	return 1
}
