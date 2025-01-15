package plot

import (
	"strconv"
)

type Tics struct {
	Pos    []float64 // tic position.
	Labels []string  // tic string.
}

// ticstring formats the number in a compact form.
func ticstring(x float64, prec int) string {
	// If tics are integers between (-10000, 10000) format with %d.
	if x > -10000 && x < 10000 && float64(int(x)) == x {
		return strconv.Itoa(int(x))
	}
	// TODO a more compact form could be possible.
	return strconv.FormatFloat(x, 'g', prec, 64)
}

// getXTics calculates the tic positions and tic labels for the x axis.
// don't care, if they fit or not, they will be reduced later.
func getXTics(lim Limits) Tics {
	a := autoscale{min: lim.Xmin, max: lim.Xmax, isInit: true}
	return a.getTics()
}
func getYTics(lim Limits) Tics {
	a := autoscale{min: lim.Ymin, max: lim.Ymax, isInit: true}
	return a.getTics()
}
func getZTics(lim Limits) Tics {
	a := autoscale{min: lim.Zmin, max: lim.Zmax, isInit: true}
	return a.getTics()
}
