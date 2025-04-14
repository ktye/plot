package plotui

import (
	"log"
	"time"

	"github.com/ktye/plot"
	"github.com/ktye/plot/xmath"
	"github.com/lxn/walk"
)

type mouseState struct {
	x, y     int
	xL, yL   int
	rect     walk.Rectangle
	button   walk.MouseButton
	time     time.Time
	modifier walk.Modifiers
}

func (ui *Plot) toHorOrVer(x, y int) (int, int) {
	dx := x - ui.mouse.x
	dy := y - ui.mouse.y
	if dx*dx > dy*dy {
		return x, ui.mouse.y
	}
	return ui.mouse.x, y
}

func (ui *Plot) mouseDown(x, y int, button walk.MouseButton) {
	ui.mouse.x = x
	ui.mouse.y = y
	ui.mouse.button = button
	ui.mouse.modifier = walk.ModifiersDown()
}

func (ui *Plot) mouseMove(x, y int, button walk.MouseButton) {
	if button != walk.LeftButton {
		return
	}
	if ui.mouse.modifier == walk.ModShift {
		_, _, n := plot.ClickCoords(ui.iplots, ui.mouse.x, ui.mouse.y)
		if !ui.iplots.IsPolar(n) {
			x, y = ui.toHorOrVer(x, y)
		}
	}
	ui.mouse.xL = x
	ui.mouse.yL = y

	x0, x1 := ui.mouse.x, x
	y0, y1 := ui.mouse.y, y
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	ui.mouse.rect.X = x0
	ui.mouse.rect.Y = y0
	ui.mouse.rect.Width = x1 - x0
	ui.mouse.rect.Height = y1 - y0
	ui.canvas.Invalidate()
}

func (ui *Plot) mouseUp(x, y int, button walk.MouseButton) {
	ui.mouse.rect.Width = 0
	ui.mouse.rect.Height = 0
	elapsed := time.Since(ui.mouse.time)
	dx := x - ui.mouse.x
	dy := y - ui.mouse.y
	//bounds := ui.canvas.ClientBoundsPixels()
	if dx*dx+dy*dy > 100 {
		// Click and Move (zoom, pan or draw line)
		if ui.mouse.modifier == walk.ModShift {
			mi, ok := plot.Measure(ui.iplots, ui.mouse.x, ui.mouse.y, x, y)
			if ok {
				if ui.MainWindow != nil {
					if r, ok := MeasureDialog(ui.MainWindow, mi); ok {
						plot.Annotate(ui.iplots, r.MeasureInfo, r.label, r.circle, r.color, r.linewidth)
						ui.update(nil)
					}
				} else {
					if vec, ok := plot.LineIPlotters(ui.iplots, ui.mouse.x, ui.mouse.y, x, y); ok {
						log.Print("vector: ", xmath.Absang(vec, "%v@%.0f"))
						ui.update(nil)
					}
				}
			}
			ui.mouse.rect = walk.Rectangle{}
		} else {
			// Zoom or pan.
			X := ui.mouse.x
			Y := ui.mouse.y
			if ui.mouse.button == walk.LeftButton {
				// Left Click and move: zoom
				if dx < 0 {
					X = x
					dx = -dx
				}
				if dy < 0 {
					Y = y
					dy = -dy
				}
				if ok, n := plot.ZoomIPlotters(ui.iplots, X, Y, dx, dy); ok {
					ui.update(nil)
					if ui.ReplotEnvelope != nil {
						ui.ReplotEnvelope(n, ui.plots)
					}
				}
			} else if ui.mouse.button == walk.RightButton {
				// Right Click and move: pan
				if ok, n := plot.PanIPlotters(ui.iplots, X, Y, dx, dy); ok {
					ui.update(nil)
					if ui.ReplotEnvelope != nil {
						ui.ReplotEnvelope(n, ui.plots)
					}
				}
			}
		}
	} else {
		if elapsed < 300*time.Millisecond {
			// Double click (without moving)
			snapToPoint, deleteLine := true, false
			if walk.ModifiersDown() == walk.ModShift {
				snapToPoint = false
			}
			//if walk.ModifiersDown() == walk.ModAlt {
			//	// delete line takes the same path as MeasurePoint
			//	deleteLine = true
			//	snapToPoint = false
			//}
			if callback, ok := plot.ClickIPlotters(ui.iplots, x, y, snapToPoint, deleteLine, true); ok {
				if callback.Type == plot.PointInfoCallback {
					pointInfo := callback.PointInfo
					hi := []plot.HighlightID{plot.HighlightID{
						Line:   pointInfo.LineID,
						Point:  pointInfo.PointNumber,
						XImage: pointInfo.X,
						YImage: pointInfo.Y,
					}}
					ui.hi = hi
					ui.update(hi)
					ui.highlightCaption([]int{pointInfo.LineID})
					ui.slidePoints = pointInfo.NumPoints
					ui.SetSlider(pointInfo.PointNumber + 1)
					if pointInfo.IsImage {
						log.Print(pointInfo)
					}
				} else if callback.Type == plot.MeasurePoint {
					defer ui.update(nil)
				} else if callback.Type == plot.UnitCallback {
					//if ui.UnitDialog != nil {
					//	ui.UnitDialog(callback.PlotIndex)
					//}
				} else if callback.Type == plot.AxisCallback {
					//if ui.AxisDialog != nil {
					//	ui.AxisDialog(callback.PlotIndex, callback.Limits)
					//}
				}
			}
		}
	}
	ui.mouse.time = time.Now()
}
