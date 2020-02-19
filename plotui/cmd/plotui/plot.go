package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ktye/plot"
	"github.com/ktye/plot/plotui"
	"github.com/lxn/walk"
)

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
	mw, ui, err := plotui.MainWindow(plots)
	if err != nil {
		log.Fatal(err)
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

func fatalw(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)

	walk.MsgBox(nil, "Plotui Error", err.Error(), walk.MsgBoxIconError)
	os.Exit(1)
}
