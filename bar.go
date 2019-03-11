package plot

// BarPlotXValues changes the X Values of a plot (which should be of PlotType AmpAng or XY).
// It moves the X Values of all lines to draw the lines of the same bin side by side.
// There are only integer X values each having a width of 1, which is the only case currently needed.
func (p *Plot) BarPlotXValues() {
	// Collect all distince X Values.
	// The map counts the number of lines having the given x value.
	x := make(map[float64]int)
	min := 0.0
	max := 0.0
	for i, l := range p.Lines {
		if len(l.X) > 0 {
			v := l.X[0]

			// This works even for the 1st value.
			n := x[v]
			x[v] = n + 1

			if i == 0 {
				min = v
				max = v
			}
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}
	}

	if min == max {
		return
	}

	// Count max number of lines per bin.
	numPerBin := 0
	for _, n := range x {
		if n > numPerBin {
			numPerBin = n
		}
	}

	// If there is only 1 value per bin, no rearranging must be done.
	if numPerBin < 2 {
		return
	}

	// Currently we support only binDist = 1
	binDist := 1.0

	// To make it "look good", we distribute the data equally, but leave out one spot
	// between two bins:
	//
	// | |   | |
	// -+-----+--  n = 2
	//  1     2
	//
	// ||| |||
	// -+---+--    n = 3
	//  1   2
	//
	// Distance between 2 lines:
	d := binDist / float64(numPerBin+2)
	off := -d * float64(numPerBin-1) / 2.0

	x = make(map[float64]int) // Yes, this overwrites the old map.
	for i, l := range p.Lines {
		if len(l.X) > 0 {
			v := l.X[0]
			n := x[v]
			lineIndex := n
			x[v] = lineIndex + 1
			xpos := v + off + float64(lineIndex)*d
			p.Lines[i].X = []float64{xpos, xpos}
		}
	}
}
