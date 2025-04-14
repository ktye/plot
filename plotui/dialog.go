package plotui

import (
	"fmt"
	"log"
	"math"
	"math/cmplx"
	"strconv"
	"strings"

	"github.com/ktye/plot"
	"github.com/ktye/plot/xmath"
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

type MeasureResult struct {
	plot.MeasureInfo
	label     string
	color     int
	linewidth int
	circle    int //0(arrow) 1(circle/radius) 2(circle/diameter)
}

func MeasureDialog(parent walk.Form, mi plot.MeasureInfo) (MeasureResult, bool) {
	var r MeasureResult
	var dlg *walk.Dialog

	rpm := false
	title := "measure distance"
	distance := "distance"
	maysnap := false
	S := func(z complex128) string {
		if mi.Polar {
			return xmath.Absang(z, "%.4g@%.0f")
		} else if mi.Vertical {
			return fmt.Sprintf("%.4g", imag(z))
		} else {
			return fmt.Sprintf("%.4g", real(z))
		}
	}
	var off float64
	if mi.Polar {
		title, distance = "measure vector", "vector"
		if cmplx.Abs(mi.A-mi.AA.C) > cmplx.Abs(mi.A) { // include origin to snap point
			mi.AA.C = 0
		}
		if cmplx.Abs(mi.B-mi.BB.C) > cmplx.Abs(mi.B) {
			mi.BB.C = 0
		}
		maysnap = cmplx.Abs(mi.AA.C-mi.BB.C) > 0
	} else if mi.Vertical {
		off = real(mi.A)
		rpm = mi.Yunit == "s"
	} else {
		off = imag(mi.A)
		rpm = mi.Xunit == "s"
	}

	parse := func(e *walk.LineEdit, def string) (complex128, bool) {
		s := e.Text()
		if mi.Polar {
			z, err := xmath.ParseComplex(s)
			if err != nil {
				e.SetText(def)
				return 0, false
			}
			return z, true
		} else {
			f, err := strconv.ParseFloat(strings.ReplaceAll(s, ",", "."), 64)
			if err != nil {
				return 0, false
			}
			if mi.Vertical {
				return complex(0, f), true
			} else {
				return complex(f, 0), true
			}
		}
	}
	rpmlabel := func(z complex128) string {
		s := S(z)
		if mi.Polar == false {
			s = strings.TrimPrefix(s, "-")
		}
		if rpm == false {
			return s
		}
		spd := math.Abs(60 / real(z))
		if mi.Vertical {
			spd = math.Abs(60 / imag(z))
		}
		if spd >= 10000 {
			return s + "s (" + fmt.Sprintf("%.4g krpm)", spd/1000)
		} else {
			return s + "s (" + fmt.Sprintf("%.4g rpm)", spd)
		}
	}

	var sta *walk.LineEdit
	var end *walk.LineEdit
	var dif *walk.LineEdit
	var color *walk.ComboBox
	var linewidth *walk.ComboBox
	var label *walk.LineEdit
	var circle *walk.ComboBox

	update := func(a, b, c bool) {
		set := func(x *walk.LineEdit, s string, b bool) {
			if b {
				x.SetText(s)
			}
		}
		set(sta, S(mi.A), a)
		set(end, S(mi.B), b)
		set(dif, S(mi.B-mi.A), c)
		if c {
			label.SetText(rpmlabel(mi.B - mi.A))
		}
	}

	children := []declarative.Widget{
		declarative.Label{Text: "start"},
		declarative.LineEdit{
			AssignTo: &sta,
			Text:     S(mi.A),
			OnEditingFinished: func() {
				if z, o := parse(sta, S(mi.A)); o {
					mi.A = z
					update(false, false, true)
				}
			},
		},
		declarative.Label{Text: "end"},
		declarative.LineEdit{
			AssignTo: &end,
			Text:     S(mi.B),
			OnEditingFinished: func() {
				if z, o := parse(end, S(mi.B)); o {
					mi.B = z
					update(false, false, true)
				}
			},
		},
		declarative.Label{Text: distance},
		declarative.LineEdit{
			AssignTo: &dif,
			Text:     S(mi.B - mi.A),
			OnEditingFinished: func() {
				if z, o := parse(dif, S(mi.B-mi.A)); o {
					mi.B = mi.A + z
					update(true, true, false)
					label.SetText(rpmlabel(mi.B - mi.A))
				}
			},
		},
		declarative.VSpacer{},
		declarative.Label{Text: "label"},
		declarative.LineEdit{
			AssignTo: &label,
			Text:     rpmlabel(mi.B - mi.A),
		},
		declarative.ComboBox{
			AssignTo:     &circle,
			Model:        []string{"arrow", "circle radius", "circle diameter"},
			Enabled:      mi.Polar,
			CurrentIndex: 0,
		},
		declarative.ComboBox{
			AssignTo:     &color,
			Model:        []string{"black:0", "color 1", "color 2", "color 3", "color 4", "color 5", "color 6", "color 7", "color 8", "color 9"},
			CurrentIndex: 0,
		},
		declarative.ComboBox{
			AssignTo:     &linewidth,
			Model:        []string{"1 px", "2 px", "3 px", "4 px"},
			CurrentIndex: 0,
		},
		declarative.VSpacer{},
		declarative.Composite{
			Layout: declarative.HBox{},
			Children: []declarative.Widget{
				declarative.PushButton{
					Text: "Ok",
					OnClicked: func() {
						if mi.Polar {
						} else if mi.Vertical {
							mi.A = complex(off, imag(mi.A))
							mi.B = complex(off, imag(mi.B))
						} else {
							mi.A = complex(real(mi.A), off)
							mi.B = complex(real(mi.B), off)
						}
						r.MeasureInfo = mi
						r.label = label.Text()
						r.color = color.CurrentIndex()
						r.linewidth = linewidth.CurrentIndex() + 1
						r.circle = circle.CurrentIndex()
						dlg.Accept()
					},
				},
				declarative.PushButton{
					Text:    "snap",
					Enabled: maysnap,
					OnClicked: func() {
						mi.A, mi.B = mi.AA.C, mi.BB.C
						update(true, true, true)
					},
				},
			},
		},
	}

	minSize := declarative.Size{300, 100}
	d := declarative.Dialog{
		AssignTo: &dlg,
		Title:    title,
		MinSize:  minSize,
		Layout:   declarative.VBox{},
		Children: children,
	}
	if err := d.Create(parent); err != nil {
		log.Println(err)
	}
	dlg.Run()
	if dlg.Result() == walk.DlgCmdOK {
		return r, true
	}
	return r, false
}
