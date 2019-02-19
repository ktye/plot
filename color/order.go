package color

import (
	"fmt"
	"strings"
)

// Color order is a palette used for different line colors.
// A color order is a list of comma separated color strings
// or a predifined order with an abbreviated name.
// If Order is empty, a default order is returned.
type Order string

// Get returns the color with index int.
// The index should start at 1. If it is 0 it is considered unset and the
// defaultNum is used.
// It is cyclic and cannot overflow.
func (co Order) Get(index int, defaultNum int) Color {
	if index == 0 {
		index = defaultNum
	}
	index--
	colors := co.split()
	if index < 0 {
		return Color("black")
	}
	return Color(colors[index%len(colors)])
}

// NumColors returns the number of colors in the color order.
func (co Order) NumColors() int {
	colors := co.split()
	return len(colors)
}

func (co Order) Verify() error {
	colors := co.split()
	for i, c := range colors {
		if err := Color(c).Verify(); err != nil {
			return fmt.Errorf("illegal color order: #%d: %s\n", i, err)
		}
	}
	return nil
}

// split returns the color vector as strings.
func (co Order) split() []string {
	s := string(co)
	if s == "" {
		s = "bright"
	}
	s = abbrevated(s)
	return strings.Split(s, ",")
}

// GetPredefinedOrders returns a list of predefined color orders.
func GetPredefinedOrders() []string {
	return []string{"deep", "muted", "pastel", "bright", "dark", "blind"}
}

// abbrevated returns a list of abbrevated built-in color orders.
func abbrevated(s string) string {
	switch s {
	case "deep":
		return "#4C72B0,#55A868,#C44E52,#8172B2,#CCB974,#64B5CD"
	case "muted":
		return "#4878CF,#6ACC65,#D65F5F,#B47CC7,#C4AD66,#77BEDB"
	case "pastel":
		return "#92C6FF,#97F0AA,#FF9F9A,#D0BBFF,#FFFEA3,#B0E0E6"
	case "bright":
		return "#003FFF,#03ED3A,#E8000B,#8A2BE2,#FFC400,#00D7FF"
	case "dark":
		return "#001C7F,#017517,#8C0900,#7600A1,#B8860B,#006374"
	case "blind":
		return "#0072B2,#009E73,#D55E00,#CC79A7,#F0E442,#56B4E9"
	}
	return s
}
