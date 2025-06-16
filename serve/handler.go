package serve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
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
	hp      plot.Iplots
	hi      []plot.HighlightID
	pi      plot.PointInfo
	e       error
	w, h, c int
}

var p t

func SetPlots(plots plot.Plots, e error) {
	p.Lock()
	defer p.Unlock()
	p.p = plots
	p.hi = nil
	p.e = e
	p.w, p.h, p.c = 0, 0, 0
}

type Option struct {
	Class    string `json:"c"`
	Color    string `json:"rgb"`
	Text     string `json:"t"`
	Selected bool   `json:"s"`
	Number   string `json:"n"`
	Line     int    `json:"l"`
}

func setSize(width, height, columns int) {
	p.w, p.h, p.c = width, height, columns
	f1, f2 := plot.Fonts()
	p.hp, _ = p.p.Iplots(vg.NewImage(width, height, f1, f2), columns)
}

func Plot(w http.ResponseWriter, r *http.Request) {
	p.Lock()
	defer p.Unlock()
	q := r.URL.Query()

	wi, hi := atoi(q.Get("w")), atoi(q.Get("h"))
	hl := atois(q.Get("hl"))
	x, y := atoi(q.Get("x")), atoi(q.Get("y"))
	z := atois(q.Get("z"))
	//fmt.Println("w", wi, "h", hi, "z", z, "x", x, "y", y, "hl", hl)

	if q.Get("caption") != "" {
		writeCaption(w, r)
		return
	}
	if p.e != nil {
		errImage(w, wi, hi, p.e)
		return
	}
	if p.p == nil {
		errImage(w, wi, hi, fmt.Errorf("no plot"))
		return
	}
	if q.Get("gethi") != "" {
		if p.hi == nil {
			p.pi.LineID = -1
		}
		if e := json.NewEncoder(w).Encode(p.pi); e != nil {
			http.Error(w, e.Error(), 400)
		}
		return
	}
	if p.w != wi || p.h != hi {
		setSize(wi, hi, 4)
	}

	if s := q.Get("pt"); s != "" { // slider
		pt := atoi(s)
		if len(p.hi) == 0 {
			p.hi = []plot.HighlightID{plot.HighlightID{Line: 0}}
		}
		for i := range p.hi {
			p.hi[i].Point = pt
		}
	} else if len(hl) == 0 {
		p.hi = nil
	} else { // caption line
		p.hi = make([]plot.HighlightID, len(hl))
		for i, n := range hl {
			p.hi[i] = plot.HighlightID{Line: n - 1, Point: -1}
		}
	}
	if len(z) == 4 {
		if q.Get("draw") != "" {
			v, ok := plot.LineIPlotters(p.hp, z[0], z[1], z[2]+z[0], z[3]+z[1])
			fmt.Println("draw vector", v, ok)
		} else {
			ok, n := plot.ZoomIPlotters(p.hp, z[0], z[1], z[2], z[3])
			fmt.Println("zoom", ok, n)
		}
	}
	if x != 0 && y != 0 {
		snapToPoint := true
		if callback, ok := plot.ClickIPlotters(p.hp, x, y, snapToPoint, false, true); ok {
			if callback.Type == plot.PointInfoCallback {
				p.pi = callback.PointInfo
				p.hi = []plot.HighlightID{plot.HighlightID{
					Line:   p.pi.LineID,
					Point:  p.pi.PointNumber,
					XImage: p.pi.X,
					YImage: p.pi.Y,
				}}
				//ww.Header().Set("l", itoa(pointInfo.LineID))
				//ww.Header().Set("m", itoa(pointInfo.NumPoints-1))
				//ww.Header().Set("p", itoa(pointInfo.PointNumber))
			}
		}
	}

	im := p.hp.Image(p.hi)
	if im == nil {
		errImage(w, wi, hi, fmt.Errorf("plot area is too small"))
		return
	}
	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, im)
}
func writeCaption(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if len(p.p) == 0 || p.p[0].Caption == nil {
		w.Write([]byte("[]"))
		return
	}
	c := p.p[0].Caption

	var buf bytes.Buffer
	lineOffset, e := c.WriteTable(&buf, plot.Numbers)
	if e != nil {
		errImage(w, p.w, p.h, e)
		return
	}
	s := strings.Split(string(buf.Bytes()), "\n")
	if len(s) > 0 && len(s[len(s)-1]) == 0 {
		s = s[:len(s)-1]
	}
	o := make([]Option, len(s))
	for i := 0; i < len(s); i++ {
		o[i].Text = s[i]
		if n := i - lineOffset; n >= 0 {
			o[i].Number = strconv.Itoa(n + 1)
		}
		if co := c.Color(i, lineOffset); co != nil {
			o[i].Color = htmlColor(co)
		}
	}
	if c.Title != "" && len(o) > 0 {
		o[0].Class = "header"
	}
	for i := range o {
		o[i].Text = html.EscapeString(o[i].Text)
	}
	if e := json.NewEncoder(w).Encode(o); e != nil {
		fmt.Println(e)
	}
}
func htmlColor(c color.Color) string {
	r, g, b, a := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x%02x", r>>8, g>>8, b>>8, a>>8)
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

func Plotjs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")
	w.Write([]byte(Js))
}

