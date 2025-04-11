package plot

// LineStyle is the line style definition.
type LineStyle struct {
	Width    int // line width
	Color    int // line color index.
	Arrow    int // arrow head if != 0
	EndMarks int // line end markers (horizontal/vertical only)
}

// Marker may be Point | Circle | Cross | Vertical.
type Marker string

const (
	NoMarker     Marker = ""
	PointMarker  Marker = "point"
	CircleMarker Marker = "circle"
	CrossMarker  Marker = "cross"
	Bar          Marker = "bar"
)

type MarkerStyle struct {
	Marker Marker // marker type point/circle/cross/vertical
	Size   int    // marker size (diameter)
	Color  int    // marker color index
}
