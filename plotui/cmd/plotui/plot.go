package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ktye/plot"
	"github.com/ktye/plot/plotui"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

var mw *walk.MainWindow

func main() {
	var r io.Reader
	if len(os.Args) == 1 { // Read plots from file or stdin.
		r = os.Stdin
	} else if len(os.Args) == 2 {
		if f, err := os.Open(os.Args[1]); err != nil {
			fatalw(err)
		} else {
			defer f.Close()
			r = f
		}
	} else {
		fatalw(fmt.Errorf("wrong number of arguments"))
	}
	plots, err := plot.DecodeAny(r)
	if err != nil {
		fatalw(err)
	}

	withCaption := false
	for _, p := range plots {
		if p.Caption != nil {
			withCaption = true
			break
		}
	}
	//panic(withCaption)

	var ui plotui.Plot
	var children []declarative.Widget
	if withCaption {
		children = append(children, declarative.VSplitter{
			Children: []declarative.Widget{
				ui.BuildPlot(Menu(&ui)),
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
				ui.BuildPlot(Menu(&ui)),
				ui.BuildSlider(),
			},
		})
	}

	err = declarative.MainWindow{
		AssignTo:  &mw,
		Title:     "Plot",
		Size:      declarative.Size{800, 800},
		OnKeyDown: keyHandler,
		Layout:    declarative.VBox{MarginsZero: true},
		Children:  children,
	}.Create()
	if err != nil {
		fatalw(err)
	}

	if ico, err := walk.NewIconFromResourceId(11); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		mw.SetIcon(ico)
	}

	ui.SetPlot(plots, nil)
	mw.Run()
}

func keyHandler(key walk.Key) {
	if key == walk.KeyQ {
		os.Exit(0)
	}
}

func fatalw(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)

	declarative.MainWindow{
		Title:  "Plot startup error",
		Size:   declarative.Size{800, 50},
		Layout: declarative.VBox{},
		Children: []declarative.Widget{
			declarative.Label{Text: err.Error()},
		},
	}.Run()
}

func Menu(ui *plotui.Plot) []declarative.MenuItem {
	return []declarative.MenuItem{
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
						log.Print(err)
					} else {
						defer f.Close()
						if err := plots.Encode(f); err != nil {
							log.Print(err)
						}
					}
				}
			},
		},
		declarative.Action{
			Text: "Copy line data",
			OnTriggered: func() {
				if err := ui.CopyLineData(); err != nil {
					log.Print(err)
				}
			},
		},
	}
}
