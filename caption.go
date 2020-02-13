package plot

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"math"
	"math/cmplx"
	"strings"
	"text/tabwriter"

	"github.com/ktye/plot/xmath"
)

const (
	Numbers    uint = 1 << iota // Add '#' column (run numbers).
	HtmlColors                  // Add a block character with the data color.
)

// Caption represents the caption table of a Plot.
// Multiple captions can be merged, if the plots are "similar".
type Caption struct {
	Title    string          // title of the caption
	LeadText []string        // caption header lines (free text before actual table).
	Columns  []CaptionColumn // column vector, len: #rows
	colors   []color.Color   // row color, set by SetCaptionColors (ignoring Title and LeadText)
}

// CaptionColumn is a single column of a Caption object.
type CaptionColumn struct {
	Class   string      // column class id
	Name    string      // column title
	Unit    string      // column unit
	Format  string      // numeric format for fmt.Fprint for the data
	Data    interface{} // data must be any of the types: []string, []int, []float64, []complex128
	isEmpty bool        // must be checked with SetEmptyColumns
}

// SetAt adds a data value to a caption column.
// If the data slice is nil, it is initialized to the given length.
func (c *CaptionColumn) SetAt(i int, data interface{}, maxlen int) error {
	if c.Data == nil {
		switch v := data.(type) {
		case string:
			c.Data = make([]string, maxlen)
		case float64:
			c.Data = make([]float64, maxlen)
		case int:
			c.Data = make([]int, maxlen)
		case complex128:
			c.Data = make([]complex128, maxlen)
		default:
			return errUnknownColType(v)
		}
	}

	switch coldata := c.Data.(type) {
	case []string:
		if v, ok := data.(string); ok {
			coldata[i] = v
		} else {
			return fmt.Errorf("cannot set to caption data: expected a string value, got %v", data)
		}
	case []int:
		if v, ok := data.(int); ok {
			coldata[i] = v
		} else {
			return fmt.Errorf("cannot set to caption data: expected an int value, got %v", data)
		}
	case []float64:
		if v, ok := data.(float64); ok {
			coldata[i] = v
		} else {
			return fmt.Errorf("cannot set to caption data: expected a float value, got %v", data)
		}
	case []complex128:
		if v, ok := data.(complex128); ok {
			coldata[i] = v
		} else {
			return fmt.Errorf("cannot set to caption data: expected a complex value, got %v", data)
		}
	default:
		return errUnknownColType(coldata)
	}
	return nil
}

func (c CaptionColumn) ValueAt(i int) (interface{}, error) {
	max := 0
	switch v := c.Data.(type) {
	case []string:
		max = len(v)
		if i >= 0 && i < len(v) {
			return v[i], nil
		}
	case []int:
		max = len(v)
		if i >= 0 && i < len(v) {
			return v[i], nil
		}
	case []float64:
		max = len(v)
		if i >= 0 && i < len(v) {
			return v[i], nil
		}
	case []complex128:
		max = len(v)
		if i >= 0 && i < len(v) {
			return v[i], nil
		}
	default:
		return nil, errUnknownColType(v)
	}
	return nil, fmt.Errorf("caption column value: column '%s' index out of range: %d [0,%d)", c.Name, i, max)
}

func errUnknownColType(v interface{}) error {
	return fmt.Errorf("internal error: caption column type is unsupported: %T", v)
}

// Rows returns the number of data rows of a caption, excluding the lead text.
func (c *Caption) Rows() int {
	if len(c.Columns) < 1 {
		return 0
	}
	data := c.Columns[0].Data
	switch t := data.(type) {
	case []string:
		return len(data.([]string))
	case []int:
		return len(data.([]int))
	case []float64:
		return len(data.([]float64))
	case []complex128:
		return len(data.([]complex128))
	default:
		panic(fmt.Sprintf("table data has unknown type: %s", t))
	}
	return 0
}

// RemoveColumns removes columns from the caption.
func (c *Caption) RemoveColumns(remCols map[string]bool) {
	rem := 0
	for _, v := range remCols {
		if v == true {
			rem++
		}
	}
	if rem == 0 {
		return
	}
	var columns []CaptionColumn
	for _, c := range c.Columns {
		if v := remCols[c.Class]; v == false {
			columns = append(columns, c)
		}
	}
	c.Columns = columns
}

