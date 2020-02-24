package main

import (
	"fmt"
	"image"
	"time"

	"github.com/ktye/plot"
)

func animate(p plot.Plots) {
	if len(p) < 1 {
		return
	}
	if ani == 1 {
		animateLines(p)
	} else {
		animatePlots(p)
	}
}
func animateLines(p plot.Plots) { // each frame: multiple plots side-by-side, single line each
	n := len(p[0].Lines)
	for i := range p {
		if len(p[i].Lines) != n {
			fatal(fmt.Errorf("plots have different number of lines"))
		}
	}
	pp := make(plot.Plots, len(p))
	for i := range p {
		pp[i] = p[i]
		pp[i].Lines = make([]plot.Line, 1)
	}
	w, h := screensize()
	frames := make([]*image.RGBA, n)
	for f := 0; f < n; f++ {
		for i := range pp {
			pp[i].Lines[0] = p[i].Lines[f]
			pp[i].Lines[0].Id = 0
		}
		frames[f] = frame(pp, w, h)
	}
	loop(frames, w, h)
}
func animatePlots(p plot.Plots) { // each frame: single plot, multiple lines
	n := len(p)
	for i := range p {
		if p[i].Type != p[0].Type {
			fatal(fmt.Errorf("plots have different types"))
		}
	}
	lim, e := p.EqualLimits()
	fatal(e)
	for i := range p {
		p[i].Limits = lim
	}
	w, h := screensize()
	frames := make([]*image.RGBA, n)
	for f := 0; f < n; f++ {
		frames[f] = frame(plot.Plots{p[f]}, w, h)
	}
	loop(frames, w, h)
}
func frame(p plot.Plots, w, h int) *image.RGBA {
	ip, e := p.IPlots(w, h, 0)
	fatal(e)
	return plot.Image(ip, nil, w, h, 0).(*image.RGBA)
}
func loop(imgs []*image.RGBA, w, h int) {
	d := time.Millisecond * time.Duration(fps)
	for {
		for _, m := range imgs {
			draw(w, h, m.Pix)
			if d > 0 {
				time.Sleep(d)
			}
		}
	}
}
