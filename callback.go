package plot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ktye/plot/xmath"
)

// Callback is returns on a double click.
// It's type is context sensitive. It depends on the are where it has been clicked.
// It may be "" (PointInfo), AxisCallback, UnitCallback
// If Type is empty, the point info is valid, otherwise it has been clicked on a special area.
type Callback struct {
	Type      string    // PointInfoCallback AxisCallback, UnitCallback
	PlotIndex int       // Plot index, filled by ClickIPlotters
	PointInfo PointInfo // PointInfo, valid only for Type == ""
	Limits    Limits    // Axis Limits, valid only for Type == AxisCallback
}

const (
	PointInfoCallback string = ""
	AxisCallback      string = "Axis"
	UnitCallback      string = "Unit"
	MeasurePoint      string = "Measure Point"
	DeleteLine        string = "Delete Line"
)

type HighlightID struct {
	Line           int     // line number starting at 0
	Point          int     // point number, -1 is complete line, NaN's don't count.
	XImage, YImage float64 // image coordinates that have been clicked.
}

func (h HighlightID) String() string {
	if h.XImage == 0 && h.YImage == 0 {
		return fmt.Sprintf("%v,%v", h.Line, h.Point)
	} else {
		return fmt.Sprintf("%v,%v,%v,%v", h.Line, h.Point, h.XImage, h.YImage)
	}
}

func ParseHighlightID(s string) (HighlightID, error) {
	hi := HighlightID{}
	v := strings.Split(s, ",")
	if len(v) != 2 && len(v) != 4 {
		return hi, fmt.Errorf("cannot parse highlight id: 2 or 4 elements are needed: '%s'", s)
	}
	if len(v) >= 2 {
		if n, err := strconv.Atoi(v[0]); err != nil {
			return hi, err
		} else {
			hi.Line = n
		}
		if n, err := strconv.Atoi(v[1]); err != nil {
			return hi, err
		} else {
			hi.Point = n
		}
	}
	if len(v) == 4 {
		if n, err := strconv.ParseFloat(v[2], 64); err != nil {
			return hi, err
		} else {
			hi.XImage = n
		}
		if n, err := strconv.ParseFloat(v[3], 64); err != nil {
			return hi, err
		} else {
			hi.YImage = n
		}
	}
	return hi, nil
}

// PointInfo is the result for a click event.
// If a user clicks on a point in the plot, this is returned.
type PointInfo struct {
	LineID      int        // Id of the line being clicked.
	PointNumber int        // Number of the data point within the line, NaNs are ignored when counting, otherwise amp/ang plots are not aligned.
	NumPoints   int        // Number of data points of the clicked line.
	IsEnvelope  bool       // Mark if the line type is an envelope.
	X, Y, Z     float64    // Points coordinates, Z is only used for images.
	C           complex128 // Complex data value of the line point, if it has one.
	IsImage     bool       // Mark if the PlotType is an Image.
	Zmin, Zmax  float64    // Min, Max value of zaxis.
}

func (p PointInfo) String() string {
	if p.IsImage {
		order := p.Y / real(p.C)
		return fmt.Sprintf("x=%v y=%v c=%v z=%v scale=[%v,%v] Order=%.1f\n", p.X, p.Y, real(p.C), p.Z, p.Zmin, p.Zmax, order)
	} else if p.IsEnvelope {
		return fmt.Sprintf("line#%d point %d/%d x=%v ymin=%v ymax=%v\n", p.LineID+1, p.PointNumber+1, p.NumPoints, p.X, real(p.C), imag(p.C))
	} else {
		return fmt.Sprintf("line#%d point %d/%d x=%v y=%v c=%v\n", p.LineID+1, p.PointNumber+1, p.NumPoints, p.X, p.Y, xmath.Absang(p.C, ""))
	}
}
