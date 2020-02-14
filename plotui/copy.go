package plotui

import (
	"fmt"
	"math"
	"math/cmplx"
	"strings"

	"github.com/ktye/plot"
	"github.com/lxn/walk"
)

func (ui *Plot) CopyLineData() error {
	if ui.plots == nil {
		return fmt.Errorf("there are no plots")
	}
	if len(ui.hi) != 1 {
		return fmt.Errorf("one line must be selected")
	}
	lineNumber := ui.hi[0].Line
	mklinedata := func(l plot.Line) ([]string, error) {
		if len(l.X) == 0 {
			return nil, fmt.Errorf("line contains no X values")
		}
		if len(l.Y) > 0 {

		}
		var v []string
		for i := range l.X {
			s := fmt.Sprintf("%v", l.X[i])
			if len(l.Y) > 0 {
				add := "\tNaN"
				if len(l.Y) > i {
					add = fmt.Sprintf("\t%v", l.Y[i])
				}
				s += add
			}
			if len(l.C) > 0 {
				add := "\tNaN\tNaN"
				if len(l.C) > i {
					amp := cmplx.Abs(l.C[i])
					ang := cmplx.Phase(l.C[i]) / math.Pi * 180.0
					if ang < 0 {
						ang += 360.0
					}
					add = fmt.Sprintf("\t%v\t%v", amp, ang)
				}
				s += add
			}
			v = append(v, s)
		}
		return v, nil
	}
	var lineData []string
	for i, p := range *ui.plots {
		if lineNumber < 0 || lineNumber >= len(p.Lines) {
			return fmt.Errorf("plot %d has not line number %d", i+1, lineNumber+1)
		} else {
			d, err := mklinedata(p.Lines[lineNumber])
			if err != nil {
				return err
			}
			if i == 0 {
				lineData = d
			} else {
				if nl := len(d); nl != len(lineData) {
					return fmt.Errorf("plot %d has different number of values", i+1)
				}
				for k := range lineData {
					lineData[k] += "\t" + d[k]
				}
			}
		}
	}
	s := strings.Join(lineData, "\n")
	cb := walk.Clipboard()
	if err := cb.Clear(); err != nil {
		return err
	}
	if err := cb.SetText(s); err != nil {
		return err
	}
	return nil
}
