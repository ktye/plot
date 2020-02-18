package plot

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// AxisLimits returns the Limits for the plot number.
func AxisLimits(a []Limits, plotNumber int) Limits {
	if len(a) > plotNumber {
		return a[plotNumber]
	} else if len(a) > 0 {
		return a[0]
	}
	return Limits{}
}

func AxisFromString(w string) (ax Limits, err error) {
	if w == "" {
		return Limits{}, nil
	}
	if w == "equal" {
		return Limits{Equal: true}, nil
	}
	a := strings.Split(w, ",")
	if len(a) > 6 {
		return ax, fmt.Errorf("axis string has too many entries")
	}
	var v [6]float64
	for k, b := range a {
		if f, err := strconv.ParseFloat(b, 64); err != nil {
			return ax, fmt.Errorf("wrong axis value (%s): %s", b, err)
		} else {
			v[k] = f
		}
	}
	if len(a) == 1 {
		ax.Ymax = v[0]
	} else if len(a) == 2 {
		ax.Xmin = v[0]
		ax.Xmax = v[1]
	} else if len(a) == 3 {
		ax.Xmin = v[0]
		ax.Xmax = v[1]
		ax.Ymax = v[2]
	} else if len(a) == 4 {
		ax.Xmin = v[0]
		ax.Xmax = v[1]
		ax.Ymin = v[2]
		ax.Ymax = v[3]
	} else if len(a) == 5 {
		ax.Xmin = v[0]
		ax.Xmax = v[1]
		ax.Ymin = v[2]
		ax.Ymax = v[3]
		ax.Zmin = v[4] - 100.0
		ax.Zmax = v[4]
	} else if len(a) == 6 {
		ax.Xmin = v[0]
		ax.Xmax = v[1]
		ax.Ymin = v[2]
		ax.Ymax = v[3]
		ax.Zmin = v[4]
		ax.Zmax = v[5]
	}
	return ax, nil
}
func AxisFromStrings(v []string) (ax []Limits, err error) {
	if len(v) == 0 {
		return ax, nil
	}
	ax = make([]Limits, len(v))
	for i, w := range v {
		ax[i], err = AxisFromString(w)
		if err != nil {
			return ax, err
		}
	}
	return ax, nil
}
func AxisFromEnv(plts Plots) Plots {
	flt := func(s string) float64 { f, _ := strconv.ParseFloat(s, 64); return f }
	xmin, xmax := flt(os.Getenv("xmin")), flt(os.Getenv("xmax"))
	ymin, ymax := flt(os.Getenv("ymin")), flt(os.Getenv("ymax"))
	zmin, zmax := flt(os.Getenv("zmin")), flt(os.Getenv("zmax"))
	for i, p := range plts {
		if p.Limits.Xmin == 0 && p.Limits.Xmax == 0 {
			plts[i].Limits.Xmin, plts[i].Limits.Xmax = xmin, xmax
		}
		if p.Limits.Ymin == 0 && p.Limits.Ymax == 0 {
			plts[i].Limits.Ymin, plts[i].Limits.Ymax = ymin, ymax
		}
		if p.Limits.Zmin == 0 && p.Limits.Zmax == 0 {
			plts[i].Limits.Zmin, plts[i].Limits.Zmax = zmin, zmax
		}
	}
	return plts
}