// Color returns the color marker for the given row (including the lead text),
// which may be nil.
func (c *Caption) Color(row int, lineOffset int) color.Color {
	row -= lineOffset
	if row < 0 {
		return nil
	}
	if row >= len(c.colors) {
		return nil
	}
	return c.colors[row]
}

// WriteTable writes an aligned text table to the writer.
func (c *Caption) WriteTable(w io.Writer, flags uint) (int, error) {
	// lineOffset is the offset of the first data line.
	lineOffset := 0

	// Write Title.
	if c.Title != "" {
		lineOffset++
		fmt.Fprintln(w, c.Title)
	}

	// Write LeadText. If it does not end with an empty line, append one.
	leadText := c.LeadText
	if len(leadText) > 0 && leadText[len(leadText)-1] != "" {
		leadText = append(leadText, "")
	}
	for _, s := range leadText {
		fmt.Fprintln(w, s)
	}
	lineOffset += len(leadText)
	if len(c.Columns) == 0 {
		return lineOffset, nil
	}

	// Build the table.
	var b bytes.Buffer
	tw := tabwriter.NewWriter(&b, 0, 8, 2, ' ', 0)
	c.SetEmptyColumns()

	// Write header.
	if flags&Numbers != 0 {
		fmt.Fprintf(tw, "#\t")
	}
	for i, column := range c.Columns {
		ap := "\t"
		if i == len(c.Columns)-1 {
			ap = "\n"
		}
		if column.isEmpty {
			if i == len(c.Columns)-1 {
				fmt.Fprintf(tw, "\n")
			}
			continue
		}
		fmt.Fprintf(tw, "%s", column.Name+ap)
	}
	lineOffset++

	// Write units.
	noUnits := true
	for _, column := range c.Columns {
		if column.isEmpty {
			continue
		}
		if column.Unit != "" {
			noUnits = false
			break
		}
	}
	if noUnits == false {
		if flags&Numbers != 0 {
			fmt.Fprintf(tw, "\t")
		}
		for i, column := range c.Columns {
			ap := "\t"
			if i == len(c.Columns)-1 {
				ap = "\n"
			}
			if column.isEmpty {
				if i == len(c.Columns)-1 {
					fmt.Fprintf(tw, "\n")
				}
				continue
			}
			fmt.Fprintf(tw, "%s"+ap, string(column.Unit))
		}
		lineOffset++
	}
	nonEmptyColumns := 0
	for _, column := range c.Columns {
		if !column.isEmpty {
			nonEmptyColumns++
		}
	}

	// Row must be defined eary, otherwise goto jumps over declaration.
	row := 0

	if len(c.Columns) == 0 || nonEmptyColumns == 0 {
		if flags&Numbers != 0 {
			for i := 0; i < c.Rows(); i++ {
				fmt.Fprintf(tw, "%d\n", i+1)
			}
		}
		goto E
	}

	// Write Content.
	for {
		i := 0
		for _, column := range c.Columns {
			if column.isEmpty {
				continue
			}
			i++
			ap := "\t"
			if i == nonEmptyColumns {
				ap = "\n"
			}
			var s string
			format := "%v"
			if column.Format != "" {
				format = column.Format
			}
			switch t := column.Data.(type) {
			case []int:
				v := column.Data.([]int)
				if row == len(v) {
					goto E
				}
				s = fmt.Sprintf(format, v[row])
			case []float64:
				v := column.Data.([]float64)
				if row == len(v) {
					goto E
				}
				if math.IsNaN(v[row]) {
					s = ""
				} else {
					s = fmt.Sprintf(format, v[row])
				}
			case []complex128:
				if column.Format == "" {
					format = "%v@%v"
				}
				v := column.Data.([]complex128)
				if row == len(v) {
					goto E
				}
				if cmplx.IsNaN(v[row]) {
					s = ""
				} else {
					s = xmath.Absang(v[row], format)
				}
			case []string:
				v := column.Data.([]string)
				if row == len(v) {
					goto E
				}
				s = v[row]
			default:
				return lineOffset, fmt.Errorf("table data has unknown type: %s", t)
			}
			if i == 1 && flags&Numbers != 0 {
				fmt.Fprintf(tw, "%d\t", row+1)
			}
			fmt.Fprintf(tw, s+ap)
		}
		row++
	}

E:
	tw.Flush()

	if flags&HtmlColors != 0 {
		lines := strings.Split(string(b.Bytes()), "\n")
		for i, s := range lines {
			color := "white"
			if co := c.Color(i, lineOffset); co != nil {
				r, g, b, _ := co.RGBA()
				color = fmt.Sprintf("#%02X%02X%02X", r>>8, g>>8, b>>8)
			}
			fmt.Fprintf(w, "<span style=\"color:%s\">&#x25A0; </span>%s\n", color, s)
		}
	} else {
		w.Write(b.Bytes())
	}
	return lineOffset, nil
}

