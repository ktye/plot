// Package plotui provides a windows gui frontend to plot.Plots
package plotui

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"strings"
	"sync"

	"github.com/ktye/plot"
	"github.com/lxn/walk"
	"github.com/lxn/win"
)

// SetCaption updates the caption table.
func (ui *Plot) SetCaption() {
	if s, lineOffset, err := ui.CaptionStrings(); err != nil {
		log.Print(err)
		return
	} else if len(s) == 0 {
		ui.table.SetVisible(false)
	} else {
		ui.table.SetVisible(true)
		ui.lineOffset = lineOffset
		if ui.plots != nil && len(*ui.plots) > 0 {
			ui.model.Lock()
			ui.model.colors = make([]color.Color, len(s))
			ui.model.lineOffset = lineOffset

			for i := 0; i < len(ui.model.colors); i++ {
				co := ui.caption.Color(i, lineOffset)
				ui.model.colors[i] = co
				if ui.model.images == nil {
					ui.model.images = make(map[color.Color]*walk.Bitmap)
				}
				if co != nil {
					if _, ok := ui.model.images[co]; !ok {
						im := image.NewRGBA(image.Rect(0, 0, smallIconWidth, smallIconHeight))
						draw.Draw(im, im.Bounds(), &image.Uniform{co}, image.ZP, draw.Src)
						if bm, err := walk.NewBitmapFromImage(im); err == nil {
							ui.model.images[co] = bm
						} else {
							log.Print("cannot create bitmap: ", err)
						}
					}
				}
			}

			ui.model.Unlock()
		}
		ui.ignore = true
		ui.table.SetSelectedIndexes(nil)
		ui.ignore = false
		ui.model.lines = s
		ui.model.PublishRowsReset()
		ui.table.StretchLastColumn()
	}
}

func (ui *Plot) lineClicked() {
	if ui.ignore {
		return
	}
	selectedIDs := func() []plot.HighlightID {
		selected := ui.table.SelectedIndexes()
		var ids []plot.HighlightID
		for _, id := range selected {
			if n := id - ui.lineOffset; n >= 0 {
				ids = append(ids, plot.HighlightID{Line: n, Point: -1})
			}
		}
		return ids
	}

	ids := selectedIDs()
	if len(ids) == 1 {
		ui.slidePoints = 0
		ui.SetSlider(ids[0].Line + 1)
	}
	ui.hi = ids
	ui.update(ids)
}

func (ui *Plot) CaptionHeader() []string {
	if ui.caption == nil {
		return nil
	}
	return ui.caption.LeadText
}

func (ui *Plot) CaptionStrings() ([]string, int, error) {
	if ui.caption == nil {
		return nil, 0, fmt.Errorf("caption is empty")
	}
	var b bytes.Buffer
	if lineOffset, err := ui.caption.WriteTable(&b, plot.Numbers); err != nil {
		return nil, 0, err
	} else {
		s := strings.Split(string(b.Bytes()), "\n")
		if len(s) > 0 && len(s[len(s)-1]) == 0 {
			s = s[:len(s)-1]
		}
		return s, lineOffset, nil
	}
}

func (ui *Plot) ClipboardCaption() error {
	v, _, err := ui.CaptionStrings()
	if err != nil {
		return err
	}
	s := strings.Join(v, "\r\n")
	cb := walk.Clipboard()
	if err := cb.Clear(); err != nil {
		return err
	}
	if err := cb.SetText(s); err != nil {
		return err
	}
	return nil
}

func (ui *Plot) highlightCaption(ids []int) {
	if ui.table == nil {
		return
	}
	for i := range ids {
		ids[i] += ui.lineOffset
	}
	ui.ignore = true
	ui.table.SetSelectedIndexes(ids)
	ui.ignore = false
}

type CaptionModel struct {
	sync.Mutex
	walk.TableModelBase
	lines      []string
	colors     []color.Color // Color for each row including lead text.
	lineOffset int
	images     map[color.Color]*walk.Bitmap
}

func (m *CaptionModel) RowCount() int {
	return len(m.lines)
}

func (m *CaptionModel) Value(row, col int) interface{} {
	if row >= len(m.lines) {
		log.Printf("requested caption row number out of range: row %d, col %d", row, col)
		return ""
	}
	return m.lines[row]
}

func (m *CaptionModel) Image(row int) interface{} {
	m.Lock()
	defer m.Unlock()
	if row < 0 || row >= len(m.colors) {
		return whiteBitmap
	}
	co := m.colors[row]
	if bm, ok := m.images[co]; ok {
		return bm
	}
	return whiteBitmap
}

var whiteBitmap *walk.Bitmap
var smallIconWidth, smallIconHeight int

func init() {
	smallIconWidth, smallIconHeight = int(win.GetSystemMetrics(win.SM_CXSMICON)), int(win.GetSystemMetrics(win.SM_CYSMICON))
	im := image.NewRGBA(image.Rect(0, 0, smallIconWidth, smallIconHeight))
	draw.Draw(im, im.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	if bm, err := walk.NewBitmapFromImage(im); err == nil {
		whiteBitmap = bm
	}
}
