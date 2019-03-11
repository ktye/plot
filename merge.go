package plot

import "fmt"

type VerticallyMergedPlots struct {
	Plotters []Plotters
	Order    map[[2]int]int
	Colors   []int
}

func (p VerticallyMergedPlots) Plots() (Plots, error) {
	// In a first step, we store all plots in a table.
	var tab [][]Plot // col, row (sic!)
	rows := len(p.Plotters)
	if rows < 1 {
		return nil, fmt.Errorf("cannot merge plots: no input")
	}

	var cols int
	for i, pltrs := range p.Plotters {
		plots, err := pltrs.Plots()
		if err != nil {
			return nil, err
		}
		if i == 0 {
			cols = len(plots)
			tab = make([][]Plot, cols)
		} else {
			if len(plots) != cols {
				return nil, fmt.Errorf("cannot merge plots which are not uniform")
			}
		}
		for k := range plots {
			if i == 0 {
				tab[k] = make([]Plot, rows)
			}
			tab[k][i] = plots[k]
		}
	}

	// Merge each row.
	ret := make(Plots, cols)
	for i := range tab {
		if plt, err := mergeVertical(tab[i], p.Order, p.Colors); err != nil {
			return nil, err
		} else {
			ret[i] = plt
		}
	}
	return ret, nil
}

// MergeVertical merges p vertically into a single Plot.
// The plot must be uniform: same Type, Title and Units.
// It also merges the captions, if the first plot has one.
// In this case, all columns must have identical "form":
// Number, Name and Unit must match.
// The number of Rows has to match the number of Lines.
// The original line order may be stored in the map.
// It maps from the {plotindex, lineindex} to the new line (row) order.
// MergeVertical changes the LineID and their order to reflect the original order.
// If order is nil, lines are concatenated plot by plot.
func mergeVertical(p []Plot, order map[[2]int]int, colors []int) (Plot, error) {
	if len(p) == 0 {
		return Plot{}, fmt.Errorf("merge plots: no input")
	}

	// Check for uniformity.
	typ := p[0].Type
	xlabel, ylabel, title := p[0].Xlabel, p[0].Ylabel, p[0].Title
	xunit, yunit, zunit := p[0].Xunit, p[0].Yunit, p[0].Zunit
	for i, pi := range p {
		if i > 0 {
			if pi.Type != typ {
				return Plot{}, fmt.Errorf("merge plots: different types")
			}
			if pi.Xlabel != xlabel || pi.Ylabel != ylabel || pi.Title != title {
				return Plot{}, fmt.Errorf("merge plots: different labels/titles")
			}
			if pi.Xunit != xunit || pi.Yunit != yunit || pi.Zunit != zunit {
				return Plot{}, fmt.Errorf("merge plots: different units")
			}
		}
	}

	checkCol := func(err error, name string, a, b string) error {
		if err != nil {
			return err
		}
		if a != b {
			return fmt.Errorf("merge plots: caption column (%s): '%s' != '%s'", name, a, b)
		}
		return nil
	}
	if p[0].Caption != nil {
		c0 := p[0].Caption
		for i, pi := range p {
			p[i].Caption.SetEmptyColumns()
			if i > 0 {
				if pi.Caption == nil {
					return Plot{}, fmt.Errorf("merge plots: caption missing")
				}
				if len(pi.Caption.Columns) != len(c0.Columns) {
					return Plot{}, fmt.Errorf("merge plots: different number of caption columns")
				}
				for k, col := range pi.Caption.Columns {
					c0k := c0.Columns[k]
					name, unit := c0k.Name, c0k.Unit
					var err error
					// We do not test column's class, as unbalance columns have a class id
					// speed:unbalance, which could fail.
					err = checkCol(err, "name", col.Name, name)
					err = checkCol(err, "unit", string(col.Unit), string(unit))
					if err != nil {
						return Plot{}, err
					}
				}
			}
			if p[i].Caption.Rows() != len(p[i].Lines) {
				return Plot{}, fmt.Errorf("merge plots: caption rows != number of lines")
			}
		}
	}

	// Build default order.
	if order == nil {
		id := 0
		order = make(map[[2]int]int)
		for i := range p {
			for k := range p[i].Lines {
				order[[2]int{i, k}] = id
				id++
			}
		}
	}

	lines := make([]Line, len(order))
	for ik, n := range order {
		i := ik[0]
		k := ik[1]
		if i < 0 || i >= len(p) || k < 0 || k >= len(p[i].Lines) {
			return Plot{}, fmt.Errorf("merge plots: order[%d][%d] does not exist", i, k)
		}
		if n < 0 || n >= len(lines) {
			return Plot{}, fmt.Errorf("merge plots: order[%d][%d] target %d does not exist", i, k, n)
		}
		lines[n] = p[i].Lines[k]
	}
	for i := range lines {
		lines[i].Id = i
		if i < len(colors) {
			lines[i].Style.Line.Color = colors[i]
			lines[i].Style.Marker.Color = colors[i]
		}
	}

	if p[0].Caption != nil {
		// We cannot just append to the first caption, as the order may be off,
		// so we make a copy but remove the Data.
		p0c := p[0].Caption
		var targetCaption Caption = *p0c
		targetCaption.Columns = make([]CaptionColumn, len(p0c.Columns))
		for i, col := range p0c.Columns {
			targetCaption.Columns[i] = col
			targetCaption.Columns[i].Data = nil
		}

		for ik, n := range order {
			i := ik[0]         // plot index
			k := ik[1]         // row index
			ci := p[i].Caption // c != nil, tested above.
			if numrows := ci.Rows(); k >= numrows {
				return Plot{}, fmt.Errorf("merge plots: plot[%d].caption %d >= num rows (%d)", i, k, numrows)
			}

			for c := range ci.Columns {
				if ci.Columns[c].isEmpty == false {
					if val, err := ci.Columns[c].ValueAt(k); err != nil {
						return Plot{}, fmt.Errorf("merge plots: %s", err)
					} else if err := targetCaption.Columns[c].SetAt(n, val, len(lines)); err != nil {
						return Plot{}, fmt.Errorf("merge plots: %s", err)
					}
				}
			}
		}

		p[0].Caption = &targetCaption
	}

	ret := p[0]
	ret.Lines = lines
	return ret, nil
}
