package plot

import "math"

// NiceLimits returns axis limits which include data values in [min, max]
// with nice numbers.
// Reference: Heckbert algorithm (Nice numbers for graph labels)
func NiceLimits(min, max float64) (niceMin, niceMax, spacing float64) {
	maxTicks := 5

	extent := NiceNumber(max-min, false)
	spacing = NiceNumber(extent/float64(maxTicks-1), true)
	niceMin = math.Floor(min/spacing) * spacing
	niceMax = math.Ceil(max/spacing) * spacing

	return niceMin, niceMax, spacing
}

func NiceTics(min, max float64) Tics {
	var t Tics
	extent := NiceNumber(max-min, false)
	niceMin, _, spacing := NiceLimits(min, max)
	pos := niceMin
	n := int(math.Ceil(extent/spacing + 0.5)) // + 1
	for i := 0; i < n; i++ {
		pos = niceMin + float64(i)*spacing
		if pos >= min && pos <= max {
			t.Pos = append(t.Pos, pos)
		}
	}

	// We increase the precision, until all labels are different.
	for prec := 2; prec < 7; prec++ {
		t.Labels = make([]string, len(t.Pos))
		ambiguous := false
		for i := range t.Labels {
			t.Labels[i] = ticstring(t.Pos[i], prec)
			if i > 0 && t.Labels[i] == t.Labels[i-1] {
				ambiguous = true
			}
		}
		if ambiguous == false {
			break
		}
	}
	return t
}

func NiceNumber(extent float64, round bool) float64 {
	exponent := math.Floor(math.Log10(extent))
	fraction := extent / math.Pow10(int(exponent))
	var niceFraction float64
	if round {
		if fraction < 1.5 {
			niceFraction = 1
		} else if fraction < 3 {
			niceFraction = 2
		} else if fraction < 7 {
			niceFraction = 5
		} else {
			niceFraction = 10
		}
	} else {
		if fraction <= 1 {
			niceFraction = 1
		} else if fraction <= 2 {
			niceFraction = 2
		} else if fraction <= 5 {
			niceFraction = 5
		} else {
			niceFraction = 10
		}
	}
	return niceFraction * math.Pow10(int(exponent))
}
