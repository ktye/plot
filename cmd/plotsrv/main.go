package main

// example for plot/serve

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/ktye/plot"
	"github.com/ktye/plot/serve"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":2020", "server address")
	flag.Parse()

	serve.SetPlots(exampleplots(), nil)

	http.HandleFunc("/plot.js", serve.Plotjs)
	http.HandleFunc("/plot", serve.Plot)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(page())) })
	fmt.Println("plotsrv -addr", addr)
	http.ListenAndServe(addr, nil)
}

func page() string {
	w, h := 800, 400
	wpx := "width:800px;"
	wh := wpx + ";height:400px;"
	return strings.Join([]string{
		`<!DOCTYPE html><html><head><meta charset="utf-8"><title>plotsrv</title>`,
		`<style>`,
		`#plot-cnv{position:absolute;right:10px;top:10px; ` + wh + `z-index:1;}`,
		`#plot-img{position:absolute;right:10px;top:10px; ` + wh + `z-index:0;}`,
		`#plot-sld{position:absolute;right:10px;top:410px;` + wpx + `;height:20px}`,
		`#plot-cap{position:absolute;right:10px;top:430px;` + wpx + `;height:200px;font-family:monospace}`,
		`</style>`,
		`</head><body>`,
		serve.Html(w, h),
		`<script type=module>
		 import { plot, caption } from './plot.js'
		 plot();caption();
		</script>`,
		`</body>`,
		`</html>`,
	}, "\n")
}

func exampleplots() plot.Plots {
	x := make([]float64, 100)
	for i := range x {
		x[i] = float64(i) / 99.0
	}
	rnd := rand.NormFloat64
	ln := func(id int, re, im func() float64) plot.Line {
		var y []float64
		var c []complex128
		if im == nil {
			y = make([]float64, len(x))
			for i := range x {
				y[i] = re()
			}
		} else {
			c = make([]complex128, len(x))
			for i := range x {
				c[i] = complex(re(), im())
			}
		}
		return plot.Line{Id: id, X: x, Y: y, C: c}
	}
	c := caption()
	plts := plot.Plots{
		plot.Plot{
			Type: plot.XY,
			Lines: []plot.Line{
				ln(0, rnd, nil),
				ln(1, rnd, nil),
				ln(2, rnd, nil),
			},
			Caption: &c,
		},
		plot.Plot{
			Type: plot.Polar,
			Lines: []plot.Line{
				ln(0, rnd, rnd),
				ln(1, rnd, rnd),
				ln(2, rnd, rnd),
			},
		},
		plot.Plot{
			Type: plot.AmpAng,
			Lines: []plot.Line{
				ln(0, rnd, rnd),
				ln(1, rnd, rnd),
				ln(2, rnd, rnd),
			},
		},
	}
	plts[0].SetCaptionColors()
	//fmt.Println(plts[0].Caption)
	return plts
}
func caption() (c plot.Caption) {
	c.Title = "Example plot"
	c.LeadText = []string{"description 1", "description 2"}
	c.Columns = []plot.CaptionColumn{
		plot.CaptionColumn{
			Class:  "c1",
			Name:   "alpha",
			Unit:   "m/s",
			Format: "%.1f",
			Data:   []float64{1.2, 3.4, 5.5, 6.7},
		},
		plot.CaptionColumn{
			Class: "c1",
			Name:  "beta",
			Unit:  "",
			Data:  []string{"alpha beta", "gamma", "delta epsilon", "omega"},
		},
	}
	return c
}
