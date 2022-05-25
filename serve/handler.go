package serve

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/ktye/plot"
	"github.com/ktye/plot/vg"
)

type t struct {
	sync.Mutex
	p       plot.Plots
	hp      []plot.IPlotter
	hi      []plot.HighlightID
	w, h, c int
}

var p t

func SetPlots(plots plot.Plots) {
	p.Lock()
	defer p.Unlock()
	p.p = plots
	p.hp = nil
	p.hi = nil
	p.w, p.h, p.c = 0, 0, 0
}

func setSize(width, height, columns int) {
	p.w, p.h, p.c = width, height, columns
	p.hp, _ = p.p.IPlots(width, height, columns)
}

func Plot(w http.ResponseWriter, r *http.Request) {
	p.Lock()
	defer p.Unlock()
	q := r.URL.Query()

	wi, hi := atoi(q.Get("w")), atoi(q.Get("h"))
	x, y := atoi(q.Get("x")), atoi(q.Get("y"))
	z := atois(q.Get("z"))
	fmt.Println("w", wi, "h", hi, "z", z, "x", x, "y", y)

	if p.p == nil {
		http.Error(w, "not plot is set", 400)
		return
	}
	if p.hp == nil || p.w != wi || p.h != hi {
		setSize(wi, hi, 4)
	}
	if p.hp == nil {
		http.Error(w, "no plot (bad size?)", 400)
		return
	}
	if len(z) == 4 {
		if q.Get("draw") != "" {
			v, ok := plot.LineIPlotters(p.hp, z[0], z[1], z[2]+z[0], z[3]+z[1], p.w, p.h, p.c)
			fmt.Println("draw vector", v, ok)
		} else {
			ok, n := plot.ZoomIPlotters(p.hp, z[0], z[1], z[2], z[3], p.w, p.h, p.c)
			fmt.Println("zoom", ok, n)
		}
	}
	if x != 0 && y != 0 {
		snapToPoint := true
		if callback, ok := plot.ClickIPlotters(p.hp, x, y, p.w, p.h, p.c, snapToPoint, false); ok {
			if callback.Type == plot.PointInfoCallback {
				pointInfo := callback.PointInfo
				p.hi = []plot.HighlightID{plot.HighlightID{
					Line:   pointInfo.LineID,
					Point:  pointInfo.PointNumber,
					XImage: pointInfo.X,
					YImage: pointInfo.Y,
				}}
				//ww.Header().Set("l", itoa(pointInfo.LineID))
				//ww.Header().Set("m", itoa(pointInfo.NumPoints-1))
				//ww.Header().Set("p", itoa(pointInfo.PointNumber))
			}
		}
	}

	im := plot.Image(p.hp, p.hi, p.w, p.h, p.c)
	if im == nil {
		errImage(w, wi, hi, fmt.Errorf("plot area is too small"))
		return
	}
	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, im)
}

func atoi(s string) int {
	i, e := strconv.Atoi(s)
	if e != nil {
		return 0
	}
	return i
}
func atois(s string) []int {
	if s == "" {
		return nil
	}
	v := strings.Split(s, ",")
	r := make([]int, len(v))
	for i := range v {
		r[i] = atoi(v[i])
	}
	return r
}

func errImage(w http.ResponseWriter, width, height int, err error) {
	w.Header().Set("Content-Type", "image/png")
	m := image.NewRGBA(image.Rectangle{Max: image.Point{width, height}})
	draw.Draw(m, m.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	p := vg.NewPainter(m)
	//p.SetFont(font1)
	p.Add(vg.Text{X: 30, Y: 30, Align: 6, S: err.Error()})
	p.Paint()
	if err := png.Encode(w, m); err != nil {
		log.Print(err)
	}
}

const Js = `
var plotcnv = document.getElementById("plot-cnv")
var plotsld = document.getElementById("plot-sld")
var plotimg = document.getElementById("plot-img")
var ctx = plotcnv.getContext("2d")
var zoom = false
function zoomStart(e){ zoom=[e.offsetX, e.offsetY, 0, 0]; plotcnv.style.cursor="crosshair" }
function zoomMove(e){
 if(zoom!==false){
  zoom = [zoom[0], zoom[1], e.offsetX-zoom[0], e.offsetY-zoom[1]]
  clearZoom();ctx.beginPath();ctx.rect(...zoom);ctx.strokeStyle='red';ctx.stroke()
}}
function zoomEnd(e){
 clearZoom(); plotcnv.style.cursor = ""
 if (Math.abs(zoom[2]) < 5 || Math.abs(zoom[3]) < 5) { zoom=false; return }
 let w = plotcnv.clientWidth; var h = plotcnv.clientHeight
 let p = "hi?&w="+w+"&h="+h+"&z="+zoom
 //console.log("zoom:", p, e.shiftKey)
 let draw = e.shiftKey ? "&draw=1" : ""
 plot("&z="+zoom+draw)
 //plot.src = p
 zoom = false
}
function clearZoom(){ctx.clearRect(0,0,plotcnv.width,plotcnv.height)}
function clickPlot(e){
 plot("&x="+e.offsetX+"&y="+e.offsetY)
}
function slideWheel(e) {
 plotsld.value = Number(plotsld.value) + ((e.deltaY<0) ? 1 : -1)
 plotSlide()
}
function plotSlide(){ console.log("slide:", plotsld.value) }

function plot(x){
 x = (x === undefined) ? "" : x
 console.log("plot", x)
 plotimg.src = "plot?w=" + plotimg.width + "&h=" + plotimg.height + x
}

plotcnv.ondblclick  = clickPlot
plotcnv.onwheel     = slideWheel
plotcnv.onmousedown = zoomStart
plotcnv.onmousemove = zoomMove
plotcnv.onmouseup   = zoomEnd
`

func Html(w, h int) string {
	wh := fmt.Sprintf(`width="%d" height="%d"`, w, h)
	return fmt.Sprintf(`<div id="plot-div">
 <canvas id="plot-cnv" %s></canvas>
 <image  id="plot-img" %s></image>
 <input  id="plot-sld" type="range" min="1" max="100" value="50" />
 <select id="plot-cap" multiple></select>
</div>
`, wh, wh)
}
