package plot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ktye/plot/xmath"
)

// HistogramIntervals are the interval settings for a histogram calculation.
type HistogramIntervals struct {
	Min, Max float64
	N        int
}

// ParseHistogramIntervals parses a string with interval settings.
func ParseHistogramIntervals(s string) (HistogramIntervals, error) {
	var h HistogramIntervals
	var parseError = fmt.Errorf("cannot parse interval settings: %s", s)

	if s == "" {
		return h, nil
	}
	v := strings.Split(s, ",")
	if len(v) == 1 {
		if n, err := strconv.Atoi(v[0]); err != nil {
			return h, parseError
		} else {
			h.N = n
			return h, nil
		}
	} else if len(v) == 2 || len(v) == 3 {
		if min, err := strconv.ParseFloat(v[0], 64); err != nil {
			return h, parseError
		} else {
			h.Min = min
		}
		if max, err := strconv.ParseFloat(v[1], 64); err != nil {
			return h, parseError
		} else {
			h.Max = max
		}
		if len(v) == 3 {
			if n, err := strconv.Atoi(v[2]); err != nil {
				return h, parseError
			} else {
				h.N = n
			}
		}
		return h, nil
	}
	return h, parseError
}

// Histogram return coordinates for a histogram plot.
// It is a single continuous line connecting the end points of the bars.
func (intervals HistogramIntervals) Histogram(u []float64, lineIndex, numLines int) (x, y []float64) {
	if len(u) == 0 {
		return nil, nil
	}
	if intervals.Min == intervals.Max {
		intervals.Min, intervals.Max = xmath.MinMax(u)
	}
	if intervals.N < 1 {
		intervals.N = 30
	}

	bins := make([]float64, intervals.N)
	for _, v := range u {
		i := int(xmath.Scale(v, intervals.Min, intervals.Max, 0, float64(intervals.N)) + 0.5)
		if i < 0 {
			i = 0
		}
		if i >= intervals.N {
			i = intervals.N - 1
		}
		bins[i]++
	}
	return Bars(bins, intervals.Min, intervals.Max, intervals.N, lineIndex, numLines)
}

// Bars creates line data for Plot.Type=XY and Line:Style: DataStyle{Marker: MarkerStyle{Marker: Bar, Size: 1}}.
// v has nx values at nominal points in the interval [xmin,xmax].
// bars are plotted side by side (not stacked).
func Bars(v []float64, xmin, xmax float64, nx int, lineIndex, numLines int) (x, y []float64) {
	w := (xmax - xmin) / float64(nx)
	dw := 0.9 * w / float64(numLines)
	for i, b := range v {
		f := xmath.Scale(float64(i), 0, float64(nx-1), xmin, xmax)
		x0 := f - w/2
		x0 += 0.05*w + float64(lineIndex)*dw
		x1 := x0 + dw
		x = append(x, x0, x1)
		y = append(y, 0, float64(b))
	}
	return x, y
}
