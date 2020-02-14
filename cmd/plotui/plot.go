// Plot data from files or stdin.
//
// Input data from a file (argv) or stdin may be in the format written by (plot.Plots).Encode
// or in the simpler data-only plot.TextDataFormat
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ktye/plot"
	"github.com/ktye/plot/plotui"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

func main() {

	// Read plots from file or stdin.
	var r io.Reader
	if len(os.Args) == 1 {
		r = os.Stdin
	} else if len(os.Args) == 2 {
		if f, err := os.Open(os.Args[1]); err != nil {
			fatalw(err)
		} else {
			defer f.Close()
			r = f
		}
	}

	br := bufio.NewReader(r)
	c, err := br.Peek(1)
	if err != nil {
		fatalw(err)
	}
	var plots plot.Plots
	if c[0] == 'P' {
		plots, err = plot.DecodePlots(br)
	} else {
		plots, err = plot.TextDataPlot(br)
	}
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

	var ui plotui.Plot
	var mw *walk.MainWindow
	var children []declarative.Widget
	var status Status
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
		children = append(children, ui.BuildPlot(Menu(&ui)))
	}

	err = declarative.MainWindow{
		AssignTo:  &mw,
		Title:     "Plot",
		Size:      declarative.Size{800, 800},
		OnKeyDown: keyHandler,
		StatusBarItems: []declarative.StatusBarItem{
			declarative.StatusBarItem{
				AssignTo: &status.sb,
			},
		},
		// TODO ContextMenuItems
		Layout:   declarative.VBox{MarginsZero: true},
		Children: children,
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

	log.SetOutput(status)
	ui.SetPlot(plots, nil)
	mw.Run()
}

type Status struct {
	sb *walk.StatusBarItem
}

// Write implements the io.Writer interface, such that Status can be used as a logger.
func (s Status) Write(b []byte) (int, error) {
	str := string(b)
	s.sb.SetText(str)
	s.sb.SetToolTipText(str)
	return len(b), nil
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
	}
}
