package plot

import (
	"math"

	"github.com/ktye/plot/xmath"
)

// imageValueAt returns the the float value of the Image data of the Line,
// scaled with ImageMax, which is located at the axis coordinates x and y.
func (l Line) imageValueAt(x, y float64) (int, int, float64) {
	i, k := l.imageCoordsToIndexes(x, y)
	if i == -1 || k == -1 {
		return 0, 0, math.NaN()
	}
	if i < 0 || i > len(l.Image) || k < 0 || k > len(l.Image[i]) {
		return 0, 0, math.NaN()
	}
	z := xmath.Scale(float64(l.Image[i][k]), 0.0, 255.0, l.ImageMin, l.ImageMax)
	return i, k, z
}

// imageCoordsToIndexes returns x and y indexes to l.Image
// for the given floating point coordinates x and y within the axes l.X and l.Y.
// It returns -1, -1 for errors.
func (l Line) imageCoordsToIndexes(x, y float64) (int, int) {
	if len(l.X) < 2 || len(l.Y) < 2 || len(l.Image) < 2 || len(l.Image[0]) < 2 {
		return -1, -1
	}

	xmin, xmax := l.X[0], l.X[len(l.X)-1]
	ymin, ymax := l.Y[0], l.Y[len(l.Y)-1]
	if x < xmin || x > xmax || y < ymin || y > ymax {
		return -1, -1
	}

	xdim, ydim := len(l.Image), len(l.Image[0])
	i := int(float64(xdim) * (x - xmin) / (xmax - xmin))
	k := int(float64(ydim) * (y - ymin) / (ymax - ymin))
	return i, k
}
