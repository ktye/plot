package xmath

// Scale does a linear scale of input x to output y.
func Scale(x, xmin, xmax, ymin, ymax float64) (y float64) {
	return ymin + (x-xmin)*(ymax-ymin)/(xmax-xmin)
}
