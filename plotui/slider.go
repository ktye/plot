package plotui

// SetSlider sets the slider value and the min max range.
// It is set to line or point mode, depending on the value of ui.slidePoints.
func (ui *Plot) SetSlider(value int) {
	if ui.slider == nil {
		return
	}
	if ui.slidePoints == 0 {
		ui.model.Lock()
		rows := ui.model.RowCount()
		ui.model.Unlock()
		maxValues := rows - ui.lineOffset
		if maxValues <= 1 {
			ui.slider.SetEnabled(false)
			ui.slider.SetRange(-1, -1)
		} else {
			ui.slider.SetEnabled(true)
			ui.slider.SetRange(1, maxValues)
		}
	} else {
		ui.slider.SetEnabled(true)
		ui.slider.SetRange(1, ui.slidePoints)
	}
	ui.ignore = true
	ui.slider.SetValue(value)
	ui.ignore = false
}

func (ui *Plot) sliderChanged() {
	if ui.ignore {
		return
	}
	val := ui.slider.Value()
	if ui.slidePoints == 0 {
		if ui.table != nil {
			ui.table.SetSelectedIndexes([]int{val + ui.lineOffset - 1})
		}
	} else if len(ui.hi) == 1 {
		ui.hi[0].Point = val - 1
		ui.update(ui.hi)
	}
}
