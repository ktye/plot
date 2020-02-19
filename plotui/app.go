package plotui

import (
	"fmt"
	"os"

	"github.com/ktye/plot"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

// MainWindow returns the plotui widgets within a main window.
func MainWindow(plots plot.Plots) (*walk.MainWindow, *Plot, error) {
	var mw *walk.MainWindow
	withCaption := false
	for _, p := range plots {
		if p.Caption != nil {
			withCaption = true
			break
		}
	}

	var ui Plot
	var children []declarative.Widget
	if withCaption {
		children = append(children, declarative.VSplitter{
			Children: []declarative.Widget{
				ui.BuildPlot(mainMenu(mw, &ui)),
				declarative.Composite{
					Layout: declarative.VBox{MarginsZero: true, SpacingZero: true},
					Children: []declarative.Widget{
						ui.BuildSlider(),
						ui.BuildCaption(nil),
					},
				},
			},
		})
	} else {
		children = append(children, declarative.Composite{
			Layout: declarative.VBox{MarginsZero: true, SpacingZero: true},
			Children: []declarative.Widget{
				ui.BuildPlot(mainMenu(mw, &ui)),
				ui.BuildSlider(),
			},
		})
	}

	err := declarative.MainWindow{
		AssignTo: &mw,
		Title:    "Plot",
		Size:     declarative.Size{800, 800},
		OnKeyDown: func(key walk.Key) {
			if key == walk.KeyQ {
				os.Exit(0)
			}
		},
		Layout:   declarative.VBox{MarginsZero: true},
		Children: children,
	}.Create()
	return mw, &ui, err
}

func mainMenu(mw *walk.MainWindow, ui *Plot) []declarative.MenuItem {
	return []declarative.MenuItem{
		declarative.Action{
			Text:        "Reset zoom",
			OnTriggered: ui.ResetZoom,
		},
		declarative.Action{
			Text:        "Screenshot (to clipboard)",
			OnTriggered: ui.Screenshot,
		},
		declarative.Action{
			Text: "Save plt file",
			OnTriggered: func() {
				plots := ui.GetPlots()
				if plots == nil {
					return
				}
				d := walk.FileDialog{
					Title:    "Save plt file",
					FilePath: "plot.plt",
					Filter:   "Plot files (*.plt)|*.plt||",
				}
				if ok, err := d.ShowSave(mw); ok && err == nil {
					if f, err := os.Create(d.FilePath); err != nil {
						warnDialog(mw, err)
					} else {
						defer f.Close()
						if err := plots.Encode(f); err != nil {
							warnDialog(mw, err)
						}
					}
				}
			},
		},
		declarative.Action{
			Text: "Copy line data",
			OnTriggered: func() {
				if err := ui.CopyLineData(); err != nil {
					warnDialog(mw, err)
				}
			},
		},
		declarative.Action{
			Text: "Caption (to clipboard)",
			OnTriggered: func() {
				if err := ui.ClipboardCaption(); err != nil {
					warnDialog(mw, err)
				}
			},
		},
		declarative.Action{
			Text:        "Measure (help)",
			OnTriggered: func() { warnDialog(mw, fmt.Errorf(MeasureHelp)) },
		},
	}
}

func warnDialog(mw *walk.MainWindow, err error) {
	// walk.MsgBox(mw, "Error", err.Error(), walk.MsgBoxIconWarning) // panic?
	walk.MsgBox(nil, "Error", err.Error(), walk.MsgBoxIconWarning)
}

// Interactive plot: Mouse langugage
// Left click and move: zoom
// Right click and move: pan
// Double-click: mark data point
// Shift and double-click: draw point (no data point)
// Shift, left click and move: draw vector
// Alt, left click and move: draw horizontal/vertical vector
var MeasureHelp string = "Double-click: Mark closest data point\r\nShift + Double-click: Draw point\r\nShift + Click + Mouse move: Draw line (vector)\r\nAlt + Click + Mouse move: Draw line (vector), snap horizontal or vertical"