// MergeCaptions merges multiple captions to a single Caption table.
// It returns an error, if the captions are not uniform.
// Columns at the end may be the same, in which case they are included only once.
func MergeCaptions(captions []Caption) (single Caption, err error) {
	// Algorithm:
	// Merging two captions checks if any columns (starting at the end)
	// are equal. The non-equal columns are appended next to each
	// other, followed by the equal column only once.

	if len(captions) == 0 {
		return single, nil
	}
	single = captions[0]
	for i := 1; i < len(captions); i++ {
		single, err = mergeTwoCaptions(single, captions[i])
		if err != nil {
			return single, err
		}
	}
	return single, nil
}

// SplitLeadText splits the text at newlines but removes a final newline.
func SplitLeadText(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(strings.TrimSuffix(s, "\n"), "\n")
}

func mergeTwoCaptions(left, right Caption) (merged Caption, err error) {
	// See where the captions differ.
	lidx, ridx := len(left.Columns), len(right.Columns)
	for {
		if lidx < 1 || ridx < 1 {
			break
		}
		if columnsDiffer(left.Columns[lidx-1], right.Columns[ridx-1]) {
			break
		}
		lidx--
		ridx--
	}
	merged = left
	merged.Columns = merged.Columns[:lidx]
	merged.Columns = append(merged.Columns, right.Columns...)
	if len(merged.LeadText) == 0 {
		merged.LeadText = right.LeadText
	} else if len(right.LeadText) > 0 {
		var b bytes.Buffer
		w := tabwriter.NewWriter(&b, 4, 8, 2, ' ', 0)
		l := len(merged.LeadText)
		if len(right.LeadText) > l {
			l = len(right.LeadText)
		}
		for i := 0; i < l; i++ {
			var le, ri string
			if i < len(merged.LeadText) {
				le = merged.LeadText[i]
			}
			if i < len(right.LeadText) {
				ri = right.LeadText[i]
			}
			fmt.Fprintf(w, "%s\t%s\n", le, ri)
		}
		w.Flush()
		merged.LeadText = SplitLeadText(string(b.Bytes()))
	}
	return merged, nil
}

// columnsDiffer checks two columns for inequality.
// The number of rows has been checked before.
// The actual data section is not checked.
func columnsDiffer(c1, c2 CaptionColumn) bool {
	if c1.Name == c2.Name && c1.Class == c2.Class {
		return false
	}
	return true
}

// SetEmptyColumns checks if the data column only has empty values and marks these columns.
func (c *Caption) SetEmptyColumns() {
	for i, col := range c.Columns {
		switch t := col.Data.(type) {
		case []string:
			c.Columns[i].isEmpty = allEmptyStrings(col.Data.([]string))
		case []int:
			c.Columns[i].isEmpty = false // ints have no empty value
		case []float64:
			c.Columns[i].isEmpty = allEmptyFloats(col.Data.([]float64))
		case []complex128:
			c.Columns[i].isEmpty = allEmptyComplex(col.Data.([]complex128))
		case nil:
			c.Columns[i].isEmpty = true
		default:
			panic(fmt.Sprintf("unknown caption data type: '%v'", t))
		}
	}
}

func allEmptyStrings(v []string) bool {
	for _, s := range v {
		if s != "" {
			return false
		}
	}
	return true
}

func allEmptyFloats(v []float64) bool {
	for _, f := range v {
		if !(math.IsNaN(f) || f == 0) {
			return false
		}
	}
	return true
}

func allEmptyComplex(v []complex128) bool {
	for _, z := range v {
		if !cmplx.IsNaN(z) {
			return false
		}
	}
	return true
}

// SetCaptionColors sets the colors of the caption associated with the plot.
func (p *Plot) SetCaptionColors() {
	if p.Caption == nil {
		return
	}
	colors := make([]color.Color, p.Caption.Rows())
	for i, l := range p.Lines {
		co := p.Style.Order.Get(l.Style.Line.Color, i+1).Color()
		if i < len(colors) {
			colors[i] = co
		}
	}
	p.Caption.colors = colors
}
