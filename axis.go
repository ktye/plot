package plot

import (
	"fmt"
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
