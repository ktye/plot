// Package plotui provides a windows gui frontend to plot.Plots
package plotui

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/ktye/plot"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

type Plot struct {
	ReplotEnvelope func(int, *plot.Plots) // called after zoom
	//UnitDialog     func(int)              // callback after unit is clicked
	//AxisDialog     func(int, plot.Limits) // callback after axis limits are clicked
	Columns     int
	canvas      *walk.CustomWidget // plot canvas
	bitmap      *walk.Bitmap       // underlying plot bitmap
	table       *walk.TableView    // caption table
	slider      *walk.Slider
	model       CaptionModel
	slidePoints int // slide over lines (value 0) or points (value > 0, enabled by a pointclick)
	lineOffset  int
	ignore      bool
	caption     *plot.Caption
	plots       *plot.Plots
	iplots      []plot.IPlotter
	hi          []plot.HighlightID
	mouse       mouseState
	ttf         []byte
}

// SetPlot sets new plots and update the plot.
func (ui *Plot) SetPlot(p plot.Plots, hi []plot.HighlightID) error {
	ui.hi = hi
	ui.plots = &p
	if caption, err := p.MergedCaption(); err != nil {
		return err
	} else {
		ui.caption = &caption
		if ui.table != nil {
			ui.SetCaption()
		}
		ui.slidePoints = 0
		ui.SetSlider(0)
	}
	return ui.setImage()
}

func (ui *Plot) SetFont(ttf []byte) {
	ui.ttf = ttf
}

func (ui *Plot) SetFontSizes(large, small int) error {
	if ui.ttf == nil {
		return fmt.Errorf("ttf is unset")
	}
	font, err := truetype.Parse(ui.ttf)
	if err != nil {
		return err
	}
	f1 := truetype.NewFace(font, &truetype.Options{Size: float64(large), DPI: 72})
	f2 := truetype.NewFace(font, &truetype.Options{Size: float64(small), DPI: 72})
	plot.SetFonts(f1, f2)
	return nil
}

// BuildPlot returns a declarative CustomWidget for the plot image.
func (ui *Plot) BuildPlot(menu []declarative.MenuItem) declarative.CustomWidget {
	var timer *time.Timer
	resizeFunc := func() {
		if ui.canvas != nil && ui.canvas.Parent().Visible() {
			ui.setImage()
		}
		timer = nil
	}
	resize := func() {
		if timer == nil {
			timer = time.AfterFunc(50*time.Millisecond, resizeFunc)
		}
		timer.Reset(50 * time.Millisecond)
	}
	return declarative.CustomWidget{
		StretchFactor:    3,
		AssignTo:         &ui.canvas,
		MinSize:          declarative.Size{Width: 300, Height: 300},
		OnMouseDown:      ui.mouseDown,
		OnMouseUp:        ui.mouseUp,
		OnMouseMove:      ui.mouseMove,
		ContextMenuItems: menu,
		OnSizeChanged:    resize,
		Paint:            ui.paint,
		PaintMode:        declarative.PaintBuffered,
	}
}

// BuildCaption returns a declarative TableView widget for the caption.
func (ui *Plot) BuildCaption(menu []declarative.MenuItem) declarative.TableView {
	return declarative.TableView{
		AssignTo:                 &ui.table,
		Font:                     declarative.Font{Family: "Consolas", PointSize: 12},
		Columns:                  []declarative.TableViewColumn{declarative.TableViewColumn{}},
		Model:                    &ui.model,
		MultiSelection:           true,
		OnSelectedIndexesChanged: ui.lineClicked,
		//NoColumnHeader:           true,
		HeaderHidden: true,
		//NotSortableByHeaderClick: true,
		ContextMenuItems: menu,
	}
}

// BuildSlider connects a Slider to the plotui.
func (ui *Plot) BuildSlider() declarative.Slider {
	return declarative.Slider{
		AssignTo:       &ui.slider,
		Tracking:       true,
		OnValueChanged: ui.sliderChanged,
		Orientation:    declarative.Horizontal,
		Enabled:        false,
	}
}

// SetPaintMode sets the paint mode for the plot canvas.
func (ui *Plot) SetPaintMode(mode walk.PaintMode) {
	ui.canvas.SetPaintMode(mode)
}

// GetPlots returns the plots.
func (ui *Plot) GetPlots() *plot.Plots {
	return ui.plots
}