const Js = `
var plotcnv = document.getElementById("plot-cnv")
var plotsld = document.getElementById("plot-sld")
var plotcap = document.getElementById("plot-cap")
var ctx = plotcnv.getContext("2d")
var zoom = false
var bak
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
 let draw = e.shiftKey ? "&draw=1" : ""
 plot("&z="+zoom+draw)
 zoom = false
}
function clearZoom(){ctx.clearRect(0,0,plotcnv.width,plotcnv.height); ctx.drawImage(bak,0,0) }
function clickPlot(e){
 plot("&x="+e.offsetX+"&y="+e.offsetY, true)
}
function slideWheel(e) {
 plotsld.value = Number(plotsld.value) + ((e.deltaY<0) ? 1 : -1)
 plotSlide()
}
function plotSlide(){ plot("&pt="+plotsld.value) }

function plot(x,sethi){
 x = (x === undefined) ? "" : x
 let r = Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15)
 let url = "plot?w=" + plotcnv.width + "&h=" + plotcnv.height + x + "&r=" + r
 let bg = new Image();
 bg.src = url
 bg.onload = function(){ctx.drawImage(bg,0,0); bak=bg }
 if(sethi===true)get("plot?gethi=1", setHighlightAfterClick)
}
function caption(){get("plot?caption=1",setCaption)}

function deleteAll(e){var c=e.lastElementChild;while(c){e.removeChild(c);c=e.lastElementChild}}
function space(s){return s.replace(/ /g, "&nbsp;")}
function setOptions(e,opts){deleteAll(e);
 for(var i=0;i<opts.length;i++){
  var op=opts[i]
  var o=document.createElement("option")
  if(op.c != "")o.classList.add(op.c)
  if(op.n != "")o.dataset.n = op.n
  //o.style.color=op.rgb
  o.style.background=i?"white":"black"
  o.style.color=i?"black":"white"
  o.style.borderLeft="0.5em solid "+((i==0)?"black":(i==1)?"white":op.rgb)
  o.dataset.l=op.l;e.appendChild(o);o.innerHTML=space(op.t);o.selected=op.s}}
function setCaption(s) { 
 var o=JSON.parse(s)
 if(o.Error){deleteAll(caption);getSet(o.Src);return}
 setOptions(plotcap,o)
 plotcap.onchange=function(){
  Array.from(plotcap.options).forEach((o,i)=>{if(i<2)o.selected=false})
  var s=selectedNumbers(plotcap);plot("&hl="+String(s))}
}
function setHighlightAfterClick(s){
 let h=JSON.parse(s);if(h.LineID<0)return
 setSelectedCaptionLine(1+h.LineID)
 if(h.PointNumber>=0&&h.NumPoints>0)setSlider(h.PointNumber,h.NumPoints)
}
function setSlider(x,n){plotsld.min=1;plotsld.max=n;plotsld.value=x}
function setSelectedCaptionLine(x){
 for(let i=0;i<plotcap.childNodes.length;i++)plotcap.childNodes[i].selected=(plotcap.childNodes[i].dataset.n==x)}
function selectedNumbers(e){ // [1,3,9]
 if(e.selectedIndex<0){return ""}
 let n=[]
 for(var i=e.selectedIndex;i<e.length;i++)if(e[i].selected && e[i].dataset.n)n.push(Number(e[i].dataset.n))
 return n
}
function get(p,f){
 var r=new XMLHttpRequest()
 r.onreadystatechange=function(){if(this.readyState==4&&this.status==200){if(f)f(this.response,this)}}
 r.open("GET",p)
 r.send() }

plotcnv.ondblclick  = clickPlot
plotcnv.onwheel     = slideWheel
plotcnv.onmousedown = zoomStart
plotcnv.onmousemove = zoomMove
plotcnv.onmouseup   = zoomEnd
plotsld.onwheel     = slideWheel
plotsld.onchange    = plotSlide

export { plot, caption }
`

// plotimg plot-img

func Html(w, h int) string {
	return fmt.Sprintf(`<div id="plot-div">
 <canvas id="plot-cnv" width="%d" height="%d"></canvas>
 <input  id="plot-sld" type="range" min="1" max="100" value="50" />
 <select id="plot-cap" multiple></select>
</div>
`, w, h)
}