// GetCaption returns the merged caption.
func (ui *Plot) GetCaption() *plot.Caption {
	return ui.caption
}

// GetCaptionMenu returns the context menu of the caption table.
func (ui *Plot) GetCaptionMenu() *walk.Menu {
	if ui.table == nil {
		return nil
	}
	return ui.table.ContextMenu()
}

// GetLimits returns the current axis limits.
func (ui *Plot) GetLimits() []plot.Limits {
	return plot.LimitsIPlotters(ui.iplots)
}

// GetHighlights returns the lighlighted IDs.
func (ui *Plot) GetHighlights() []plot.HighlightID {
	return ui.hi
}

// GetImage generates the current image in the requested size.
func (ui *Plot) GetImage(width, height int) (image.Image, error) {
	im, _, err := ui.image(width, height)
	return im, err
}

// SetImage creates an image from the current plot and puts it on the canvas.
func (ui *Plot) setImage() error {
	if ui.plots == nil {
		return nil
	}
	if ui.canvas == nil { // headless mode
		return nil
	}
	bounds := ui.canvas.ClientBounds()
	if im, ip, err := ui.image(bounds.Width, bounds.Height); err != nil {
		return err
	} else {
		ui.iplots = ip
		if len(ui.iplots) == 0 {
			return nil
		}
		if bm, err := walk.NewBitmapFromImage(im); err != nil {
			return err
		} else {
			old := ui.bitmap
			ui.bitmap = bm
			ui.canvas.Invalidate()
			if old != nil {
				old.Dispose()
			}
		}
	}
	return nil
}

func (ui *Plot) image(width, height int) (image.Image, []plot.IPlotter, error) {
	if hp, err := ui.plots.IPlots(width, height, ui.Columns); err != nil {
		return nil, nil, err
	} else {
		if im := plot.Image(hp, ui.hi, width, height, ui.Columns); im == nil {
			return nil, nil, fmt.Errorf("could not make image (area too small?)")
		} else {
			return im, hp, nil
		}
	}
}

// paint ist the paint function for the plot canvas.
func (ui *Plot) paint(canvas *walk.Canvas, updateBounds walk.Rectangle) error {
	if ui.bitmap != nil {
		canvas.DrawImage(ui.bitmap, walk.Point{0, 0})
	} else {
		white := walk.RGB(255, 255, 255)
		brush, err := walk.NewSolidColorBrush(white)
		if err != nil {
			return nil
		}
		canvas.FillRectangle(brush, updateBounds)
		brush.Dispose()
	}
	if ui.mouse.rect.Width != 0 || ui.mouse.rect.Height != 0 {
		// Draw zoom rectangle.
		red := walk.RGB(255, 0, 0)
		pen, err := walk.NewCosmeticPen(walk.PenSolid, red)
		if err != nil {
			return nil
		}
		if ui.mouse.modifier == walk.ModShift || ui.mouse.modifier == walk.ModAlt {
			canvas.DrawLine(pen, walk.Point{ui.mouse.x, ui.mouse.y}, walk.Point{ui.mouse.xL, ui.mouse.yL})
		} else {
			canvas.DrawRectangle(pen, ui.mouse.rect)
		}
		pen.Dispose()
	}
	return nil
}

// Update updates the plot from the currentIPlots, it is used after zooming or line clicks.
func (ui *Plot) update(hiIDs []plot.HighlightID) {
	if len(ui.iplots) == 0 {
		return
	}
	bounds := ui.canvas.ClientBounds()
	if im := plot.Image(ui.iplots, hiIDs, bounds.Width, bounds.Height, ui.Columns); im == nil {
		log.Println("could not make image: ", hiIDs)
	} else {
		if bm, err := walk.NewBitmapFromImage(im); err != nil {
			log.Println(err)
		} else {
			old := ui.bitmap
			ui.bitmap = bm
			ui.canvas.Invalidate()
			if old != nil {
				old.Dispose()
			}
		}
	}
}

// ResetZoom removes custom axis settings (even initial ones) from all plots and redraws.
func (ui *Plot) ResetZoom() {
	plts := *ui.plots
	if len(plts) == 0 {
		return
	}
	for i, p := range plts {
		plts[i].Limits = plot.Limits{p.Equal, 0, 0, 0, 0, 0, 0}
	}
	ui.setImage()
}
